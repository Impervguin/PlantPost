package plantservice_test

import (
	"bytes"
	"context"
	"testing"

	"PlantSite/internal/models"
	"PlantSite/internal/models/plant"
	plantservice "PlantSite/internal/services/plant-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreatePlant(t *testing.T) {
	ctx := context.Background()
	validFileID := uuid.New()
	validCategoryName := "flowers"
	validMainPhoto := models.FileData{Name: "plant.jpg", Reader: bytes.NewReader([]byte("image data"))}

	validSpec := new(MockPlantSpecification)
	validSpec.On("Validate").Return(nil)

	validData := plantservice.CreatePlantData{
		Name:        "Rose",
		LatinName:   "Rosa",
		Description: "Beautiful flower",
		Category:    validCategoryName,
		Spec:        validSpec,
	}

	t.Run("Success", func(t *testing.T) {
		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		// Setup expectations
		user.On("HasAuthorRights").Return(true)
		crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
		frepo.On("Upload", mock.Anything, &validMainPhoto).Return(&models.File{ID: validFileID}, nil)
		prepo.On("Create", mock.Anything, mock.AnythingOfType("*plant.Plant")).Return(&plant.Plant{}, nil)

		svc := plantservice.NewPlantService(prepo, crepo, frepo)
		ctx := PutUserInContext(ctx, user)

		err := svc.CreatePlant(ctx, validData, validMainPhoto)
		require.NoError(t, err)

		// Verify all expectations were met
		user.AssertExpectations(t)
		crepo.AssertExpectations(t)
		frepo.AssertExpectations(t)
		prepo.AssertExpectations(t)
		validSpec.AssertExpectations(t)
	})

	t.Run("NotAuthorized", func(t *testing.T) {
		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		svc := plantservice.NewPlantService(prepo, crepo, frepo)

		err := svc.CreatePlant(ctx, validData, validMainPhoto)
		require.Error(t, err)
		assert.ErrorIs(t, err, plantservice.ErrNotAuthorized)
	})

	t.Run("NotAuthor", func(t *testing.T) {
		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("HasAuthorRights").Return(false)

		svc := plantservice.NewPlantService(prepo, crepo, frepo)
		ctx := PutUserInContext(ctx, user)

		err := svc.CreatePlant(ctx, validData, validMainPhoto)
		require.Error(t, err)
		assert.ErrorIs(t, err, plantservice.ErrNotAuthor)
	})

	t.Run("InvalidCategory", func(t *testing.T) {
		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("HasAuthorRights").Return(true)
		crepo.On("GetCategory", mock.Anything, validCategoryName).Return(nil, assert.AnError)

		svc := plantservice.NewPlantService(prepo, crepo, frepo)
		ctx := PutUserInContext(ctx, user)

		err := svc.CreatePlant(ctx, validData, validMainPhoto)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("FileUploadError", func(t *testing.T) {
		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("HasAuthorRights").Return(true)
		crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
		frepo.On("Upload", mock.Anything, &validMainPhoto).Return(nil, assert.AnError)

		svc := plantservice.NewPlantService(prepo, crepo, frepo)
		ctx := PutUserInContext(ctx, user)

		err := svc.CreatePlant(ctx, validData, validMainPhoto)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("InvalidPlantData", func(t *testing.T) {
		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		invalidData := plantservice.CreatePlantData{
			Name:        "", // Invalid empty name
			LatinName:   "Rosa",
			Description: "Beautiful flower",
			Category:    validCategoryName,
			Spec:        validSpec,
		}

		user.On("HasAuthorRights").Return(true)
		crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
		frepo.On("Upload", mock.Anything, &validMainPhoto).Return(&models.File{ID: validFileID}, nil)

		svc := plantservice.NewPlantService(prepo, crepo, frepo)
		ctx := PutUserInContext(ctx, user)

		err := svc.CreatePlant(ctx, invalidData, validMainPhoto)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "plant name cannot be empty")
	})

	t.Run("RepositoryCreateError", func(t *testing.T) {
		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		user.On("HasAuthorRights").Return(true)
		crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
		frepo.On("Upload", mock.Anything, &validMainPhoto).Return(&models.File{ID: validFileID}, nil)
		prepo.On("Create", mock.Anything, mock.AnythingOfType("*plant.Plant")).Return(nil, assert.AnError)

		svc := plantservice.NewPlantService(prepo, crepo, frepo)
		ctx := PutUserInContext(ctx, user)

		err := svc.CreatePlant(ctx, validData, validMainPhoto)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("InvalidSpecification", func(t *testing.T) {
		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)
		user := new(MockUser)

		invalidSpec := new(MockPlantSpecification)
		invalidSpec.On("Validate").Return(assert.AnError)

		invalidData := plantservice.CreatePlantData{
			Name:        "Rose",
			LatinName:   "Rosa",
			Description: "Beautiful flower",
			Category:    validCategoryName,
			Spec:        invalidSpec,
		}

		user.On("HasAuthorRights").Return(true)
		crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
		frepo.On("Upload", mock.Anything, &validMainPhoto).Return(&models.File{ID: validFileID}, nil)

		svc := plantservice.NewPlantService(prepo, crepo, frepo)
		ctx := PutUserInContext(ctx, user)

		err := svc.CreatePlant(ctx, invalidData, validMainPhoto)
		require.Error(t, err)

		invalidSpec.AssertExpectations(t)
	})
}
