//go:build unit

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
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// PlantMother реализует паттерн Object Mother для создания тестовых растений
type PlantMother struct{}

func (pm *PlantMother) CreateValidPlant() (*plant.Plant, *MockPlantSpecification, error) {
	spec := new(MockPlantSpecification)
	spec.On("Validate").Return(nil)
	spec.On("Category").Return("mock")

	plant, err := plant.NewPlant(
		"Rose",
		"Rosa",
		"Beautiful flower",
		uuid.New(),
		*plant.NewPlantPhotos(),
		"mock",
		spec,
	)
	return plant, spec, err
}

func (pm *PlantMother) CreatePlantWithCategory(category string) (*plant.Plant, *MockPlantSpecification, error) {
	spec := new(MockPlantSpecification)
	spec.On("Validate").Return(nil)
	spec.On("Category").Return(category)

	plant, err := plant.NewPlant(
		"Test Plant",
		"Testus Plantus",
		"Test description",
		uuid.New(),
		*plant.NewPlantPhotos(),
		category,
		spec,
	)
	return plant, spec, err
}

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

type PlantServiceTestSuite struct {
	suite.Suite
	plantMother *PlantMother
}

func (s *PlantServiceTestSuite) BeforeEach(t provider.T) {
	t.Epic("Plant Service")
	t.Feature("Plant Management")
	s.plantMother = &PlantMother{}
}

func setupAuthService(t provider.T, userID uuid.UUID, hasAuthorRights bool) (*authservice.AuthService, context.Context) {
	arepo := new(authmock.MockAuthRepository)
	sessions := new(authmock.MockSessionStorage)
	hasher := new(authmock.MockPasswdHasher)
	asvc := authservice.NewAuthService(sessions, arepo, hasher)

	sessionID := uuid.New()
	validSession := &authservice.Session{
		ID:        sessionID,
		MemberID:  userID,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	user := new(authmock.MockUser)
	user.On("ID").Return(userID)
	user.On("HasAuthorRights").Return(hasAuthorRights)

	sessions.On("Get", mock.Anything, sessionID).Return(validSession, nil)
	authCtx := asvc.Authenticate(context.Background(), sessionID)
	arepo.On("Get", authCtx, userID).Return(user, nil)

	return asvc, authCtx
}

func (s *PlantServiceTestSuite) TestUpdatePlantSpec(t provider.T) {
	t.Tags("update", "specification")
	t.Description("Test plant specification update functionality")
	t.Parallel()

	validPlantID := uuid.New()
	validOwnerID := uuid.New()
	validCategoryName := "mock"

	t.Run("Successful plant specification update", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		validPlant, _, err := s.plantMother.CreateValidPlant()
		require.NoError(t, err)

		newSpec := new(MockPlantSpecification)
		t.WithNewStep("Setup new specification", func(pctx provider.StepCtx) {
			newSpec.On("Validate").Return(nil)
			newSpec.On("Category").Return(validCategoryName)
		})

		prepo := new(MockPlantRepository)
		t.WithNewStep("Setup plant repository", func(pctx provider.StepCtx) {
			prepo.On("Update", mock.Anything, validPlantID, mock.Anything).Return(validPlant, nil).Run(func(args mock.Arguments) {
				fn, ok := args.Get(2).(func(*plant.Plant) (*plant.Plant, error))
				require.True(t, ok)
				_, err := fn(validPlant)
				require.NoError(t, err)
			})
		})

		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Update plant specification", func(pctx provider.StepCtx) {
			err := svc.UpdatePlantSpec(ctx, validPlantID, newSpec)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			prepo.AssertExpectations(t)
			newSpec.AssertExpectations(t)

		})
	})

	t.Run("Invalid specification during update", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		validPlant, _, err := s.plantMother.CreateValidPlant()
		require.NoError(t, err)

		invalidSpec := new(MockPlantSpecification)
		t.WithNewStep("Setup invalid specification", func(pctx provider.StepCtx) {
			invalidSpec.On("Validate").Return(assert.AnError)
		})

		prepo := new(MockPlantRepository)
		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		t.WithNewStep("Setup plant repository", func(pctx provider.StepCtx) {
			prepo.On("Get", mock.Anything, validPlantID).Return(validPlant, nil)
			prepo.On("Update", mock.Anything, validPlantID, mock.Anything).Return(validPlant, assert.AnError).Run(func(args mock.Arguments) {
				fn, ok := args.Get(2).(func(*plant.Plant) (*plant.Plant, error))
				require.True(t, ok)
				_, err := fn(validPlant)
				require.Error(t, err)
			})
		})

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Attempt update with invalid specification", func(pctx provider.StepCtx) {
			err := svc.UpdatePlantSpec(ctx, validPlantID, invalidSpec)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			invalidSpec.AssertExpectations(t)

		})
	})
}

func (s *PlantServiceTestSuite) TestDeletePlant(t provider.T) {
	t.Tags("delete", "positive")
	t.Description("Test plant deletion functionality")
	t.Parallel()

	validPlantID := uuid.New()
	validOwnerID := uuid.New()

	t.Run("Successful plant deletion", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		prepo := new(MockPlantRepository)
		t.WithNewStep("Setup plant deletion", func(pctx provider.StepCtx) {
			prepo.On("Delete", mock.Anything, validPlantID).Return(nil)
		})

		crepo := new(MockPlantCategoryRepository)
		frepo := new(MockFileRepository)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Delete plant", func(pctx provider.StepCtx) {
			err := svc.DeletePlant(ctx, validPlantID)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			prepo.AssertExpectations(t)
		})
	})
}

func (s *PlantServiceTestSuite) TestUploadPlantPhoto(t provider.T) {
	t.Tags("upload", "photos")
	t.Description("Test plant photo upload functionality")
	t.Parallel()

	validPlantID := uuid.New()
	validOwnerID := uuid.New()
	validFileID := uuid.New()

	t.Run("Successful plant photo upload", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		validPlant, _, err := s.plantMother.CreateValidPlant()
		require.NoError(t, err)

		fdata := models.FileData{}
		description := "test photo"
		newFile := models.File{ID: uuid.New(), Name: "new_file.jpg", URL: "http://new_file.jpg", CreatedAt: time.Now()}

		frepo := new(MockFileRepository)
		t.WithNewStep("Setup file upload", func(pctx provider.StepCtx) {
			frepo.On("Upload", mock.Anything, &fdata).Return(&newFile, nil)
		})

		prepo := new(MockPlantRepository)
		t.WithNewStep("Setup plant update", func(pctx provider.StepCtx) {
			prepo.On("Update", mock.Anything, validPlantID, mock.Anything).Return(validPlant, nil).Run(func(args mock.Arguments) {
				_, ok := args.Get(2).(func(*plant.Plant) (*plant.Plant, error))
				require.True(t, ok)
			})
		})

		crepo := new(MockPlantCategoryRepository)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Upload plant photo", func(pctx provider.StepCtx) {
			err = svc.UploadPlantPhoto(ctx, validPlantID, fdata, description)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			require.NoError(t, err)
			require.Equal(t, validPlant.GetPhotos().Len(), 1)
			frepo.AssertExpectations(t)
			prepo.AssertExpectations(t)
		})
	})

	t.Run("Photo addition error during upload", func(t provider.T) {
		t.Parallel()

		asvc, ctx := setupAuthService(t, validOwnerID, true)

		fdata := models.FileData{}
		description := "test photo"

		frepo := new(MockFileRepository)
		t.WithNewStep("Setup file upload", func(pctx provider.StepCtx) {
			frepo.On("Upload", mock.Anything, &fdata).Return(&models.File{ID: validFileID}, nil)
		})

		prepo := new(MockPlantRepository)
		t.WithNewStep("Setup plant update error", func(pctx provider.StepCtx) {
			prepo.On("Update", mock.Anything, validPlantID, mock.Anything).Return(nil, assert.AnError)
		})

		crepo := new(MockPlantCategoryRepository)

		svc := plantservice.NewPlantService(prepo, crepo, frepo, asvc)

		t.WithNewStep("Attempt upload with update error", func(pctx provider.StepCtx) {
			err := svc.UploadPlantPhoto(ctx, validPlantID, fdata, description)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			frepo.AssertExpectations(t)
			prepo.AssertExpectations(t)

		})
	})
}

func TestPlantService(t *testing.T) {
	suite.RunSuite(t, new(PlantServiceTestSuite))
}
