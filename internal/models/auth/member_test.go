//go:build unit

package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MemberTestSuite struct {
	suite.Suite
	validID    uuid.UUID
	validName  string
	validEmail string
	validHash  []byte
	validTime  time.Time
}

func (s *MemberTestSuite) BeforeEach(t provider.T) {
	t.Epic("Authentication")
	t.Feature("Member Management")

	s.validID = uuid.New()
	s.validName = "John Doe"
	s.validEmail = "john@example.com"
	s.validHash = []byte("$2a$10$hashedpassword")
	s.validTime = time.Now()
}

func (s *MemberTestSuite) TestCreateMemberSuccess(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful creation of member with valid parameters")

	var member *Member
	var err error

	t.WithNewStep("Create member with valid data", func(ctx provider.StepCtx) {
		member, err = CreateMember(s.validID, s.validName, s.validEmail, s.validHash, s.validTime)
	})

	t.WithNewStep("Verify member properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, s.validID, member.ID())
		assert.Equal(t, s.validName, member.name)
		assert.Equal(t, s.validEmail, member.email)
		assert.Equal(t, s.validHash, member.hashPasswd)
		assert.Equal(t, s.validTime, member.createdAt)
	})
}

func (s *MemberTestSuite) TestCreateMemberValidationErrors(t provider.T) {
	t.Tags("negative", "validation")
	t.Description("Test validation errors during member creation")

	testCases := []struct {
		name        string
		id          uuid.UUID
		nameStr     string
		email       string
		hashPasswd  []byte
		createdAt   time.Time
		expectError bool
	}{
		{"Empty name", s.validID, "", s.validEmail, s.validHash, s.validTime, true},
		{"Empty email", s.validID, s.validName, "", s.validHash, s.validTime, true},
		{"Invalid email", s.validID, s.validName, "invalid-email", s.validHash, s.validTime, true},
		{"Nil password hash", s.validID, s.validName, s.validEmail, nil, s.validTime, true},
		{"Zero creation time", s.validID, s.validName, s.validEmail, s.validHash, time.Time{}, true},
		{"Valid data", s.validID, s.validName, s.validEmail, s.validHash, s.validTime, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t provider.T) {
			t.Tags("negative", "validation")

			t.WithNewStep("Prepare test parameters", func(ctx provider.StepCtx) {
				ctx.WithParameters(
					allure.NewParameter("Name", tc.nameStr),
					allure.NewParameter("Email", tc.email),
					allure.NewParameter("ExpectError", tc.expectError),
				)
			})

			var err error
			t.WithNewStep("Attempt to create member", func(ctx provider.StepCtx) {
				_, err = CreateMember(tc.id, tc.nameStr, tc.email, tc.hashPasswd, tc.createdAt)
			})

			t.WithNewStep("Verify validation result", func(ctx provider.StepCtx) {
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		})
	}
}

func (s *MemberTestSuite) TestNewMemberSuccess(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful creation of member with auto-generated fields")

	var member *Member
	var err error

	t.WithNewStep("Create member with NewMember", func(ctx provider.StepCtx) {
		member, err = NewMember(s.validName, s.validEmail, s.validHash)
	})

	t.WithNewStep("Verify auto-generated properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, member.ID())
		assert.False(t, member.createdAt.IsZero())
		assert.Equal(t, s.validName, member.name)
		assert.Equal(t, s.validEmail, member.email)
		assert.Equal(t, s.validHash, member.hashPasswd)
	})
}

func (s *MemberTestSuite) TestPermissionChecks(t provider.T) {
	t.Tags("functionality", "permissions")
	t.Description("Test member permission checks")

	member := &Member{}

	t.WithNewStep("Check author rights", func(ctx provider.StepCtx) {
		assert.False(t, member.HasAuthorRights())
	})

	t.WithNewStep("Check member rights", func(ctx provider.StepCtx) {
		assert.True(t, member.HasMemberRights())
	})
}

func (s *MemberTestSuite) TestAuthSuccess(t provider.T) {
	t.Tags("positive", "authentication")
	t.Description("Successful member authentication")

	member := &Member{hashPasswd: s.validHash}
	authFunc := func(hash, plain []byte) (bool, error) {
		return true, nil
	}

	var result bool
	t.WithNewStep("Authenticate member", func(ctx provider.StepCtx) {
		result = member.Auth([]byte("password"), authFunc)
	})

	t.WithNewStep("Verify authentication success", func(ctx provider.StepCtx) {
		assert.True(t, result)
	})
}

func (s *MemberTestSuite) TestAuthFailure(t provider.T) {
	t.Tags("negative", "authentication")
	t.Description("Failed member authentication")

	member := &Member{hashPasswd: s.validHash}
	authFunc := func(hash, plain []byte) (bool, error) {
		return false, nil
	}

	var result bool
	t.WithNewStep("Attempt authentication with wrong password", func(ctx provider.StepCtx) {
		result = member.Auth([]byte("wrong"), authFunc)
	})

	t.WithNewStep("Verify authentication failure", func(ctx provider.StepCtx) {
		assert.False(t, result)
	})
}

func (s *MemberTestSuite) TestAuthError(t provider.T) {
	t.Tags("negative", "authentication")
	t.Description("Authentication error handling")

	member := &Member{hashPasswd: s.validHash}
	authFunc := func(hash, plain []byte) (bool, error) {
		return false, assert.AnError
	}

	var result bool
	t.WithNewStep("Attempt authentication with error", func(ctx provider.StepCtx) {
		result = member.Auth([]byte("password"), authFunc)
	})

	t.WithNewStep("Verify authentication handles error", func(ctx provider.StepCtx) {
		assert.False(t, result)
	})
}

func TestMember(t *testing.T) {
	suite.RunSuite(t, new(MemberTestSuite))
}
