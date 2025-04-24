package plantservice_test

import (
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

// MockPlantRepository implements plant.PlantRepository interface
type MockPlantRepository struct {
	mock.Mock
}

func (m *MockPlantRepository) Create(ctx context.Context, p *plant.Plant) (*plant.Plant, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*plant.Plant), args.Error(1)
}

func (m *MockPlantRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(*plant.Plant) (*plant.Plant, error)) (*plant.Plant, error) {
	args := m.Called(ctx, id, updateFn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	pl, err := updateFn(args.Get(0).(*plant.Plant))
	return pl, err
}

func (m *MockPlantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPlantRepository) Get(ctx context.Context, id uuid.UUID) (*plant.Plant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*plant.Plant), args.Error(1)
}

// MockPlantCategoryRepository implements plant.PlantCategoryRepository interface
type MockPlantCategoryRepository struct {
	mock.Mock
}

func (m *MockPlantCategoryRepository) GetCategory(ctx context.Context, name string) (*plant.PlantCategory, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*plant.PlantCategory), args.Error(1)
}

func (m *MockPlantCategoryRepository) GetCategories(ctx context.Context) ([]plant.PlantCategory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]plant.PlantCategory), args.Error(1)
}

// MockFileRepository implements models.FileRepository interface
type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) Upload(ctx context.Context, fdata *models.FileData) (*models.File, error) {
	args := m.Called(ctx, fdata)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.File), args.Error(1)
}

func (m *MockFileRepository) Get(ctx context.Context, id uuid.UUID) (*models.File, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.File), args.Error(1)
}

func (m *MockFileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFileRepository) Download(ctx context.Context, fileID uuid.UUID) (*models.FileData, error) {
	args := m.Called(ctx, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FileData), args.Error(1)
}

func (m *MockFileRepository) Update(ctx context.Context, fileID uuid.UUID, data *models.FileData) (*models.File, error) {
	args := m.Called(ctx, fileID, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.File), args.Error(1)
}

// MockPlantSpecification implements plant.PlantSpecification interface
type MockPlantSpecification struct {
	mock.Mock
}

func (m *MockPlantSpecification) Validate() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPlantSpecification) Category() string {
	return m.Called().String(0)
}

func TestPlantService(t *testing.T) {
	ctx := context.Background()
	validSessionID := uuid.New()
	validOwnerID := uuid.New()
	validPlantID := uuid.New()
	validFileID := uuid.New()
	validCategoryName := "mock"

	// Create a valid plant for testing
	validSpec := new(MockPlantSpecification)
	validSpec.On("Validate").Return(nil)
	validPlant, err := plant.NewPlant(
		"Rose",
		"Rosa",
		"Beautiful flower",
		validFileID,
		plant.PlantPhotos{},
		validCategoryName,
		validSpec,
	)
	require.NoError(t, err)

	t.Run("UpdatePlantSpec", func(t *testing.T) {
		newSpec := new(MockPlantSpecification)
		newSpec.On("Validate").Return(nil)
		newSpec.On("Category").Return(validCategoryName)

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

			user.On("HasAuthorRights").Return(true)
			prepo.On("Update", mock.Anything, validPlantID, mock.Anything).Return(validPlant, nil)

			svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

			err := svc.UpdatePlantSpec(ctx, validPlantID, newSpec)
			require.NoError(t, err)

			prepo.AssertExpectations(t)
			user.AssertExpectations(t)
			newSpec.AssertExpectations(t)
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

			invalidSpec := new(MockPlantSpecification)
			invalidSpec.On("Validate").Return(assert.AnError)

			prepo := new(MockPlantRepository)
			crepo := new(MockPlantCategoryRepository)
			frepo := new(MockFileRepository)

			prepo.On("Update", mock.Anything, validPlantID, mock.Anything).Return(validPlant, nil)

			svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

			err := svc.UpdatePlantSpec(ctx, validPlantID, invalidSpec)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)

			invalidSpec.AssertExpectations(t)
		})
	})

	t.Run("DeletePlant", func(t *testing.T) {
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

			prepo.On("Delete", mock.Anything, validPlantID).Return(nil)

			svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

			err := svc.DeletePlant(ctx, validPlantID)
			require.NoError(t, err)

			prepo.AssertExpectations(t)
			user.AssertExpectations(t)
		})
	})

	t.Run("UploadPlantPhoto", func(t *testing.T) {
		fdata := models.FileData{}
		description := "test photo"

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

			frepo.On("Upload", mock.Anything, &fdata).Return(&models.File{ID: validFileID}, nil)
			prepo.On("Update", mock.Anything, validPlantID, mock.Anything).Return(validPlant, nil)

			svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

			err := svc.UploadPlantPhoto(ctx, validPlantID, fdata, description)
			require.NoError(t, err)

			frepo.AssertExpectations(t)
			prepo.AssertExpectations(t)
			user.AssertExpectations(t)
		})

		t.Run("PhotoAddError", func(t *testing.T) {
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

			frepo.On("Upload", mock.Anything, &fdata).Return(&models.File{ID: validFileID}, nil)
			prepo.On("Update", mock.Anything, validPlantID, mock.Anything).Return(nil, assert.AnError)

			svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

			err := svc.UploadPlantPhoto(ctx, validPlantID, fdata, description)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})
}
