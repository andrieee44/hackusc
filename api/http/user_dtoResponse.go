package apihttp

import "github.com/andrieee44/hackusc/service"

type userDTOResponse struct {
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	ID    int64    `json:"id"`
}

func newUserDTOResponse(userDTO service.UserDTO) userDTOResponse {
	return userDTOResponse{
		userDTO.Name,
		userDTO.Email,
		userDTO.Roles,
		userDTO.ID,
	}
}
