package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthor(t *testing.T) {
	// Подготовка тестовых данных
	validMember := Member{
		id:         uuid.New(),
		name:       "Test Author",
		email:      "author@example.com",
		hashPasswd: []byte("$2a$10$hashedpassword"),
		createdAt:  time.Now().Add(-24 * time.Hour),
	}
	validGiveTime := time.Now().Add(-1 * time.Hour)
	validRevokeTime := time.Now().Add(1 * time.Hour)

	t.Run("CreateAuthor - успешное создание", func(t *testing.T) {
		author, err := CreateAuthor(validMember, validGiveTime, true, validRevokeTime)

		require.NoError(t, err)
		assert.Equal(t, validMember.ID(), author.ID())
		assert.Equal(t, validGiveTime, author.giveTime)
		assert.True(t, author.HasAuthorRights())
		assert.True(t, author.HasMemberRights())
	})

	t.Run("CreateAuthor - ошибки валидации", func(t *testing.T) {
		testCases := []struct {
			name        string
			member      Member
			giveTime    time.Time
			revokeTime  time.Time
			expectError bool
		}{
			{
				"Невалидный member",
				Member{},
				validGiveTime,
				validRevokeTime,
				true,
			},
			{
				"Zero giveTime",
				validMember,
				time.Time{},
				validRevokeTime,
				true,
			},
			{
				"RevokeTime before giveTime",
				validMember,
				validGiveTime,
				validGiveTime.Add(-1 * time.Hour),
				true,
			},
			{
				"Корректные данные",
				validMember,
				validGiveTime,
				validRevokeTime,
				false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := CreateAuthor(tc.member, tc.giveTime, true, tc.revokeTime)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("RevokeAuthorRights", func(t *testing.T) {
		author := &Author{
			Member:     validMember,
			rights:     true,
			giveTime:   validGiveTime,
			revokeTime: time.Time{},
		}

		author.RevokeAuthorRights()
		assert.False(t, author.rights)
		assert.False(t, author.revokeTime.IsZero())
		assert.True(t, author.revokeTime.After(author.giveTime))
	})

	t.Run("HasAuthorRights", func(t *testing.T) {
		authorWithRights := &Author{rights: true}
		authorWithoutRights := &Author{rights: false}

		assert.True(t, authorWithRights.HasAuthorRights())
		assert.False(t, authorWithoutRights.HasAuthorRights())
	})

	t.Run("HasMemberRights", func(t *testing.T) {
		author := &Author{}
		assert.True(t, author.HasMemberRights())
	})

	t.Run("Auth", func(t *testing.T) {
		author := &Author{Member: validMember}
		authFunc := func(hash, plain []byte) (bool, error) {
			return true, nil
		}

		assert.True(t, author.Auth([]byte("password"), authFunc))
	})

	t.Run("ID", func(t *testing.T) {
		author := &Author{Member: validMember}
		assert.Equal(t, validMember.ID(), author.ID())
	})
}
