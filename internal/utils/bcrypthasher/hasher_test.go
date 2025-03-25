package bcrypthasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestNewBcryptHasher(t *testing.T) {
	t.Run("valid cost", func(t *testing.T) {
		hasher := NewBcryptHasher(bcrypt.DefaultCost)
		assert.NotNil(t, hasher)
	})

	t.Run("panic when cost too low", func(t *testing.T) {
		assert.Panics(t, func() {
			NewBcryptHasher(bcrypt.MinCost - 1)
		}, "should panic when cost is below MinCost")
	})

	t.Run("panic when cost too high", func(t *testing.T) {
		assert.Panics(t, func() {
			NewBcryptHasher(bcrypt.MaxCost + 1)
		}, "should panic when cost is above MaxCost")
	})
}

func TestBcryptHasher_Hash(t *testing.T) {
	hasher := NewBcryptHasher(bcrypt.DefaultCost).(*bcrypthasher)

	t.Run("successful hash", func(t *testing.T) {
		password := []byte("test_password")
		hashed, err := hasher.Hash(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hashed)
		assert.NotEqual(t, password, hashed)
	})

	t.Run("empty password", func(t *testing.T) {
		hashed, err := hasher.Hash([]byte{})

		require.NoError(t, err)
		assert.NotEmpty(t, hashed)
	})
}

func TestBcryptHasher_Compare(t *testing.T) {
	hasher := NewBcryptHasher(bcrypt.DefaultCost).(*bcrypthasher)
	password := []byte("test_password")
	wrongPassword := []byte("wrong_password")
	hashed, _ := hasher.Hash(password)

	t.Run("successful compare", func(t *testing.T) {
		match, err := hasher.Compare(hashed, password)

		require.NoError(t, err)
		assert.True(t, match)
	})

	t.Run("wrong password", func(t *testing.T) {
		match, err := hasher.Compare(hashed, wrongPassword)

		require.Error(t, err)
		assert.False(t, match)
		assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err)
	})

	t.Run("empty password", func(t *testing.T) {
		emptyHashed, _ := hasher.Hash([]byte{})
		match, err := hasher.Compare(emptyHashed, []byte{})

		require.NoError(t, err)
		assert.True(t, match)
	})

	t.Run("invalid hash", func(t *testing.T) {
		invalidHash := []byte("invalid_hash")
		match, err := hasher.Compare(invalidHash, password)

		require.Error(t, err)
		assert.False(t, match)
	})
}
