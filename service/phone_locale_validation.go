package service

import (
	"unicode"

	"github.com/supermarios-hotel-management-system/authentication/model"
)

type phoneLocaleValidation struct {
	minLen int
	maxLen int
}

var phoneLocaleValidationData = map[model.UserDialCode]phoneLocaleValidation{
	model.IndiaDialCode: {minLen: 10, maxLen: 10},
}

// ValidatePhone validates the phone number based on the dial code.
func ValidatePhone(dialCode model.UserDialCode, phone string) bool {
	validations, ok := phoneLocaleValidationData[dialCode]
	if !ok {
		return false
	}

	if len(phone) < validations.minLen || len(phone) > validations.maxLen {
		return false
	}

	for _, r := range phone {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}
