package utils

import (
	"errors"
	"fmt"
	"regexp"
)

func ValidateVersion(s string) error {
	if s == "" {
		return errors.New("missing required argument 'version'")
	}
	patterns := []string{
		`^\d{1,}\.\d{1,}$`,
		`^\d{1,}\.\d{1,}-build\.\d{1,}$`,
		`^\d{1,}\.\d{1,}\.\d{1,}$`,
		`^\d{1,}\.\d{1,}\.\d{1,}-build\.\d{1,}$`,
	}
	for _, pattern := range patterns {
		if regexp.MustCompile(pattern).MatchString(s) {
			return nil
		}
	}
	return fmt.Errorf("invalid version (%s) expected format [NUMBER].[NUMBER] or "+
		"[NUMBER].[NUMBER].[NUMBER]", s)
}
