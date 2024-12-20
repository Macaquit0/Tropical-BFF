package helper

import (
	"regexp"
)

func ValidatePassword(s string) bool {
	if len(s) < 8 {
		return false
	}

	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString
	hasSpecial := regexp.MustCompile(`[!@#~$%^&*(),.?:{}|<>]`).MatchString

	if !hasLetter(s) {
		return false
	}

	if !hasNumber(s) {
		return false
	}

	if !hasSpecial(s) {
		return false
	}

	return true
}
