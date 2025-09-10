//go:build unit

package bcrypthasher

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasherTestSuite struct {
	suite.Suite
}

func (s *BcryptHasherTestSuite) BeforeEach(t provider.T) {
	t.Epic("Authentication")
	t.Feature("Bcrypt Hasher")
}

func (s *BcryptHasherTestSuite) TestNewBcryptHasher(t provider.T) {
	t.Tags("creation", "validation")
	t.Description("Test creation of BcryptHasher with different cost parameters")

	t.Run("Valid cost parameter", func(t provider.T) {
		hasher := NewBcryptHasher(bcrypt.DefaultCost)
		assert.NotNil(t, hasher)
	})

	t.Run("Panic when cost too low", func(t provider.T) {
		assert.Panics(t, func() {
			NewBcryptHasher(bcrypt.MinCost - 1)
		}, "Should panic when cost is below minimum allowed")
	})

	t.Run("Panic when cost too high", func(t provider.T) {
		assert.Panics(t, func() {
			NewBcryptHasher(bcrypt.MaxCost + 1)
		}, "Should panic when cost is above maximum allowed")
	})
}

func (s *BcryptHasherTestSuite) TestBcryptHasherHash(t provider.T) {
	t.Tags("functionality", "hashing")
	t.Description("Test password hashing functionality")

	hasher := NewBcryptHasher(bcrypt.DefaultCost).(*bcrypthasher)

	t.Run("Successful password hashing", func(t provider.T) {
		password := []byte("test_password")

		var hashed []byte
		var err error
		t.WithNewStep("Hash password", func(ctx provider.StepCtx) {
			hashed, err = hasher.Hash(password)
		})

		t.WithNewStep("Verify hash result", func(ctx provider.StepCtx) {
			require.NoError(t, err)
			assert.NotEmpty(t, hashed)
			assert.NotEqual(t, password, hashed)
		})
	})

	t.Run("Hash empty password", func(t provider.T) {
		var hashed []byte
		var err error
		t.WithNewStep("Hash empty password", func(ctx provider.StepCtx) {
			hashed, err = hasher.Hash([]byte{})
		})

		t.WithNewStep("Verify empty password hash", func(ctx provider.StepCtx) {
			require.NoError(t, err)
			assert.NotEmpty(t, hashed)
		})
	})
}

func (s *BcryptHasherTestSuite) TestBcryptHasherCompare(t provider.T) {
	t.Tags("functionality", "verification")
	t.Description("Test password comparison functionality")

	hasher := NewBcryptHasher(bcrypt.DefaultCost).(*bcrypthasher)
	password := []byte("test_password")
	wrongPassword := []byte("wrong_password")

	var hashed []byte
	t.WithNewStep("Prepare hashed password", func(ctx provider.StepCtx) {
		var err error
		hashed, err = hasher.Hash(password)
		require.NoError(t, err)
	})

	t.Run("Successful password comparison", func(t provider.T) {
		var match bool
		var err error
		t.WithNewStep("Compare correct password", func(ctx provider.StepCtx) {
			match, err = hasher.Compare(hashed, password)
		})

		t.WithNewStep("Verify successful comparison", func(ctx provider.StepCtx) {
			require.NoError(t, err)
			assert.True(t, match)
		})
	})

	t.Run("Wrong password comparison", func(t provider.T) {
		var match bool
		var err error
		t.WithNewStep("Compare wrong password", func(ctx provider.StepCtx) {
			match, err = hasher.Compare(hashed, wrongPassword)
		})

		t.WithNewStep("Verify failed comparison", func(ctx provider.StepCtx) {
			require.Error(t, err)
			assert.False(t, match)
			assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err)
		})
	})

	t.Run("Empty password comparison", func(t provider.T) {
		var emptyHashed []byte
		t.WithNewStep("Hash empty password", func(ctx provider.StepCtx) {
			var err error
			emptyHashed, err = hasher.Hash([]byte{})
			require.NoError(t, err)
		})

		var match bool
		var err error
		t.WithNewStep("Compare empty password", func(ctx provider.StepCtx) {
			match, err = hasher.Compare(emptyHashed, []byte{})
		})

		t.WithNewStep("Verify empty password comparison", func(ctx provider.StepCtx) {
			require.NoError(t, err)
			assert.True(t, match)
		})
	})

	t.Run("Invalid hash comparison", func(t provider.T) {
		invalidHash := []byte("invalid_hash")

		var match bool
		var err error
		t.WithNewStep("Compare with invalid hash", func(ctx provider.StepCtx) {
			match, err = hasher.Compare(invalidHash, password)
		})

		t.WithNewStep("Verify invalid hash error", func(ctx provider.StepCtx) {
			require.Error(t, err)
			assert.False(t, match)
		})
	})
}

func TestBcryptHasher(t *testing.T) { suite.RunSuite(t, new(BcryptHasherTestSuite)) }
