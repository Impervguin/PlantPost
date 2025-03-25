package searchservice_test

import (
	"context"
	"testing"

	"PlantSite/internal/models"
	"PlantSite/internal/models/plant"
	searchservice "PlantSite/internal/search-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockPlantSpecification implements plant.PlantSpecification interface
type MockPlantSpecification struct {
	mock.Mock
}

func (m *MockPlantSpecification) Validate() error {
	args := m.Called()
	return args.Error(0)
}

func TestGetPlantByID(t *testing.T) {
	ctx := context.Background()
	validPlantID := uuid.New()
	validFileID := uuid.New()
	validCategoryName := "flowers"
	// validPhotoID := uuid.New()
	// createdAt := time.Now()

	// Create test data
	validSpec := new(MockPlantSpecification)
	validSpec.On("Validate").Return(nil)

	mainPhotoFile := &models.File{ID: validFileID, Name: "main.jpg"}
	photoFile := &models.File{ID: uuid.New(), Name: "photo.jpg"}

	plantPhotos := plant.NewPlantPhotos()
	plantPhoto, err := plant.NewPlantPhoto(photoFile.ID, "additional photo")
	require.NoError(t, err)
	require.NoError(t, plantPhotos.Add(plantPhoto))

	validPlant, err := plant.NewPlant(
		"Rose",
		"Rosa",
		"Beautiful flower",
		mainPhotoFile.ID,
		*plantPhotos,
		validCategoryName,
		validSpec,
	)
	require.NoError(t, err)

	expectedResult := searchservice.GetPlant{
		ID:          validPlant.ID(),
		Name:        validPlant.GetName(),
		LatinName:   validPlant.GetLatinName(),
		Description: validPlant.GetDescription(),
		MainPhoto:   *mainPhotoFile,
		Photos: []searchservice.GetPlantPhoto{
			{
				ID:          plantPhoto.ID(),
				File:        *photoFile,
				Description: plantPhoto.Description(),
			},
		},
		Category:      validPlant.GetCategory(),
		Specification: validPlant.GetSpecification(),
		CreatedAt:     validPlant.CreatedAt(),
	}

	t.Run("Success", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		// Setup expectations
		srepo.On("GetPlantByID", ctx, validPlantID).Return(validPlant, nil)
		pfrepo.On("Get", ctx, mainPhotoFile.ID).Return(mainPhotoFile, nil)
		pfrepo.On("Get", ctx, photoFile.ID).Return(photoFile, nil)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		result, err := svc.GetPlantByID(ctx, validPlantID)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, *result)

		// Verify all expectations were met
		srepo.AssertExpectations(t)
		pfrepo.AssertExpectations(t)
	})

	t.Run("PlantNotFound", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		srepo.On("GetPlantByID", ctx, validPlantID).Return(nil, assert.AnError)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		_, err := svc.GetPlantByID(ctx, validPlantID)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("MainPhotoNotFound", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		srepo.On("GetPlantByID", ctx, validPlantID).Return(validPlant, nil)
		pfrepo.On("Get", ctx, mainPhotoFile.ID).Return(nil, assert.AnError)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		_, err := svc.GetPlantByID(ctx, validPlantID)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("AdditionalPhotoNotFound", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		srepo.On("GetPlantByID", ctx, validPlantID).Return(validPlant, nil)
		pfrepo.On("Get", ctx, mainPhotoFile.ID).Return(mainPhotoFile, nil)
		pfrepo.On("Get", ctx, photoFile.ID).Return(nil, assert.AnError)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		_, err := svc.GetPlantByID(ctx, validPlantID)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("NoPhotos", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		noPhotos := plant.NewPlantPhotos()
		plantNoPhotos, err := plant.NewPlant(
			"Rose",
			"Rosa",
			"Beautiful flower",
			mainPhotoFile.ID,
			*noPhotos,
			validCategoryName,
			validSpec,
		)

		srepo.On("GetPlantByID", ctx, validPlantID).Return(plantNoPhotos, nil)
		pfrepo.On("Get", ctx, mainPhotoFile.ID).Return(mainPhotoFile, nil)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		result, err := svc.GetPlantByID(ctx, validPlantID)
		require.NoError(t, err)
		assert.Empty(t, result.Photos)
	})

	t.Run("NilPlantID", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		_, err := svc.GetPlantByID(ctx, uuid.Nil)
		require.Error(t, err)
	})
}
