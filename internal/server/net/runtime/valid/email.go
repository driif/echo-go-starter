package valid

import (
	"errors"
	"regexp"
)

func Email(email string) error {
	if len(email) < 3 && len(email) > 254 {
		return errors.New("email must be between 3 and 254 characters")
	}
	regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !regex.MatchString(email) {
		return errors.New("email can not be validated. Please check the format")
	}
	return nil
}
