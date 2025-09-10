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

type AuthorTestSuite struct {
	suite.Suite
	validMember     Member
	validGiveTime   time.Time
	validRevokeTime time.Time
}

func (s *AuthorTestSuite) BeforeEach(t provider.T) {
	t.Epic("Authentication")
	t.Feature("Author Management")

	s.validMember = Member{
		id:         uuid.New(),
		name:       "Test Author",
		email:      "author@example.com",
		hashPasswd: []byte("$2a$10$hashedpassword"),
		createdAt:  time.Now().Add(-24 * time.Hour),
	}
	s.validGiveTime = time.Now().Add(-1 * time.Hour)
	s.validRevokeTime = time.Now().Add(1 * time.Hour)
}

func (s *AuthorTestSuite) TestCreateAuthorSuccess(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful creation of author with valid parameters")

	var author *Author
	var err error

	t.WithNewStep("Create author with valid data", func(ctx provider.StepCtx) {
		author, err = CreateAuthor(s.validMember, s.validGiveTime, false, s.validRevokeTime)
	})

	t.WithNewStep("Verify author properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, s.validMember.ID(), author.ID())
		assert.Equal(t, s.validGiveTime, author.giveTime)
		assert.False(t, author.HasAuthorRights())
		assert.True(t, author.HasMemberRights())
	})
}

func (s *AuthorTestSuite) TestCreateAuthorValidationErrors(t provider.T) {
	t.Tags("negative", "validation")
	t.Description("Test validation errors during author creation")

	testCases := []struct {
		name        string
		member      Member
		giveTime    time.Time
		revokeTime  time.Time
		expectError bool
	}{
		{
			"Invalid member",
			Member{},
			s.validGiveTime,
			s.validRevokeTime,
			true,
		},
		{
			"Zero give time",
			s.validMember,
			time.Time{},
			s.validRevokeTime,
			true,
		},
		{
			"Revoke time before give time",
			s.validMember,
			s.validGiveTime,
			s.validGiveTime.Add(-1 * time.Hour),
			false,
		},
		{
			"Valid data",
			s.validMember,
			s.validGiveTime,
			s.validRevokeTime,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t provider.T) {
			t.Tags("negative", "validation")

			t.WithNewStep("Prepare test parameters", func(ctx provider.StepCtx) {
				ctx.WithParameters(
					allure.NewParameter("GiveTime", tc.giveTime.String()),
					allure.NewParameter("RevokeTime", tc.revokeTime.String()),
					allure.NewParameter("ExpectError", tc.expectError),
				)
			})

			var err error
			t.WithNewStep("Attempt to create author", func(ctx provider.StepCtx) {
				_, err = CreateAuthor(tc.member, tc.giveTime, true, tc.revokeTime)
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

func (s *AuthorTestSuite) TestRevokeAuthorRights(t provider.T) {
	t.Tags("positive", "rights")
	t.Description("Test revocation of author rights")

	author := &Author{
		Member:     s.validMember,
		rights:     true,
		giveTime:   s.validGiveTime,
		revokeTime: time.Time{},
	}

	t.WithNewStep("Revoke author rights", func(ctx provider.StepCtx) {
		author.RevokeAuthorRights()
	})

	t.WithNewStep("Verify rights revocation", func(ctx provider.StepCtx) {
		assert.False(t, author.rights)
		assert.False(t, author.revokeTime.IsZero())
		assert.True(t, author.revokeTime.After(author.giveTime))
	})
}

func (s *AuthorTestSuite) TestHasAuthorRights(t provider.T) {
	t.Tags("functionality", "rights")
	t.Description("Test author rights checking")

	authorWithRights := &Author{rights: true}
	authorWithoutRights := &Author{rights: false}

	t.WithNewStep("Check author with rights", func(ctx provider.StepCtx) {
		assert.True(t, authorWithRights.HasAuthorRights())
	})

	t.WithNewStep("Check author without rights", func(ctx provider.StepCtx) {
		assert.False(t, authorWithoutRights.HasAuthorRights())
	})
}

func (s *AuthorTestSuite) TestHasMemberRights(t provider.T) {
	t.Tags("functionality", "rights")
	t.Description("Test member rights checking")

	author := &Author{}

	t.WithNewStep("Check member rights", func(ctx provider.StepCtx) {
		assert.True(t, author.HasMemberRights())
	})
}

func (s *AuthorTestSuite) TestAuth(t provider.T) {
	t.Tags("functionality", "authentication")
	t.Description("Test author authentication")

	author := &Author{Member: s.validMember}
	authFunc := func(hash, plain []byte) (bool, error) {
		return true, nil
	}

	var result bool
	t.WithNewStep("Authenticate author", func(ctx provider.StepCtx) {
		result = author.Auth([]byte("password"), authFunc)
	})

	t.WithNewStep("Verify authentication success", func(ctx provider.StepCtx) {
		assert.True(t, result)
	})
}

func (s *AuthorTestSuite) TestID(t provider.T) {
	t.Tags("functionality", "identification")
	t.Description("Test author ID retrieval")

	author := &Author{Member: s.validMember}

	t.WithNewStep("Get author ID", func(ctx provider.StepCtx) {
		assert.Equal(t, s.validMember.ID(), author.ID())
	})
}

func TestAuthor(t *testing.T) {
	suite.RunSuite(t, new(AuthorTestSuite))
}
