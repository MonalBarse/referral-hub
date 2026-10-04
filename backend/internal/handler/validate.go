package handler

import (
	"net/mail"
	"strings"
)

type fieldErrors map[string]string

func (fe fieldErrors) ok() bool { return len(fe) == 0 }

func (fe fieldErrors) required(field, value string, max int) string {
	trimmed := strings.TrimSpace(value)
	switch {
	case trimmed == "":
		fe[field] = "is required"
	case len(trimmed) > max:
		fe[field] = "is too long"
	}
	return trimmed
}

func (fe fieldErrors) optional(field, value string, max int) string {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) > max {
		fe[field] = "is too long"
	}
	return trimmed
}

func (fe fieldErrors) email(field, value string, max int) string {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	if trimmed == "" {
		fe[field] = "is required"
		return trimmed
	}
	if len(trimmed) > max {
		fe[field] = "is too long"
		return trimmed
	}
	if _, err := mail.ParseAddress(trimmed); err != nil {
		fe[field] = "must be a valid email address"
	}
	return trimmed
}

func (fe fieldErrors) oneOf(field, value, fallback string, allowed ...string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}
	for _, a := range allowed {
		if trimmed == a {
			return trimmed
		}
	}
	fe[field] = "is not a valid value"
	return trimmed
}
