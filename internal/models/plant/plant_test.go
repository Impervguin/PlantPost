package plant

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockPlantSpecification для тестирования
type MockPlantSpecification struct{}

func (m *MockPlantSpecification) Validate() error {
	return nil
}

func TestPlant(t *testing.T) {
	// Подготовка тестовых данных
	validID := uuid.New()
	validFileID := uuid.New()
	validTime := time.Now().Add(-time.Hour)
	validSpec := &MockPlantSpecification{}
	validPhoto, _ := CreatePlantPhoto(uuid.New(), validFileID, "test photo")
	validPhotos := NewPlantPhotos()
	_ = validPhotos.Add(validPhoto)

	t.Run("CreatePlant - успешное создание", func(t *testing.T) {
		plant, err := CreatePlant(
			validID,
			"Test Plant",
			"Testus Plantus",
			"Test description",
			validFileID,
			*validPhotos,
			"Test Category",
			validSpec,
			validTime,
		)

		require.NoError(t, err)
		assert.NotNil(t, plant)
		assert.Equal(t, validID, plant.ID())
		assert.Equal(t, validFileID, plant.MainPhotoID())
		assert.Equal(t, *validPhotos, plant.GetPhotos())
		assert.Equal(t, "Test Category", plant.GetCategory())
		assert.Equal(t, validSpec, plant.GetSpecification())
		assert.Equal(t, validTime, plant.CreatedAt())
		assert.Equal(t, validFileID, plant.MainPhotoID())
		assert.Equal(t, "Test Plant", plant.GetName())
		assert.Equal(t, "Testus Plantus", plant.GetLatinName())
		assert.Equal(t, "Test description", plant.GetDescription())

	})

	t.Run("CreatePlant - ошибки валидации", func(t *testing.T) {
		testCases := []struct {
			name        string
			id          uuid.UUID
			nameStr     string
			latinName   string
			description string
			mainPhotoID uuid.UUID
			photos      PlantPhotos
			category    string
			spec        PlantSpecification
			createdAt   time.Time
			expectError bool
		}{
			{
				name:        "empty ID",
				id:          uuid.Nil,
				nameStr:     "Test",
				latinName:   "Testus",
				description: "Desc",
				mainPhotoID: validFileID,
				photos:      *validPhotos,
				category:    "Cat",
				spec:        validSpec,
				createdAt:   validTime,
				expectError: true,
			},
			{
				name:        "empty name",
				id:          validID,
				nameStr:     "",
				latinName:   "Testus",
				description: "Desc",
				mainPhotoID: validFileID,
				photos:      *validPhotos,
				category:    "Cat",
				spec:        validSpec,
				createdAt:   validTime,
				expectError: true,
			},
			{
				name:        "future creation date",
				id:          validID,
				nameStr:     "Test",
				latinName:   "Testus",
				description: "Desc",
				mainPhotoID: validFileID,
				photos:      *validPhotos,
				category:    "Cat",
				spec:        validSpec,
				createdAt:   time.Now().Add(time.Hour),
				expectError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := CreatePlant(
					tc.id,
					tc.nameStr,
					tc.latinName,
					tc.description,
					tc.mainPhotoID,
					tc.photos,
					tc.category,
					tc.spec,
					tc.createdAt,
				)
				if tc.expectError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			})
		}
	})

	t.Run("NewPlant - успешное создание", func(t *testing.T) {
		plant, err := NewPlant(
			"Test Plant",
			"Testus Plantus",
			"Test description",
			validFileID,
			*validPhotos,
			"Test Category",
			validSpec,
		)

		require.NoError(t, err)
		assert.NotNil(t, plant)
		assert.NotEqual(t, uuid.Nil, plant.ID())
	})

	t.Run("UpdateSpec - успешное обновление", func(t *testing.T) {
		plant, _ := CreatePlant(
			validID,
			"Test Plant",
			"Testus Plantus",
			"Test description",
			validFileID,
			*validPhotos,
			"Test Category",
			validSpec,
			validTime,
		)

		newSpec := &MockPlantSpecification{}
		err := plant.UpdateSpec(newSpec)
		require.NoError(t, err)
	})

	t.Run("AddPhoto - успешное добавление", func(t *testing.T) {
		plant, _ := CreatePlant(
			validID,
			"Test Plant",
			"Testus Plantus",
			"Test description",
			validFileID,
			*NewPlantPhotos(),
			"Test Category",
			validSpec,
			validTime,
		)

		newPhoto, _ := CreatePlantPhoto(uuid.New(), uuid.New(), "new photo")
		err := plant.AddPhoto(newPhoto)
		require.NoError(t, err)
	})

	t.Run("AddPhoto - ошибка при дубликате", func(t *testing.T) {
		plant, _ := CreatePlant(
			validID,
			"Test Plant",
			"Testus Plantus",
			"Test description",
			validFileID,
			*validPhotos,
			"Test Category",
			validSpec,
			validTime,
		)

		err := plant.AddPhoto(validPhoto)
		require.Error(t, err)
	})

	t.Run("Getters", func(t *testing.T) {
		plant, _ := CreatePlant(
			validID,
			"Test Plant",
			"Testus Plantus",
			"Test description",
			validFileID,
			*validPhotos,
			"Test Category",
			validSpec,
			validTime,
		)

		assert.Equal(t, validID, plant.ID())
		assert.Equal(t, "Test Plant", plant.GetName())
		assert.Equal(t, validFileID, plant.MainPhotoID())
		assert.Equal(t, *validPhotos, plant.GetPhotos())
	})
}
