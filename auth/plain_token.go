package auth

import (
	"encoding/json"
	"time"

	"github.com/andrieee44/hackusc/api/http"
)

type PlainToken[T any] struct{}

type PlainTokenMetadata struct {
	Exp int64           `json:"exp"`
	Val json.RawMessage `json:"val"`
}

func (PlainToken[T]) Sign(v T, ttl time.Duration) (string, error) {
	var (
		val     json.RawMessage
		payload []byte
		err     error
	)

	val, err = json.Marshal(v)
	if err != nil {
		return "", err
	}

	payload, err = json.Marshal(
		PlainTokenMetadata{time.Now().Add(ttl).Unix(), val},
	)
	if err != nil {
		return "", err
	}

	return string(payload), nil
}

func (PlainToken[T]) Extract(token string) (T, error) {
	var (
		metadata PlainTokenMetadata
		v        T
		err      error
	)

	err = json.Unmarshal([]byte(token), &metadata)
	if err != nil {
		return *new(T), err
	}

	if time.Now().After(time.Unix(metadata.Exp, 0)) {
		return *new(T), apihttp.ErrExpiredUserToken
	}

	err = json.Unmarshal(metadata.Val, &v)
	if err != nil {
		return *new(T), err
	}

	return v, nil
}
