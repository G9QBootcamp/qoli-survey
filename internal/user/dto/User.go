package dto

type UserResponse struct {
	Name string `json:"name"`
}

type UserRequest struct {
	Page int
	Name string
}

type UserFilters struct {
	Name   string
	Offset int
	Limit  int
}
