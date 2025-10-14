//go:build unit

package plantservice_test

import (
	"context"
	"testing"

	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/plant"
	authservice "PlantSite/internal/services/auth-service"
	authmock "PlantSite/internal/services/auth-service/auth-mock"
	plantservice "PlantSite/internal/services/plant-service"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type PlantServiceGetTestSuite struct {
	suite.Suite
	plantMother *PlantMother
}

func (s *PlantServiceGetTestSuite) BeforeEach(t provider.T) {
	t.Epic("Plant Service")
	t.Feature("Plant Retrieval")
	s.plantMother = &PlantMother{}
}

func (s *PlantServiceGetTestSuite) TestGetPlant(t provider.T) {
	t.Tags("retrieval", "positive")
	t.Description("Test plant retrieval functionality")
	t.Parallel()

	validOwnerID := uuid.New()
	validPlantID := uuid.New()

	t.Run("Successful plant retrieval", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		// Create test plant using Object Mother
		validPlant, _, err := s.plantMother.CreateValidPlant()
		require.NoError(t, err)

		// Create test files
		mainPhotoFile := &models.File{ID: uuid.New(), Name: "main.jpg"}
		photoFile := &models.File{ID: uuid.New(), Name: "photo.jpg"}

		// Add photos to plant
		plantPhoto, err := plant.NewPlantPhoto(photoFile.ID, "additional photo")
		require.NoError(t, err)

		// Update plant with photos
		updatedPlant := *validPlant
		updatedPlant.UpdateMainPhotoID(mainPhotoFile.ID)
		updatedPlant.AddPhoto(plantPhoto)

		expectedResult := &plantservice.GetPlant{
			ID:          updatedPlant.ID(),
			Name:        updatedPlant.GetName(),
			LatinName:   updatedPlant.GetLatinName(),
			Description: updatedPlant.GetDescription(),
			MainPhoto:   *mainPhotoFile,
			Photos: []plantservice.GetPlantPhoto{
				{
					ID:          plantPhoto.ID(),
					File:        *photoFile,
					Description: plantPhoto.Description(),
				},
			},
			Category:      updatedPlant.GetCategory(),
			Specification: updatedPlant.GetSpecification(),
			CreatedAt:     updatedPlant.CreatedAt(),
		}

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup repositories", func(pctx provider.StepCtx) {
			prepo.On("Get", mock.Anything, validPlantID).Return(&updatedPlant, nil)
			frepo.On("Get", mock.Anything, mainPhotoFile.ID).Return(mainPhotoFile, nil)
			frepo.On("Get", mock.Anything, photoFile.ID).Return(photoFile, nil)
		})

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		var result *plantservice.GetPlant
		t.WithNewStep("Retrieve plant", func(pctx provider.StepCtx) {
			result, err = svc.GetPlant(ctx, validPlantID)
		})

		t.WithNewStep("Verify retrieval success", func(pctx provider.StepCtx) {
			require.NoError(t, err)
			assert.Equal(t, expectedResult, result)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			prepo.AssertExpectations(t)
			frepo.AssertExpectations(t)

		})
	})

	t.Run("Not authorized for plant retrieval", func(t provider.T) {
		t.Parallel()

		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)

		t.WithNewStep("Setup authentication failure", func(pctx provider.StepCtx) {
			sessions.On("Get", mock.Anything, mock.Anything).Return(nil, assert.AnError)
		})

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Attempt unauthorized retrieval", func(pctx provider.StepCtx) {
			_, err := svc.GetPlant(context.Background(), validPlantID)
			require.Error(t, err)
		})
	})

	t.Run("No author rights for plant retrieval", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, false)

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Attempt retrieval without author rights", func(pctx provider.StepCtx) {
			_, err := svc.GetPlant(ctx, validPlantID)
			require.Error(t, err)
			assert.ErrorIs(t, err, auth.ErrNoAuthorRights)
		})
	})

	t.Run("Plant not found during retrieval", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup plant not found", func(pctx provider.StepCtx) {
			prepo.On("Get", mock.Anything, validPlantID).Return(nil, assert.AnError)
		})

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Attempt retrieval of non-existent plant", func(pctx provider.StepCtx) {
			_, err := svc.GetPlant(ctx, validPlantID)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Main photo not found during retrieval", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		validPlant, _, err := s.plantMother.CreateValidPlant()
		require.NoError(t, err)

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup main photo not found", func(pctx provider.StepCtx) {
			prepo.On("Get", mock.Anything, validPlantID).Return(validPlant, nil)
			frepo.On("Get", mock.Anything, validPlant.MainPhotoID()).Return(nil, assert.AnError)
		})

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Attempt retrieval with missing main photo", func(pctx provider.StepCtx) {
			_, err := svc.GetPlant(ctx, validPlantID)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {

		})
	})

	t.Run("Plant with empty photos", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		validPlant, _, err := s.plantMother.CreateValidPlant()
		require.NoError(t, err)

		mainPhotoFile := &models.File{ID: validPlant.MainPhotoID(), Name: "main.jpg"}

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup plant with empty photos", func(pctx provider.StepCtx) {
			prepo.On("Get", mock.Anything, validPlantID).Return(validPlant, nil)
			frepo.On("Get", mock.Anything, mainPhotoFile.ID).Return(mainPhotoFile, nil)
		})

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		var result *plantservice.GetPlant
		t.WithNewStep("Retrieve plant with empty photos", func(pctx provider.StepCtx) {
			result, err = svc.GetPlant(ctx, validPlantID)
		})

		t.WithNewStep("Verify empty photos result", func(pctx provider.StepCtx) {
			require.NoError(t, err)
			assert.Empty(t, result.Photos)
			assert.Equal(t, *mainPhotoFile, result.MainPhoto)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {

		})
	})
}

func TestPlantServiceGet(t *testing.T) {
	suite.RunSuite(t, new(PlantServiceGetTestSuite))
}
