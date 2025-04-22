package plantservice_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"PlantSite/internal/models"
	"PlantSite/internal/models/plant"
	authservice "PlantSite/internal/services/auth-service"
	authmock "PlantSite/internal/services/auth-service/auth-mock"
	plantservice "PlantSite/internal/services/plant-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreatePlant(t *testing.T) {
	validSessionID := uuid.New()
	validOwnerID := uuid.New()
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
		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validOwnerID,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		user := new(authmock.MockUser)
		// user.On("ID").Return(validOwnerID)
		user.On("HasAuthorRights").Return(true)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validOwnerID).Return(user, nil)

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		// Setup expectations
		crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
		frepo.On("Upload", mock.Anything, &validMainPhoto).Return(&models.File{ID: validFileID}, nil)
		prepo.On("Create", mock.Anything, mock.AnythingOfType("*plant.Plant")).Return(&plant.Plant{}, nil)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

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
		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		// validSession := &authservice.Session{
		// 	ID:        validSessionID,
		// 	MemberID:  validOwnerID,
		// 	ExpiresAt: time.Now().Add(time.Hour),
		// }
		user := new(authmock.MockUser)
		user.On("ID").Return(validOwnerID)
		user.On("HasAuthorRights").Return(true)
		sessions.On("Get", ctx, validSessionID).Return(nil, assert.AnError)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validOwnerID).Return(user, nil)

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		err := svc.CreatePlant(ctx, validData, validMainPhoto)
		require.Error(t, err)
	})

	t.Run("NotAuthor", func(t *testing.T) {
		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validOwnerID,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		user := new(authmock.MockUser)
		user.On("ID").Return(validOwnerID)
		user.On("HasAuthorRights").Return(false)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validOwnerID).Return(user, nil)

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		err := svc.CreatePlant(ctx, validData, validMainPhoto)
		require.Error(t, err)
		assert.ErrorIs(t, err, plantservice.ErrNotAuthor)
	})

	t.Run("InvalidCategory", func(t *testing.T) {
		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validOwnerID,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		user := new(authmock.MockUser)
		user.On("ID").Return(validOwnerID)
		user.On("HasAuthorRights").Return(true)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validOwnerID).Return(user, nil)

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		crepo.On("GetCategory", mock.Anything, validCategoryName).Return(nil, assert.AnError)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		err := svc.CreatePlant(ctx, validData, validMainPhoto)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("FileUploadError", func(t *testing.T) {
		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validOwnerID,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		user := new(authmock.MockUser)
		user.On("ID").Return(validOwnerID)
		user.On("HasAuthorRights").Return(true)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validOwnerID).Return(user, nil)

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
		frepo.On("Upload", mock.Anything, &validMainPhoto).Return(nil, assert.AnError)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		err := svc.CreatePlant(ctx, validData, validMainPhoto)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("InvalidPlantData", func(t *testing.T) {
		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validOwnerID,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		user := new(authmock.MockUser)
		user.On("ID").Return(validOwnerID)
		user.On("HasAuthorRights").Return(true)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validOwnerID).Return(user, nil)

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		invalidData := plantservice.CreatePlantData{
			Name:        "", // Invalid empty name
			LatinName:   "Rosa",
			Description: "Beautiful flower",
			Category:    validCategoryName,
			Spec:        validSpec,
		}

		crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
		frepo.On("Upload", mock.Anything, &validMainPhoto).Return(&models.File{ID: validFileID}, nil)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		err := svc.CreatePlant(ctx, invalidData, validMainPhoto)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "plant name cannot be empty")
	})

	t.Run("RepositoryCreateError", func(t *testing.T) {
		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validOwnerID,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		user := new(authmock.MockUser)
		user.On("ID").Return(validOwnerID)
		user.On("HasAuthorRights").Return(true)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validOwnerID).Return(user, nil)

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
		frepo.On("Upload", mock.Anything, &validMainPhoto).Return(&models.File{ID: validFileID}, nil)
		prepo.On("Create", mock.Anything, mock.AnythingOfType("*plant.Plant")).Return(nil, assert.AnError)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		err := svc.CreatePlant(ctx, validData, validMainPhoto)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("InvalidSpecification", func(t *testing.T) {
		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		validSession := &authservice.Session{
			ID:        validSessionID,
			MemberID:  validOwnerID,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		user := new(authmock.MockUser)
		user.On("ID").Return(validOwnerID)
		user.On("HasAuthorRights").Return(true)
		sessions.On("Get", ctx, validSessionID).Return(validSession, nil)
		ctx := asvc.Authenticate(ctx, validSessionID)
		arepo.On("Get", ctx, validOwnerID).Return(user, nil)

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		invalidSpec := new(MockPlantSpecification)
		invalidSpec.On("Validate").Return(assert.AnError)

		invalidData := plantservice.CreatePlantData{
			Name:        "Rose",
			LatinName:   "Rosa",
			Description: "Beautiful flower",
			Category:    validCategoryName,
			Spec:        invalidSpec,
		}

		crepo.On("GetCategory", mock.Anything, validCategoryName).Return(&plant.PlantCategory{}, nil)
		frepo.On("Upload", mock.Anything, &validMainPhoto).Return(&models.File{ID: validFileID}, nil)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		err := svc.CreatePlant(ctx, invalidData, validMainPhoto)
		require.Error(t, err)

		invalidSpec.AssertExpectations(t)
	})
}
