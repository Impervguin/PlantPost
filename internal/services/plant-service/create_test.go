//go:build unit

package plantservice_test

import (
	"bytes"
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

type PlantServiceCreateTestSuite struct {
	suite.Suite
	plantMother *PlantMother
}

func (s *PlantServiceCreateTestSuite) BeforeEach(t provider.T) {
	t.Epic("Plant Service")
	t.Feature("Plant Creation")
	s.plantMother = &PlantMother{}
}

func (s *PlantServiceCreateTestSuite) TestCreatePlant(t provider.T) {
	t.Tags("creation", "positive")
	t.Description("Test plant creation functionality")
	t.Parallel()

	validOwnerID := uuid.New()
	validFileID := uuid.New()
	validCategoryName := "flowers"
	validMainPhoto := models.FileData{Name: "plant.jpg", Reader: bytes.NewReader([]byte("image data"))}

	t.Run("Successful plant creation", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		validSpec := new(MockPlantSpecification)
		t.WithNewStep("Setup valid specification", func(pctx provider.StepCtx) {
			validSpec.On("Validate").Return(nil)
		})

		validData := plantservice.CreatePlantData{
			Name:        "Rose",
			LatinName:   "Rosa",
			Description: "Beautiful flower",
			Category:    validCategoryName,
			Spec:        validSpec,
		}

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup repositories", func(pctx provider.StepCtx) {
			crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
			frepo.On("Upload", mock.Anything, &validMainPhoto).Return(&models.File{ID: validFileID}, nil)
			prepo.On("Create", mock.Anything, mock.AnythingOfType("*plant.Plant")).Return(&plant.Plant{}, nil)
		})

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Create plant", func(pctx provider.StepCtx) {
			err := svc.CreatePlant(ctx, validData, validMainPhoto)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			crepo.AssertExpectations(t)
			frepo.AssertExpectations(t)
			prepo.AssertExpectations(t)

		})
	})

	t.Run("Not authorized for plant creation", func(t provider.T) {
		t.Parallel()

		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)

		t.WithNewStep("Setup authentication failure", func(pctx provider.StepCtx) {
			sessions.On("Get", mock.Anything, mock.Anything).Return(nil, assert.AnError)
		})

		validSpec := new(MockPlantSpecification)
		validData := plantservice.CreatePlantData{
			Name:        "Rose",
			LatinName:   "Rosa",
			Description: "Beautiful flower",
			Category:    validCategoryName,
			Spec:        validSpec,
		}

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Attempt unauthorized creation", func(pctx provider.StepCtx) {
			err := svc.CreatePlant(context.Background(), validData, validMainPhoto)
			require.Error(t, err)
		})
	})

	t.Run("No author rights for plant creation", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, false)

		validSpec := new(MockPlantSpecification)
		validData := plantservice.CreatePlantData{
			Name:        "Rose",
			LatinName:   "Rosa",
			Description: "Beautiful flower",
			Category:    validCategoryName,
			Spec:        validSpec,
		}

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Attempt creation without author rights", func(pctx provider.StepCtx) {
			err := svc.CreatePlant(ctx, validData, validMainPhoto)
			require.Error(t, err)
			assert.ErrorIs(t, err, auth.ErrNoAuthorRights)
		})
	})

	t.Run("Invalid category during creation", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		validSpec := new(MockPlantSpecification)
		validData := plantservice.CreatePlantData{
			Name:        "Rose",
			LatinName:   "Rosa",
			Description: "Beautiful flower",
			Category:    validCategoryName,
			Spec:        validSpec,
		}

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup category error", func(pctx provider.StepCtx) {
			crepo.On("GetCategory", mock.Anything, validCategoryName).Return(nil, assert.AnError)
		})

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Attempt creation with invalid category", func(pctx provider.StepCtx) {
			err := svc.CreatePlant(ctx, validData, validMainPhoto)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("File upload error during creation", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		validSpec := new(MockPlantSpecification)
		validData := plantservice.CreatePlantData{
			Name:        "Rose",
			LatinName:   "Rosa",
			Description: "Beautiful flower",
			Category:    validCategoryName,
			Spec:        validSpec,
		}

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup file upload error", func(pctx provider.StepCtx) {
			crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
			frepo.On("Upload", mock.Anything, &validMainPhoto).Return(nil, assert.AnError)
		})

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Attempt creation with upload error", func(pctx provider.StepCtx) {
			err := svc.CreatePlant(ctx, validData, validMainPhoto)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Invalid plant data during creation", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		validSpec := new(MockPlantSpecification)
		invalidData := plantservice.CreatePlantData{
			Name:        "", // Invalid empty name
			LatinName:   "Rosa",
			Description: "Beautiful flower",
			Category:    validCategoryName,
			Spec:        validSpec,
		}

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup valid category and file upload", func(pctx provider.StepCtx) {
			crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
			frepo.On("Upload", mock.Anything, &validMainPhoto).Return(&models.File{ID: validFileID}, nil)
		})

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Attempt creation with invalid data", func(pctx provider.StepCtx) {
			err := svc.CreatePlant(ctx, invalidData, validMainPhoto)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "plant name cannot be empty")
		})
	})

	t.Run("Repository create error during creation", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		validSpec := new(MockPlantSpecification)
		t.WithNewStep("Setup valid specification", func(pctx provider.StepCtx) {
			validSpec.On("Validate").Return(nil)
		})

		validData := plantservice.CreatePlantData{
			Name:        "Rose",
			LatinName:   "Rosa",
			Description: "Beautiful flower",
			Category:    validCategoryName,
			Spec:        validSpec,
		}

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup repository create error", func(pctx provider.StepCtx) {
			crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
			frepo.On("Upload", mock.Anything, &validMainPhoto).Return(&models.File{ID: validFileID}, nil)
			prepo.On("Create", mock.Anything, mock.AnythingOfType("*plant.Plant")).Return(nil, assert.AnError)
		})

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Attempt creation with repository error", func(pctx provider.StepCtx) {
			err := svc.CreatePlant(ctx, validData, validMainPhoto)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {

		})
	})

	t.Run("Invalid specification during creation", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		invalidSpec := new(MockPlantSpecification)
		t.WithNewStep("Setup invalid specification", func(pctx provider.StepCtx) {
			invalidSpec.On("Validate").Return(assert.AnError)
		})

		invalidData := plantservice.CreatePlantData{
			Name:        "Rose",
			LatinName:   "Rosa",
			Description: "Beautiful flower",
			Category:    validCategoryName,
			Spec:        invalidSpec,
		}

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup valid category and file upload", func(pctx provider.StepCtx) {
			crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
			frepo.On("Upload", mock.Anything, &validMainPhoto).Return(&models.File{ID: validFileID}, nil)
		})

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Attempt creation with invalid specification", func(pctx provider.StepCtx) {
			err := svc.CreatePlant(ctx, invalidData, validMainPhoto)
			require.Error(t, err)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			invalidSpec.AssertExpectations(t)
		})
	})
}

func TestPlantServiceCreate(t *testing.T) {
	suite.RunSuite(t, new(PlantServiceCreateTestSuite))
}
