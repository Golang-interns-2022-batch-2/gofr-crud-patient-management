package service

import (
	"regexp"

	"gopkg.in/guregu/null.v4"
)

func validateName(name null.String) bool {
	if name.IsZero() || name.String == "" {
		return false
	}

	return true
}

func validatePhone(phone null.String) bool {
	phoneRegex := regexp.MustCompile(`^\s*(?:\+?(\d{1,3}))?[-. (]*(\d{3})[-. )]*(\d{3})[-. ]*(\d{4})(?: *x(\d+))?\s*$`)
	return phoneRegex.MatchString(phone.String)
}

func validateBloodGroup(bloodgroup null.String) bool {
	switch bloodgroup.String {
	case "+A", "-A", "+B", "-B", "+O", "-O", "+AB", "-AB":
		return true
	default:
		return false
	}
}
