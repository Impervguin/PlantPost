package album

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlbum(t *testing.T) {
	// Подготовка тестовых данных
	validID := uuid.New()
	validName := "Test Album"
	validDescription := "Test Description"
	validPlantIDs := uuid.UUIDs{uuid.New(), uuid.New()}
	validOwnerID := uuid.New()
	validCreatedAt := time.Now().Add(-time.Hour)
	validUpdatedAt := time.Now()

	t.Run("CreateAlbum - успешное создание", func(t *testing.T) {
		album, err := CreateAlbum(
			validID,
			validName,
			validDescription,
			validPlantIDs,
			validOwnerID,
			validCreatedAt,
			validUpdatedAt,
		)

		require.NoError(t, err)
		assert.Equal(t, validID, album.id)
		assert.Equal(t, validName, album.name)
		assert.Equal(t, validDescription, album.description)
		assert.Equal(t, validPlantIDs, album.plantIDs)
		assert.Equal(t, validOwnerID, album.ownerID)
		assert.Equal(t, validCreatedAt, album.createdAt)
	})

	t.Run("CreateAlbum - ошибки валидации", func(t *testing.T) {
		testCases := []struct {
			name        string
			id          uuid.UUID
			nameStr     string
			description string
			plantIDs    uuid.UUIDs
			ownerID     uuid.UUID
			createdAt   time.Time
			updatedAt   time.Time
			expectError bool
		}{
			{
				"Пустой ID",
				uuid.Nil,
				validName,
				validDescription,
				validPlantIDs,
				validOwnerID,
				validCreatedAt,
				validUpdatedAt,
				true,
			},
			{
				"Пустое название",
				validID,
				"",
				validDescription,
				validPlantIDs,
				validOwnerID,
				validCreatedAt,
				validUpdatedAt,
				true,
			},
			{
				"Пустой ownerID",
				validID,
				validName,
				validDescription,
				validPlantIDs,
				uuid.Nil,
				validCreatedAt,
				validUpdatedAt,
				true,
			},
			{
				"Nil plantIDs",
				validID,
				validName,
				validDescription,
				nil,
				validOwnerID,
				validCreatedAt,
				validUpdatedAt,
				true,
			},
			{
				"Пустой plantID в списке",
				validID,
				validName,
				validDescription,
				uuid.UUIDs{uuid.New(), uuid.Nil},
				validOwnerID,
				validCreatedAt,
				validUpdatedAt,
				true,
			},
			{
				"Обновление в будущем",
				validID,
				validName,
				validDescription,
				validPlantIDs,
				validOwnerID,
				validCreatedAt,
				time.Now().Add(time.Hour),
				true,
			},
			{
				"Создание после обновления",
				validID,
				validName,
				validDescription,
				validPlantIDs,
				validOwnerID,
				validUpdatedAt,
				validCreatedAt,
				true,
			},
			{
				"Корректные данные",
				validID,
				validName,
				validDescription,
				validPlantIDs,
				validOwnerID,
				validCreatedAt,
				validUpdatedAt,
				false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := CreateAlbum(
					tc.id,
					tc.nameStr,
					tc.description,
					tc.plantIDs,
					tc.ownerID,
					tc.createdAt,
					tc.updatedAt,
				)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("NewAlbum - успешное создание", func(t *testing.T) {
		album, err := NewAlbum(
			validName,
			validDescription,
			validPlantIDs,
			validOwnerID,
		)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, album.id)
		assert.True(t, album.createdAt.Before(time.Now()) || album.createdAt.Equal(time.Now()))
		assert.True(t, album.updatedAt.Before(time.Now()) || album.updatedAt.Equal(time.Now()))
	})

	t.Run("GetOwnerID", func(t *testing.T) {
		album := &Album{ownerID: validOwnerID}
		assert.Equal(t, validOwnerID, album.GetOwnerID())
	})

	t.Run("UpdateName", func(t *testing.T) {
		album := &Album{name: "Old Name", updatedAt: validUpdatedAt}
		newName := "New Name"

		err := album.UpdateName(newName)
		require.NoError(t, err)
		assert.Equal(t, newName, album.name)
		assert.True(t, album.updatedAt.After(validUpdatedAt))
	})

	t.Run("UpdateDescription", func(t *testing.T) {
		album := &Album{description: "Old Desc", updatedAt: validUpdatedAt}
		newDesc := "New Desc"

		err := album.UpdateDescription(newDesc)
		require.NoError(t, err)
		assert.Equal(t, newDesc, album.description)
		assert.True(t, album.updatedAt.After(validUpdatedAt))
	})

	t.Run("AddPlant", func(t *testing.T) {
		album := &Album{plantIDs: validPlantIDs, updatedAt: validUpdatedAt}
		newPlantID := uuid.New()

		err := album.AddPlant(newPlantID)
		require.NoError(t, err)
		assert.Contains(t, album.plantIDs, newPlantID)
		assert.True(t, album.updatedAt.After(validUpdatedAt))
	})

	t.Run("RemovePlant - успешное удаление", func(t *testing.T) {
		plantToRemove := validPlantIDs[0]
		album := &Album{plantIDs: validPlantIDs, updatedAt: validUpdatedAt}

		err := album.RemovePlant(plantToRemove)
		require.NoError(t, err)
		assert.NotContains(t, album.plantIDs, plantToRemove)
		assert.True(t, album.updatedAt.After(validUpdatedAt))
	})

	t.Run("RemovePlant - растение не найдено", func(t *testing.T) {
		album := &Album{plantIDs: validPlantIDs, updatedAt: validUpdatedAt}
		nonExistentPlant := uuid.New()

		err := album.RemovePlant(nonExistentPlant)
		assert.Error(t, err)
		assert.Equal(t, validPlantIDs, album.plantIDs)
	})
}
