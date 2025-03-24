package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"Valid email", "user@example.com", true},
		{"Valid with subdomain", "user@sub.example.com", true},
		{"Valid with plus", "user+tag@example.com", true},
		{"Missing @", "userexample.com", false},
		{"Missing domain", "user@", false},
		{"Invalid chars", "user@exa mple.com", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
