package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMember(t *testing.T) {
	validID := uuid.New()
	validName := "John Doe"
	validEmail := "john@example.com"
	validHash := []byte("$2a$10$hashedpassword")
	validTime := time.Now()

	t.Run("CreateMember - успешное создание", func(t *testing.T) {
		member, err := CreateMember(validID, validName, validEmail, validHash, validTime)

		require.NoError(t, err)
		assert.Equal(t, validID, member.ID())
		assert.Equal(t, validName, member.name)
		assert.Equal(t, validEmail, member.email)
		assert.Equal(t, validHash, member.hashPasswd)
		assert.Equal(t, validTime, member.createdAt)
	})

	t.Run("CreateMember - ошибки валидации", func(t *testing.T) {
		testCases := []struct {
			name        string
			id          uuid.UUID
			nameStr     string
			email       string
			hashPasswd  []byte
			createdAt   time.Time
			expectError bool
		}{
			{"Пустое имя", validID, "", validEmail, validHash, validTime, true},
			{"Пустой email", validID, validName, "", validHash, validTime, true},
			{"Невалидный email", validID, validName, "invalid-email", validHash, validTime, true},
			{"Nil хеш пароля", validID, validName, validEmail, nil, validTime, true},
			{"Zero время создания", validID, validName, validEmail, validHash, time.Time{}, true},
			{"Корректные данные", validID, validName, validEmail, validHash, validTime, false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := CreateMember(tc.id, tc.nameStr, tc.email, tc.hashPasswd, tc.createdAt)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("NewMember - создание с генерацией ID и времени", func(t *testing.T) {
		member, err := NewMember(validName, validEmail, validHash)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, member.ID())
		assert.False(t, member.createdAt.IsZero())
	})

	t.Run("Проверка прав", func(t *testing.T) {
		member := &Member{}
		assert.False(t, member.HasAuthorRights())
		assert.True(t, member.HasMemberRights())
	})

	t.Run("Auth - успешная аутентификация", func(t *testing.T) {
		member := &Member{hashPasswd: validHash}
		authFunc := func(hash, plain []byte) (bool, error) {
			return true, nil
		}

		assert.True(t, member.Auth([]byte("password"), authFunc))
	})

	t.Run("Auth - неудачная аутентификация", func(t *testing.T) {
		member := &Member{hashPasswd: validHash}
		authFunc := func(hash, plain []byte) (bool, error) {
			return false, nil
		}

		assert.False(t, member.Auth([]byte("wrong"), authFunc))
	})

	t.Run("Auth - ошибка при аутентификации", func(t *testing.T) {
		member := &Member{hashPasswd: validHash}
		authFunc := func(hash, plain []byte) (bool, error) {
			return false, assert.AnError
		}

		assert.False(t, member.Auth([]byte("password"), authFunc))
	})
}
