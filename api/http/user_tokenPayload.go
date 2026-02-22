package apihttp

type UserTokenPayload struct {
	Token string `json:"token"`
}

func newUserTokenPayload(token string) UserTokenPayload {
	return UserTokenPayload{token}
}
