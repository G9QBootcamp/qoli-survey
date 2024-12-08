package validation

import (
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator"
)

func ValidateIranianNationalID(fl validator.FieldLevel) bool {
	nationalID := fl.Field().String()
	if len(nationalID) != 10 {
		return false
	}

	match, _ := regexp.MatchString(`^\d{10}$`, nationalID)
	if !match {
		return false
	}

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

func ValidateDateFormat(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()

	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

func RegisterCustomValidation(v *validator.Validate) {
	v.RegisterValidation("national_id", ValidateIranianNationalID)
	v.RegisterValidation("date", ValidateDateFormat)

}
