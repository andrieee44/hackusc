package main

import (
	"fmt"
	"net/http"

	"github.com/andrieee44/hackusc/api"
	"github.com/andrieee44/hackusc/service"
	"github.com/andrieee44/hackusc/store"
)

func main() {
	var (
		userService *service.UserService
		mux         *http.ServeMux
	)

	userService = service.NewUserService(store.NewMemStore())

	mux = http.NewServeMux()
	mux.HandleFunc("POST /users", api.CreateUser(userService))
	mux.HandleFunc("GET /users/{id}", api.GetUserByID(userService))
	mux.HandleFunc("GET /users/email/{email}", api.GetUserByEmail(userService))
	mux.HandleFunc("PATCH /users/{id}", api.UpdateUser(userService))

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
