package validation

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator"
)

func IsValidNationalID(fl validator.FieldLevel) bool {
	nationalID := fl.Field().String()
	if len(nationalID) != 10 {
		return false
	}

	matched, _ := regexp.MatchString(`^\d{10}$`, nationalID)
	if !matched {
		return false
	}

	allSame := true
	for i := 1; i < len(nationalID); i++ {
		if nationalID[i] != nationalID[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}

	sum := 0
	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(nationalID[i]))
		sum += digit * (10 - i)
	}
	checksum := sum % 11
	lastDigit, _ := strconv.Atoi(string(nationalID[9]))
	return (checksum < 2 && checksum == lastDigit) || (checksum >= 2 && lastDigit == 11-checksum)
}

func RegisterCustomValidation(v *validator.Validate) {
	v.RegisterValidation("national_id", IsValidNationalID)
}
