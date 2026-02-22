package apihttp

import "github.com/andrieee44/hackusc/service"

type userUpdateRequestPayload struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func newUserUpdateRequest(
	req userUpdateRequestPayload,
) service.UserUpdateRequest {
	return service.UserUpdateRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}
