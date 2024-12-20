package helper

import (
	"regexp"
)

func IsValidZipcode(zipcode string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(zipcode)
}
