//go:build unit

package auth

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
)

type AuthValidationTestSuite struct {
	suite.Suite
}

func (s *AuthValidationTestSuite) BeforeEach(t provider.T) {
	t.Epic("Authentication")
	t.Feature("Email Validation")
}

func (s *AuthValidationTestSuite) TestValidateEmail(t provider.T) {
	t.Tags("validation", "email")
	t.Description("Test email validation with various input cases")

	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"Valid email address", "user@example.com", true},
		{"Valid email with subdomain", "user@sub.example.com", true},
		{"Valid email with plus addressing", "user+tag@example.com", true},
		{"Missing @ symbol", "userexample.com", false},
		{"Missing domain part", "user@", false},
		{"Invalid characters in domain", "user@exa mple.com", false},
		{"Empty email string", "", false},
		{"Email with multiple @ symbols", "user@ex@mple.com", false},
		{"Email with special characters", "user.name@example.com", true},
		{"Email with numeric domain", "user@123.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t provider.T) {
			t.WithNewStep("Prepare test parameters", func(ctx provider.StepCtx) {
				ctx.WithParameters(
					allure.NewParameter("Email", tt.email),
					allure.NewParameter("ExpectedValid", tt.valid),
				)
			})

			var err error
			t.WithNewStep("Validate email address", func(ctx provider.StepCtx) {
				err = validateEmail(tt.email)
			})

			t.WithNewStep("Verify validation result", func(ctx provider.StepCtx) {
				if tt.valid {
					assert.NoError(t, err, "Expected email '%s' to be valid", tt.email)
				} else {
					assert.Error(t, err, "Expected email '%s' to be invalid", tt.email)
				}
			})
		})
	}
}

func TestAuthValidation(t *testing.T) {
	suite.RunSuite(t, new(AuthValidationTestSuite))
}
