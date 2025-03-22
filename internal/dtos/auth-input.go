package dtos
type AuthLoginInput struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type AuthRegisterInput struct {
	UserName string `json:"userName"`
	Email string `json:"email"`
	Password string `json:"password"`
}