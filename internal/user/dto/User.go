package dto

type UserResponse struct {
	ID          uint   `json:"id"`
	NationalID  string `json:"national_id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	City        string `json:"city"`
	DateOfBirth string `json:"date_of_birth"`
}

type UserGetRequest struct {
	Page int
	Name string
}

type UserFilters struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	NationalID  string `json:"national_id"`
	City        string `json:"city"`
	YearOfBirth int    `json:"year_of_birth"`
	Offset      int
	Limit       int
}

type SignupRequest struct {
	NationalID  string `json:"national_id" validate:"required,national_id"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	City        string `json:"city"`
	DateOfBirth string `json:"date_of_birth" validate:"omitempty,date"`
}

type UpdateUserRequest struct {
	FirstName   string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName    string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
	DateOfBirth string `json:"date_of_birth,omitempty" validate:"omitempty,date"`
	City        string `json:"city,omitempty" validate:"omitempty,min=1,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
