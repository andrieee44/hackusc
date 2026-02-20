package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andrieee44/hackusc/domain"
	"github.com/andrieee44/hackusc/service"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type reqJson struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var (
			req  reqJson
			hash []byte
			id   int64
			err  error
		)

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		defer r.Body.Close()

		hash, err = bcrypt.GenerateFromPassword(
			[]byte(req.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		id, err = userService.CreateUser(
			r.Context(),
			domain.WithName(req.Name),
			domain.WithEmail(req.Email),
			domain.WithPasswordHash(hash),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(struct {
			ID int64 `json:"id"`
		}{id})
	}
}

func GetUserByID(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			id  int64
			u   domain.User
			err error
		)

		id, err = strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		u, err = userService.GetUserByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(struct {
			ID    int64  `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}{u.ID, u.Name, u.Email})
	}
}

func GetUserByEmail(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			u   domain.User
			err error
		)

		u, err = userService.GetUserByEmail(r.Context(), r.PathValue("email"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(struct {
			ID    int64  `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}{u.ID, u.Name, u.Email})
	}
}

func UpdateUser(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type reqJson struct {
			Name     *string `json:"name"`
			Email    *string `json:"email"`
			Password *string `json:"password"`
		}

		var (
			id   int64
			req  reqJson
			opts []domain.UserOpt
			hash []byte
			err  error
		)

		id, err = strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		defer r.Body.Close()

		if req.Name != nil {
			opts = append(opts, domain.WithName(*req.Name))
		}

		if req.Email != nil {
			opts = append(opts, domain.WithEmail(*req.Email))
		}

		if req.Password != nil {
			hash, err = bcrypt.GenerateFromPassword(
				[]byte(*req.Password),
				bcrypt.DefaultCost,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			opts = append(opts, domain.WithPasswordHash(hash))
		}

		err = userService.UpdateUser(r.Context(), id, opts...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
