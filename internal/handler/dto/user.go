package dto

type UserInfoResponse struct {
	Id    string         `json:"id"`
	Name  string         `json:"name"`
	Login string         `json:"login"`
	Roles []RoleResponse `json:"roles"`
	Cats  []CatResponse  `json:"cats"`
}

type ShelterCatRequest struct {
	Id string `json:"id"`
}
