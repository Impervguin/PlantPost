//go:build unit

package albumservice_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"PlantSite/internal/models/album"
	"PlantSite/internal/models/auth"
	albumservice "PlantSite/internal/services/album-service"
	authservice "PlantSite/internal/services/auth-service"
	authmock "PlantSite/internal/services/auth-service/auth-mock"
	"PlantSite/internal/utils/logs"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// AlbumBuilder реализует паттерн Data Builder для создания тестовых альбомов
type AlbumBuilder struct {
	name        string
	description string
	plantIDs    uuid.UUIDs
	ownerID     uuid.UUID
}

func NewAlbumBuilder() *AlbumBuilder {
	return &AlbumBuilder{
		name:        "Test Album",
		description: "Test Description",
		plantIDs:    uuid.UUIDs{},
		ownerID:     uuid.New(),
	}
}

func (b *AlbumBuilder) WithName(name string) *AlbumBuilder {
	b.name = name
	return b
}

func (b *AlbumBuilder) WithDescription(description string) *AlbumBuilder {
	b.description = description
	return b
}

func (b *AlbumBuilder) WithPlantIDs(plantIDs uuid.UUIDs) *AlbumBuilder {
	b.plantIDs = plantIDs
	return b
}

func (b *AlbumBuilder) WithOwnerID(ownerID uuid.UUID) *AlbumBuilder {
	b.ownerID = ownerID
	return b
}

func (b *AlbumBuilder) Build() (*album.Album, error) {
	return album.NewAlbum(b.name, b.description, b.plantIDs, b.ownerID)
}

// MockAlbumRepository implements AlbumRepository interface
type MockAlbumRepository struct {
	mock.Mock
}

func (m *MockAlbumRepository) Create(ctx context.Context, alb *album.Album) (*album.Album, error) {
	args := m.Called(ctx, alb)
	return args.Get(0).(*album.Album), args.Error(1)
}

func (m *MockAlbumRepository) Get(ctx context.Context, id uuid.UUID) (*album.Album, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*album.Album), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAlbumRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(*album.Album) (*album.Album, error)) (*album.Album, error) {
	args := m.Called(ctx, id, updateFn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	alb, err := updateFn(args.Get(0).(*album.Album))
	return alb, err
}

func (m *MockAlbumRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAlbumRepository) List(ctx context.Context, ownerID uuid.UUID) ([]*album.Album, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).([]*album.Album), args.Error(1)
}

type AlbumServiceTestSuite struct {
	suite.Suite
}

func (s *AlbumServiceTestSuite) BeforeAll(t provider.T) {
	logs.InitNoopLogger()
}

func (s *AlbumServiceTestSuite) BeforeEach(t provider.T) {
	t.Epic("Album Service")
	t.Feature("Album Management")
}

func setupAuthService(t provider.T, userID uuid.UUID, hasMemberRights bool) (*authservice.AuthService, context.Context) {
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
	user.On("HasMemberRights").Return(hasMemberRights)

	sessions.On("Get", mock.Anything, sessionID).Return(validSession, nil)
	authCtx := asvc.Authenticate(context.Background(), sessionID)
	arepo.On("Get", authCtx, userID).Return(user, nil)

	return asvc, authCtx
}

func (s *AlbumServiceTestSuite) TestCreateAlbum(t provider.T) {
	t.Tags("creation", "positive")
	t.Description("Test album creation functionality")
	t.Parallel()

	ownerID := uuid.New()
	asvc, ctx := setupAuthService(t, ownerID, true)

	albumBuilder := NewAlbumBuilder().WithOwnerID(ownerID)
	validAlbum, err := albumBuilder.Build()
	require.NoError(t, err)

	t.Run("Successful album creation", func(t provider.T) {
		t.Parallel()

		repo := new(MockAlbumRepository)
		repo.On("Create", mock.Anything, validAlbum).Return(validAlbum, nil)

		svc := albumservice.NewAlbumService(repo, asvc)

		var result *album.Album
		t.WithNewStep("Create album", func(pctx provider.StepCtx) {
			result, err = svc.CreateAlbum(ctx, validAlbum)
		})

		t.WithNewStep("Verify creation result", func(pctx provider.StepCtx) {
			require.NoError(t, err)
			assert.Equal(t, validAlbum, result)
			repo.AssertExpectations(t)
		})
	})

	t.Run("Not authorized for album creation", func(t provider.T) {
		t.Parallel()

		arepo := new(authmock.MockAuthRepository)
		sessions := new(authmock.MockSessionStorage)
		hasher := new(authmock.MockPasswdHasher)
		asvc := authservice.NewAuthService(sessions, arepo, hasher)
		sessions.On("Get", mock.Anything, mock.Anything).Return(nil, assert.AnError)

		albumBuilder := NewAlbumBuilder()
		validAlbum, err := albumBuilder.Build()
		require.NoError(t, err)

		repo := new(MockAlbumRepository)
		svc := albumservice.NewAlbumService(repo, asvc)

		var result *album.Album
		t.WithNewStep("Attempt unauthorized creation", func(pctx provider.StepCtx) {
			result, err = svc.CreateAlbum(context.Background(), validAlbum)
		})

		t.WithNewStep("Verify authorization error", func(pctx provider.StepCtx) {
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("No member rights for album creation", func(t provider.T) {
		t.Parallel()

		ownerID := uuid.New()
		asvc, ctx := setupAuthService(t, ownerID, false)

		albumBuilder := NewAlbumBuilder().WithOwnerID(ownerID)
		validAlbum, err := albumBuilder.Build()
		require.NoError(t, err)

		repo := new(MockAlbumRepository)
		svc := albumservice.NewAlbumService(repo, asvc)

		var result *album.Album
		t.WithNewStep("Attempt creation without member rights", func(pctx provider.StepCtx) {
			result, err = svc.CreateAlbum(ctx, validAlbum)
		})

		t.WithNewStep("Verify member rights error", func(pctx provider.StepCtx) {
			assert.ErrorIs(t, err, auth.ErrNoMemberRights)
			assert.Nil(t, result)
		})
	})
}

func (s *AlbumServiceTestSuite) TestGetAlbum(t provider.T) {
	t.Tags("retrieval", "positive")
	t.Description("Test album retrieval functionality")
	t.Parallel()

	t.Run("Successful album retrieval", func(t provider.T) {
		t.Parallel()

		ownerID := uuid.New()
		albumID := uuid.New()
		asvc, ctx := setupAuthService(t, ownerID, true)

		albumBuilder := NewAlbumBuilder().WithOwnerID(ownerID)
		validAlbum, err := albumBuilder.Build()
		require.NoError(t, err)

		repo := new(MockAlbumRepository)
		repo.On("Get", mock.Anything, albumID).Return(validAlbum, nil)

		svc := albumservice.NewAlbumService(repo, asvc)

		var result *album.Album
		t.WithNewStep("Retrieve album", func(pctx provider.StepCtx) {
			result, err = svc.GetAlbum(ctx, albumID)
		})

		t.WithNewStep("Verify retrieval result", func(pctx provider.StepCtx) {
			require.NoError(t, err)
			assert.Equal(t, validAlbum, result)
		})
	})

	t.Run("Not owner trying to access album", func(t provider.T) {
		t.Parallel()

		ownerID := uuid.New()
		otherUserID := uuid.New()
		albumID := uuid.New()
		asvc, ctx := setupAuthService(t, otherUserID, true)

		albumBuilder := NewAlbumBuilder().WithOwnerID(ownerID)
		validAlbum, err := albumBuilder.Build()
		require.NoError(t, err)

		repo := new(MockAlbumRepository)
		repo.On("Get", mock.Anything, albumID).Return(validAlbum, nil)

		svc := albumservice.NewAlbumService(repo, asvc)

		var result *album.Album
		t.WithNewStep("Attempt access by non-owner", func(pctx provider.StepCtx) {
			result, err = svc.GetAlbum(ctx, albumID)
		})

		t.WithNewStep("Verify ownership error", func(pctx provider.StepCtx) {
			assert.ErrorIs(t, err, albumservice.ErrNotOwner)
			assert.Nil(t, result)
		})
	})

	t.Run("Album not found during retrieval", func(t provider.T) {
		t.Parallel()
		var err error

		ownerID := uuid.New()
		albumID := uuid.New()
		asvc, ctx := setupAuthService(t, ownerID, true)

		repo := new(MockAlbumRepository)
		repo.On("Get", mock.Anything, albumID).Return(nil, errors.New("not found"))

		svc := albumservice.NewAlbumService(repo, asvc)

		var result *album.Album
		t.WithNewStep("Attempt to retrieve non-existent album", func(pctx provider.StepCtx) {
			result, err = svc.GetAlbum(ctx, albumID)
		})

		t.WithNewStep("Verify not found error", func(pctx provider.StepCtx) {
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}

func (s *AlbumServiceTestSuite) TestUpdateAlbumName(t provider.T) {
	t.Tags("update", "positive")
	t.Description("Test album name update functionality")
	t.Parallel()

	t.Run("Successful album name update", func(t provider.T) {
		t.Parallel()

		ownerID := uuid.New()
		albumID := uuid.New()
		newName := "Updated Album Name"
		asvc, ctx := setupAuthService(t, ownerID, true)

		albumBuilder := NewAlbumBuilder().WithOwnerID(ownerID)
		validAlbum, err := albumBuilder.Build()
		require.NoError(t, err)

		repo := new(MockAlbumRepository)
		repo.On("Update", mock.Anything, albumID, mock.Anything).Return(validAlbum, nil)

		svc := albumservice.NewAlbumService(repo, asvc)

		t.WithNewStep("Update album name", func(pctx provider.StepCtx) {
			err = svc.UpdateAlbumName(ctx, albumID, newName)
		})

		t.WithNewStep("Verify update success", func(pctx provider.StepCtx) {
			require.NoError(t, err)
		})
	})

	t.Run("Not owner trying to update album name", func(t provider.T) {
		t.Parallel()

		ownerID := uuid.New()
		otherUserID := uuid.New()
		albumID := uuid.New()
		newName := "Updated Album Name"
		asvc, ctx := setupAuthService(t, otherUserID, true)

		albumBuilder := NewAlbumBuilder().WithOwnerID(ownerID)
		validAlbum, err := albumBuilder.Build()
		require.NoError(t, err)

		repo := new(MockAlbumRepository)
		repo.On("Update", mock.Anything, albumID, mock.Anything).Return(validAlbum, nil)

		svc := albumservice.NewAlbumService(repo, asvc)

		t.WithNewStep("Attempt update by non-owner", func(pctx provider.StepCtx) {
			err = svc.UpdateAlbumName(ctx, albumID, newName)
		})

		t.WithNewStep("Verify ownership error", func(pctx provider.StepCtx) {
			assert.ErrorIs(t, err, albumservice.ErrNotOwner)
		})
	})
}

func (s *AlbumServiceTestSuite) TestAddPlantToAlbum(t provider.T) {
	t.Tags("update", "plants")
	t.Description("Test adding plants to album functionality")
	t.Parallel()

	t.Run("Successful plant addition to album", func(t provider.T) {
		t.Parallel()

		ownerID := uuid.New()
		albumID := uuid.New()
		plantID := uuid.New()
		asvc, ctx := setupAuthService(t, ownerID, true)

		albumBuilder := NewAlbumBuilder().WithOwnerID(ownerID)
		validAlbum, err := albumBuilder.Build()
		require.NoError(t, err)

		repo := new(MockAlbumRepository)
		repo.On("Update", mock.Anything, albumID, mock.Anything).Return(validAlbum, nil)

		svc := albumservice.NewAlbumService(repo, asvc)

		t.WithNewStep("Add plant to album", func(pctx provider.StepCtx) {
			err = svc.AddPlantToAlbum(ctx, albumID, plantID)
		})

		t.WithNewStep("Verify plant addition success", func(pctx provider.StepCtx) {
			require.NoError(t, err)
		})
	})

	t.Run("Plant already exists in album", func(t provider.T) {
		t.Parallel()

		ownerID := uuid.New()
		albumID := uuid.New()
		plantID := uuid.New()
		asvc, ctx := setupAuthService(t, ownerID, true)

		albumBuilder := NewAlbumBuilder().WithOwnerID(ownerID).WithPlantIDs(uuid.UUIDs{plantID})
		albumWithPlant, err := albumBuilder.Build()
		require.NoError(t, err)

		repo := new(MockAlbumRepository)
		repo.On("Update", mock.Anything, albumID, mock.Anything).Return(albumWithPlant, nil)

		svc := albumservice.NewAlbumService(repo, asvc)

		t.WithNewStep("Attempt to add duplicate plant", func(pctx provider.StepCtx) {
			err = svc.AddPlantToAlbum(ctx, albumID, plantID)
		})

		t.WithNewStep("Verify duplicate plant error", func(pctx provider.StepCtx) {
			assert.ErrorIs(t, err, album.ErrPlantAlreadyInAlbum)
		})
	})
}

func (s *AlbumServiceTestSuite) TestRemovePlantFromAlbum(t provider.T) {
	t.Tags("update", "plants")
	t.Description("Test removing plants from album functionality")
	t.Parallel()

	t.Run("Successful plant removal from album", func(t provider.T) {
		t.Parallel()

		ownerID := uuid.New()
		albumID := uuid.New()
		plantID := uuid.New()
		asvc, ctx := setupAuthService(t, ownerID, true)

		albumBuilder := NewAlbumBuilder().WithOwnerID(ownerID).WithPlantIDs(uuid.UUIDs{plantID})
		validAlbum, err := albumBuilder.Build()
		require.NoError(t, err)

		repo := new(MockAlbumRepository)
		repo.On("Update", mock.Anything, albumID, mock.Anything).Return(validAlbum, nil)

		svc := albumservice.NewAlbumService(repo, asvc)

		t.WithNewStep("Remove plant from album", func(pctx provider.StepCtx) {
			err = svc.RemovePlantFromAlbum(ctx, albumID, plantID)
		})

		t.WithNewStep("Verify plant removal success", func(pctx provider.StepCtx) {
			require.NoError(t, err)
		})
	})

	t.Run("Plant not found in album during removal", func(t provider.T) {
		t.Parallel()

		ownerID := uuid.New()
		albumID := uuid.New()
		plantID := uuid.New()
		asvc, ctx := setupAuthService(t, ownerID, true)

		albumBuilder := NewAlbumBuilder().WithOwnerID(ownerID).WithPlantIDs(uuid.UUIDs{})
		emptyAlbum, err := albumBuilder.Build()
		require.NoError(t, err)

		repo := new(MockAlbumRepository)
		repo.On("Update", mock.Anything, albumID, mock.Anything).Return(emptyAlbum, nil)

		svc := albumservice.NewAlbumService(repo, asvc)

		t.WithNewStep("Attempt to remove non-existent plant", func(pctx provider.StepCtx) {
			err = svc.RemovePlantFromAlbum(ctx, albumID, plantID)
		})

		t.WithNewStep("Verify plant not found error", func(pctx provider.StepCtx) {
			assert.ErrorIs(t, err, album.ErrPlantNotFound)
		})
	})
}

func (s *AlbumServiceTestSuite) TestListAlbums(t provider.T) {
	t.Tags("retrieval", "listing")
	t.Description("Test album listing functionality")
	t.Parallel()

	t.Run("Successful album listing", func(t provider.T) {
		t.Parallel()

		ownerID := uuid.New()
		asvc, ctx := setupAuthService(t, ownerID, true)

		albumBuilder := NewAlbumBuilder().WithOwnerID(ownerID)
		validAlbum, err := albumBuilder.Build()
		require.NoError(t, err)

		repo := new(MockAlbumRepository)
		expectedAlbums := []*album.Album{validAlbum}
		repo.On("List", mock.Anything, ownerID).Return(expectedAlbums, nil)

		svc := albumservice.NewAlbumService(repo, asvc)

		var result []*album.Album
		t.WithNewStep("List albums", func(ppctx provider.StepCtx) {
			result, err = svc.ListAlbums(ctx)
		})

		t.WithNewStep("Verify album list", func(pctx provider.StepCtx) {
			require.NoError(t, err)
			assert.Equal(t, expectedAlbums, result)
		})
	})

	t.Run("Empty album list", func(t provider.T) {
		t.Parallel()
		var err error

		ownerID := uuid.New()
		asvc, ctx := setupAuthService(t, ownerID, true)

		repo := new(MockAlbumRepository)
		repo.On("List", mock.Anything, ownerID).Return([]*album.Album{}, nil)

		svc := albumservice.NewAlbumService(repo, asvc)

		var result []*album.Album
		t.WithNewStep("List empty albums", func(ppctx provider.StepCtx) {
			result, err = svc.ListAlbums(ctx)
		})

		t.WithNewStep("Verify empty list", func(pctx provider.StepCtx) {
			require.NoError(t, err)
			assert.Empty(t, result)
		})
	})
}

func TestAlbumService(t *testing.T) {
	suite.RunSuite(t, new(AlbumServiceTestSuite))
}
