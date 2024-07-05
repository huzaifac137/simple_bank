package grpc_util

import (
	"fmt"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidEmail    = regexp.MustCompile(`^[a-z0-9_.%+-]+@[a-z0-9_.-]+\.[a-z]{2,}$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLen int, maxLen int) error {
	if len(value) < minLen || len(value) > maxLen {
		return fmt.Errorf("value must be  %d - %d characters long", minLen, maxLen)
	}
	return nil
}

func ValidateUsername(value string) error {
	err := ValidateString(value, 3, 100)
	if err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("username must only contain lowercase alphanumeric characters and underscores")
	}

	return nil

}

func ValidatePassword(value string) error {
	err := ValidateString(value, 6, 100)
	if err != nil {
		return err
	}

	return nil
}

func ValidateEmail(value string) error {
	err := ValidateString(value, 3, 100)
	if err != nil {
		return err
	}

	if !isValidEmail(value) {
		err = fmt.Errorf("invalid email address")
		return err
	}

	return nil
}

func ValidateFullName(value string) error {
	err := ValidateString(value, 3, 100)
	if err != nil {
		return err
	}

	if !isValidFullName(value) {
		err = fmt.Errorf("full name must contain only letters and spaces")
		return err
	}

	return nil
}
