package dto

type LoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
