package plantservice_test

import (
	"context"
	"testing"
	"time"

	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/plant"
	authservice "PlantSite/internal/services/auth-service"
	authmock "PlantSite/internal/services/auth-service/auth-mock"
	plantservice "PlantSite/internal/services/plant-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetPlant(t *testing.T) {
	validSessionID := uuid.New()
	validOwnerID := uuid.New()

	ctx := context.Background()
	validPlantID := uuid.New()
	validFileID := uuid.New()
	validCategoryName := "flowers"

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

	expectedResult := &plantservice.GetPlant{
		ID:          validPlant.ID(),
		Name:        validPlant.GetName(),
		LatinName:   validPlant.GetLatinName(),
		Description: validPlant.GetDescription(),
		MainPhoto:   *mainPhotoFile,
		Photos: []plantservice.GetPlantPhoto{
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
		prepo.On("Get", mock.Anything, validPlantID).Return(validPlant, nil)
		frepo.On("Get", mock.Anything, mainPhotoFile.ID).Return(mainPhotoFile, nil)
		frepo.On("Get", mock.Anything, photoFile.ID).Return(photoFile, nil)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		result, err := svc.GetPlant(ctx, validPlantID)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)

		// Verify all expectations were met
		user.AssertExpectations(t)
		prepo.AssertExpectations(t)
		frepo.AssertExpectations(t)
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

		_, err := svc.GetPlant(ctx, validPlantID)
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

		_, err := svc.GetPlant(ctx, validPlantID)
		require.Error(t, err)
		assert.ErrorIs(t, err, auth.ErrNoAuthorRights)
	})

	t.Run("PlantNotFound", func(t *testing.T) {
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

		prepo.On("Get", mock.Anything, validPlantID).Return(nil, assert.AnError)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		_, err := svc.GetPlant(ctx, validPlantID)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("MainPhotoNotFound", func(t *testing.T) {
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

		prepo.On("Get", mock.Anything, validPlantID).Return(validPlant, nil)
		frepo.On("Get", mock.Anything, mainPhotoFile.ID).Return(nil, assert.AnError)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		_, err := svc.GetPlant(ctx, validPlantID)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("AdditionalPhotoNotFound", func(t *testing.T) {
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

		prepo.On("Get", mock.Anything, validPlantID).Return(validPlant, nil)
		frepo.On("Get", mock.Anything, mainPhotoFile.ID).Return(mainPhotoFile, nil)
		frepo.On("Get", mock.Anything, photoFile.ID).Return(nil, assert.AnError)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		_, err := svc.GetPlant(ctx, validPlantID)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("EmptyPhotos", func(t *testing.T) {
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

		// Plant with no additional photos
		plantNoPhotos, err := plant.NewPlant(
			"Rose",
			"Rosa",
			"Beautiful flower",
			mainPhotoFile.ID,
			*plant.NewPlantPhotos(),
			validCategoryName,
			validSpec,
		)
		require.NoError(t, err)

		prepo.On("Get", mock.Anything, validPlantID).Return(plantNoPhotos, nil)
		frepo.On("Get", mock.Anything, mainPhotoFile.ID).Return(mainPhotoFile, nil)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		result, err := svc.GetPlant(ctx, validPlantID)
		require.NoError(t, err)
		assert.Empty(t, result.Photos)
		assert.Equal(t, *mainPhotoFile, result.MainPhoto)
	})
}
