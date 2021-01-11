package services

type CCXLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func CCXAuth(Login string, Password string) *CCXLogin {
	return &CCXLogin{
		Login:    Login,
		Password: Password,
	}
}
