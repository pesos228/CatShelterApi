package dto

type CatResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int16  `json:"age"`
}

type CatRequest struct {
	Name string `json:"name"`
	Age  int16  `json:"age"`
}

type CatsPaginatedResponse struct {
	Data       []CatResponse    `json:"data"`
	Pagination PaginationResult `json:"pagination"`
}
