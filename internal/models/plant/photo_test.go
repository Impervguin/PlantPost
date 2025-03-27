package plant

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlantPhoto(t *testing.T) {
	validID := uuid.New()
	validFileID := uuid.New()
	validDescription := "Test description"

	t.Run("CreatePlantPhoto - успешное создание", func(t *testing.T) {
		photo, err := CreatePlantPhoto(validID, validFileID, validDescription)
		require.NoError(t, err)
		assert.Equal(t, validID, photo.id)
		assert.Equal(t, validFileID, photo.fileID)
		assert.Equal(t, validDescription, photo.description)
	})

	t.Run("CreatePlantPhoto - ошибки валидации", func(t *testing.T) {
		testCases := []struct {
			name        string
			id          uuid.UUID
			fileID      uuid.UUID
			expectedErr string
		}{
			{
				name:        "empty ID",
				id:          uuid.Nil,
				fileID:      validFileID,
				expectedErr: "plant photo ID cannot be empty",
			},
			{
				name:        "empty file ID",
				id:          validID,
				fileID:      uuid.Nil,
				expectedErr: "plant photo file ID cannot be empty",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := CreatePlantPhoto(tc.id, tc.fileID, validDescription)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
			})
		}
	})

	t.Run("NewPlantPhoto - успешное создание", func(t *testing.T) {
		photo, err := NewPlantPhoto(validFileID, validDescription)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, photo.id)
		assert.Equal(t, validFileID, photo.fileID)
		assert.Equal(t, validDescription, photo.description)
	})

	t.Run("Compare - сравнение фото", func(t *testing.T) {
		photo1, _ := CreatePlantPhoto(validID, validFileID, "Photo 1")
		photo2, _ := CreatePlantPhoto(uuid.New(), uuid.New(), "Photo 2")
		photo3, _ := CreatePlantPhoto(validID, uuid.New(), "Photo 3")
		photo4, _ := CreatePlantPhoto(uuid.New(), validFileID, "Photo 4")

		assert.True(t, photo1.Compare(photo1))  // Тот же объект
		assert.True(t, photo1.Compare(photo3))  // Тот же ID
		assert.True(t, photo1.Compare(photo4))  // Тот же fileID
		assert.False(t, photo1.Compare(photo2)) // Разные ID и fileID
		assert.False(t, photo1.Compare(nil))    // nil сравнение
	})

	t.Run("Getters", func(t *testing.T) {
		photo, _ := CreatePlantPhoto(validID, validFileID, validDescription)
		assert.Equal(t, validID, photo.ID())
		assert.Equal(t, validFileID, photo.FileID())
		assert.Equal(t, validDescription, photo.Description())
	})

	t.Run("Validate", func(t *testing.T) {
		validPhoto := &PlantPhoto{id: validID, fileID: validFileID}
		assert.NoError(t, validPhoto.Validate())

		invalidPhoto := &PlantPhoto{id: uuid.Nil, fileID: validFileID}
		assert.Error(t, invalidPhoto.Validate())
	})
}

func TestPlantPhotos(t *testing.T) {
	photo1, _ := CreatePlantPhoto(uuid.New(), uuid.New(), "Photo 1")
	photo2, _ := CreatePlantPhoto(uuid.New(), uuid.New(), "Photo 2")

	t.Run("NewPlantPhotos", func(t *testing.T) {
		pp := NewPlantPhotos()
		assert.NotNil(t, pp)
		assert.Equal(t, 0, pp.Len())
	})

	t.Run("Add - успешное добавление", func(t *testing.T) {
		pp := NewPlantPhotos()
		require.NoError(t, pp.Add(photo1))
		assert.Equal(t, 1, pp.Len())
		require.NoError(t, pp.Add(photo2))
		assert.Equal(t, 2, pp.Len())
	})

	t.Run("Add - дубликат фото", func(t *testing.T) {
		pp := NewPlantPhotos()
		require.NoError(t, pp.Add(photo1))
		err := pp.Add(photo1)
		require.Error(t, err)
		assert.Equal(t, "photo already exists", err.Error())
		assert.Equal(t, 1, pp.Len())
	})

	t.Run("Remove - успешное удаление", func(t *testing.T) {
		pp := NewPlantPhotos()
		require.NoError(t, pp.Add(photo1))
		require.NoError(t, pp.Add(photo2))
		require.NoError(t, pp.Remove(photo1))
		assert.Equal(t, 1, pp.Len())
	})

	t.Run("Remove - фото не найдено", func(t *testing.T) {
		pp := NewPlantPhotos()
		require.NoError(t, pp.Add(photo1))
		err := pp.Remove(photo2)
		require.Error(t, err)
		assert.Equal(t, "photo not found", err.Error())
		assert.Equal(t, 1, pp.Len())
	})

	t.Run("Iterate", func(t *testing.T) {
		pp := NewPlantPhotos()
		require.NoError(t, pp.Add(photo1))
		require.NoError(t, pp.Add(photo2))

		var collected []PlantPhoto
		err := pp.Iterate(func(e PlantPhoto) error {
			collected = append(collected, e)
			return nil
		})
		require.NoError(t, err)
		assert.Len(t, collected, 2)
		assert.Contains(t, collected, *photo1)
		assert.Contains(t, collected, *photo2)
	})

	t.Run("Iterate - с ошибкой", func(t *testing.T) {
		pp := NewPlantPhotos()
		require.NoError(t, pp.Add(photo1))
		require.NoError(t, pp.Add(photo2))

		testErr := errors.New("test error")
		err := pp.Iterate(func(e PlantPhoto) error {
			if e.Compare(photo2) {
				return testErr
			}
			return nil
		})
		require.Error(t, err)
		assert.Equal(t, testErr, err)
	})

	t.Run("IterateUpdate - успешное обновление", func(t *testing.T) {
		pp := NewPlantPhotos()
		require.NoError(t, pp.Add(photo1))
		require.NoError(t, pp.Add(photo2))

		newDesc := "Updated description"
		err := pp.IterateUpdate(func(e *PlantPhoto) error {
			if e.Compare(photo2) {
				e.description = newDesc
			}
			return nil
		})
		require.NoError(t, err)

		var found bool
		_ = pp.Iterate(func(e PlantPhoto) error {
			if e.Compare(photo2) {
				assert.Equal(t, newDesc, e.description)
				found = true
			}
			return nil
		})
		assert.True(t, found)
	})

	t.Run("IterateUpdate - с ошибкой валидации", func(t *testing.T) {
		pp := NewPlantPhotos()
		require.NoError(t, pp.Add(photo1))

		err := pp.IterateUpdate(func(e *PlantPhoto) error {
			e.fileID = uuid.Nil // Делаем невалидным
			return nil
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "plant photo file ID cannot be empty")
	})

	t.Run("Len", func(t *testing.T) {
		pp := NewPlantPhotos()
		assert.Equal(t, 0, pp.Len())
		require.NoError(t, pp.Add(photo1))
		assert.Equal(t, 1, pp.Len())
		require.NoError(t, pp.Add(photo2))
		assert.Equal(t, 2, pp.Len())
	})

	t.Run("Validate", func(t *testing.T) {
		pp := NewPlantPhotos()
		assert.NoError(t, pp.Validate())
		require.NoError(t, pp.Add(photo1))
		assert.NoError(t, pp.Validate())
	})
}
