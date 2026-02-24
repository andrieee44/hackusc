package apihttp

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/andrieee44/hackusc/domain"
	"github.com/andrieee44/hackusc/service"
)

type UserTokenSigner interface {
	Sign(service.UserActor, time.Duration) (string, error)
	Extract(string) (service.UserActor, error)
}

type UserHTTPHandler struct {
	userTokenSigner UserTokenSigner
	userService     *service.UserService
}

const userTokenTTL = 24 * time.Hour

var ErrExpiredUserToken = errors.New("user token is expired")

func NewUserHTTPHandler(
	userTokenSigner UserTokenSigner,
	userService *service.UserService,
) *UserHTTPHandler {
	return &UserHTTPHandler{userTokenSigner, userService}
}

func (api *UserHTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	type reqJson struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var (
		req   reqJson
		actor service.UserActor
		token string
		err   error
	)

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeErrorStatus(w, http.StatusBadRequest)

		return
	}

	actor, err = api.userService.Create(
		r.Context(),
		req.Name,
		req.Email,
		req.Password,
	)
	if err != nil {
		writeErrorStatus(w, errToStatus(err))

		return
	}

	token, err = api.userTokenSigner.Sign(actor, userTokenTTL)
	if err != nil {
		println(err.Error())
		writeErrorStatus(w, http.StatusInternalServerError)

		return
	}

	writeJSON(w, http.StatusCreated, newUserTokenPayload(token))
}

func (api *UserHTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	var (
		id      int64
		userDTO service.UserDTO
		err     error
	)

	id, err = strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeErrorStatus(w, http.StatusBadRequest)

		return
	}

	userDTO, err = api.userService.GetByID(r.Context(), id)
	if err != nil {
		writeErrorStatus(w, errToStatus(err))

		return
	}

	writeJSON(w, http.StatusOK, newUserDTOResponse(userDTO))
}

func (api *UserHTTPHandler) GetByEmail(
	w http.ResponseWriter,
	r *http.Request,
) {
	var (
		userDTO service.UserDTO
		err     error
	)

	userDTO, err = api.userService.GetByEmail(
		r.Context(),
		r.PathValue("email"),
	)
	if err != nil {
		writeErrorStatus(w, errToStatus(err))

		return
	}

	writeJSON(w, http.StatusOK, newUserDTOResponse(userDTO))
}

func (api *UserHTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		actor service.UserActor
		req   userUpdateRequestPayload
		ok    bool
		err   error
	)

	actor, ok = api.requireActor(w, r)
	if !ok {
		return
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeErrorStatus(w, http.StatusBadRequest)

		return
	}

	err = api.userService.Update(
		r.Context(),
		actor,
		newUserUpdateRequest(req),
	)
	if err != nil {
		writeErrorStatus(w, errToStatus(err))

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *UserHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	type reqJson struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var (
		req   reqJson
		actor service.UserActor
		token string
		err   error
	)

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeErrorStatus(w, http.StatusBadRequest)

		return
	}

	actor, err = api.userService.Login(
		r.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		writeErrorStatus(w, errToStatus(err))

		return
	}

	token, err = api.userTokenSigner.Sign(actor, userTokenTTL)
	if err != nil {
		writeErrorStatus(w, http.StatusInternalServerError)

		return
	}

	writeJSON(w, http.StatusOK, newUserTokenPayload(token))
}

func (api *UserHTTPHandler) requireActor(
	w http.ResponseWriter,
	r *http.Request,
) (service.UserActor, bool) {
	const bearerPrefix string = "Bearer "

	var (
		authHeader string
		token      string
		actor      service.UserActor
		err        error
	)

	authHeader = r.Header.Get("Authorization")
	if authHeader == "" {
		writeErrorStatus(w, http.StatusUnauthorized)

		return service.UserActor{}, false
	}

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		writeErrorStatus(w, http.StatusBadRequest)

		return service.UserActor{}, false
	}

	token = strings.TrimPrefix(authHeader, bearerPrefix)

	actor, err = api.userTokenSigner.Extract(token)
	if err != nil {
		if errors.Is(err, ErrExpiredUserToken) {
			writeErrorStatus(w, http.StatusUnauthorized)
		} else {
			writeErrorStatus(w, http.StatusInternalServerError)
		}

		return service.UserActor{}, false
	}

	return actor, true
}

func errToStatus(err error) int {
	switch {
	case errors.Is(err, service.ErrUserAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, service.ErrUserDoesNotExist):
		return http.StatusNotFound
	case errors.Is(err, service.ErrUserLacksRole):
		return http.StatusForbidden
	case errors.Is(err, service.ErrIncorrectPassword):
		return http.StatusUnauthorized
	case errors.Is(err, service.ErrEmptyUserUpdateRequest):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrInvalidArgument):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrInternal):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func writeErrorStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func writeJSON(w http.ResponseWriter, successStatus int, v any) {
	var (
		payload []byte
		err     error
	)

	payload, err = json.Marshal(v)
	if err != nil {
		writeErrorStatus(w, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(successStatus)

	_, err = w.Write(payload)
	if err != nil {
		log.Printf("response write error: %v", err)
	}
}
