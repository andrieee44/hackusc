package auth

import (
	"errors"
	"time"

	"github.com/andrieee44/hackusc/api/http"
	"github.com/golang-jwt/jwt/v5"
)

type customJWTClaims[T any] struct {
	jwt.RegisteredClaims
	Val T `json:"val"`
}

type JWTSigner[T any] struct {
	method jwt.SigningMethod
	secret any
}

func NewJWTSigner[T any](method jwt.SigningMethod, secret any) *JWTSigner[T] {
	return &JWTSigner[T]{method, secret}
}

func (s *JWTSigner[T]) Sign(v T, ttl time.Duration) (string, error) {
	return jwt.NewWithClaims(s.method, customJWTClaims[T]{
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
		v,
	}).SignedString(s.secret)
}

func (s *JWTSigner[T]) Extract(tokenString string) (T, error) {
	var (
		token  *jwt.Token
		claims *customJWTClaims[T]
		ok     bool
		err    error
	)

	token, err = jwt.ParseWithClaims(
		tokenString,
		new(customJWTClaims[T]),
		func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != s.method.Alg() {
				return nil, jwt.ErrTokenSignatureInvalid
			}

			return s.secret, nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return *new(T), apihttp.ErrExpiredUserToken
		}

		return *new(T), err
	}

	claims, ok = token.Claims.(*customJWTClaims[T])
	if !ok {
		return *new(T), jwt.ErrTokenInvalidClaims
	}

	return claims.Val, nil
}
