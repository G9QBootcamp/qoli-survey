package dto

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required"`
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
