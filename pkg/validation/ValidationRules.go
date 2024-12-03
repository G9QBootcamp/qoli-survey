package validation

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator"
)

// ValidateIranianNationalID validates the format and checksum of an Iranian national ID.
func ValidateIranianNationalID(fl validator.FieldLevel) bool {
	nationalID := fl.Field().String()

	match, _ := regexp.MatchString(`^\d{10}$`, nationalID)
	if !match {
		return false
	}

	// Checksum validation
	var checksum int
	for i := 0; i < 9; i++ {
		num, _ := strconv.Atoi(string(nationalID[i]))
		checksum += num * (10 - i)
	}
	controlDigit, _ := strconv.Atoi(string(nationalID[9]))
	calculatedDigit := checksum % 11

	if (calculatedDigit < 2 && calculatedDigit != controlDigit) || (calculatedDigit >= 2 && 11-calculatedDigit != controlDigit) {
		return false
	}

	return true
}

func RegisterCustomValidation(v *validator.Validate) {
	v.RegisterValidation("national_id", ValidateIranianNationalID)
}
