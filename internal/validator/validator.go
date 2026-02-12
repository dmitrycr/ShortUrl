package validator

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrInvalidURL    = errors.New("invalid URL format")
	ErrURLTooLong    = errors.New("URL exceeds maximum length")
	ErrInvalidScheme = errors.New("URL must use http or https scheme")
	ErrInvalidCode   = errors.New("invalid short code format")
	ErrCodeTooLong   = errors.New("short code exceeds maximum length")
)

const (
	MaxURLLength  = 2048
	MaxCodeLength = 20
	MinCodeLength = 3
)

type URLValidator struct{}

func NewURLValidator() *URLValidator {
	return &URLValidator{}
}

// ValidateURL проверяет корректность URL
func (v *URLValidator) ValidateURL(rawURL string) error {
	// Проверка длины
	if len(rawURL) > MaxURLLength {
		return ErrURLTooLong
	}

	if len(rawURL) == 0 {
		return ErrInvalidURL
	}

	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return ErrInvalidURL
	}
	return nil
}

func (v *URLValidator) ValidateCustomCode(code string) error {
	if len(code) == 0 {
		return ErrCodeTooLong
	}

	if len(code) < MinCodeLength {
		return ErrInvalidCode
	}

	if len(code) > MaxCodeLength {
		return ErrInvalidCode
	}

	for _, char := range code {
		if !isValidCodeChar(char) {
			return ErrInvalidCode
		}
	}

	return nil
}

func (v *URLValidator) NormalizeURL(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return "https://" + rawURL
	}
	return rawURL
}

// isValidCodeChar проверяет, допустим ли символ в коде
func isValidCodeChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '-' || char == '_'
}
