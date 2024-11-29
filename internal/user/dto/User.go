package dto

type UserResponse struct {
	Name string `json:"name"`
}

type UserRequest struct {
	Name string
}
