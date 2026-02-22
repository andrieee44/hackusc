package main

import (
	"fmt"
	"net/http"

	"github.com/andrieee44/hackusc/api/http"
	"github.com/andrieee44/hackusc/auth"
	"github.com/andrieee44/hackusc/service"
	"github.com/andrieee44/hackusc/store"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	var (
		userHTTPHandler *apihttp.UserHTTPHandler
		mux             *http.ServeMux
	)

	userHTTPHandler = apihttp.NewUserHTTPHandler(
		auth.NewJWTSigner[service.UserActor](
			jwt.SigningMethodHS256,
			[]byte("test"),
		),
		service.NewUserService(
			store.NewMemStore(),
			auth.Bcrypt{},
		),
	)

	mux = http.NewServeMux()

	mux.HandleFunc(
		"POST /users",
		userHTTPHandler.Create,
	)

	mux.HandleFunc(
		"GET /users/{id}",
		userHTTPHandler.GetByID,
	)

	mux.HandleFunc(
		"GET /users/email/{email}",
		userHTTPHandler.GetByEmail,
	)

	mux.HandleFunc(
		"PATCH /users",
		userHTTPHandler.Update,
	)

	mux.HandleFunc(
		"POST /auth/login",
		userHTTPHandler.Login,
	)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
