package albumservice_test

import (
	"context"
	"errors"
	"testing"

	"PlantSite/internal/models/album"
	"PlantSite/internal/models/auth"
	albumservice "PlantSite/internal/services/album-service"
	authservice "PlantSite/internal/services/auth-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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

type MockUser struct {
	mock.Mock
}

func (m *MockUser) ID() uuid.UUID {
	return m.Called().Get(0).(uuid.UUID)
}

func (m *MockUser) HasMemberRights() bool {
	return m.Called().Bool(0)
}

func (m *MockUser) Auth(_ []byte, _ func(hashPasswd []byte, plainPasswd []byte) (bool, error)) bool {
	return m.Called().Bool(0)
}

func (m *MockUser) HasAuthorRights() bool {
	return m.Called().Bool(0)
}

func PutUserInContext(ctx context.Context, user auth.User) context.Context {
	return context.WithValue(ctx, authservice.AuthContextKey, user)
}

func TestAlbumService(t *testing.T) {
	ctx := context.Background()
	validAlbumID := uuid.New()
	validOwnerID := uuid.New()
	validPlantID := uuid.New()
	validAlbum, err := album.NewAlbum("Test Album", "Test Description", uuid.UUIDs{}, validOwnerID)
	require.NoError(t, err)

	t.Run("CreateAlbum", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(validOwnerID)
			user.On("HasMemberRights").Return(true)

			repo.On("Create", mock.Anything, mock.Anything).Return(validAlbum, nil)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			result, err := svc.CreateAlbum(ctx, validAlbum)
			require.NoError(t, err)
			assert.Equal(t, validAlbum, result)

			repo.AssertExpectations(t)
			user.AssertExpectations(t)
		})

		t.Run("NotAuthorized", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			svc := albumservice.NewAlbumService(repo)

			_, err := svc.CreateAlbum(ctx, validAlbum)
			assert.ErrorIs(t, err, albumservice.ErrNotAuthorized)
		})

		t.Run("NotMember", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("HasMemberRights").Return(false)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			_, err := svc.CreateAlbum(ctx, validAlbum)
			assert.ErrorIs(t, err, albumservice.ErrNotMember)
		})
	})

	t.Run("GetAlbum", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(validOwnerID)
			user.On("HasMemberRights").Return(true)

			repo.On("Get", mock.Anything, validAlbumID).Return(validAlbum, nil)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			result, err := svc.GetAlbum(ctx, validAlbumID)
			require.NoError(t, err)
			assert.Equal(t, validAlbum, result)

			repo.AssertExpectations(t)
			user.AssertExpectations(t)
		})

		t.Run("NotOwner", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(uuid.New()) // Different owner
			user.On("HasMemberRights").Return(true)

			repo.On("Get", mock.Anything, validAlbumID).Return(validAlbum, nil)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			_, err := svc.GetAlbum(ctx, validAlbumID)
			assert.ErrorIs(t, err, albumservice.ErrNotOwner)
		})

		t.Run("NotFound", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(validOwnerID)
			user.On("HasMemberRights").Return(true)

			repo.On("Get", mock.Anything, validAlbumID).Return(nil, errors.New("not found"))

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			_, err := svc.GetAlbum(ctx, validAlbumID)
			assert.Error(t, err)
		})
	})

	t.Run("UpdateAlbumName", func(t *testing.T) {
		newName := "New Album Name"
		newAlb := validAlbum
		err := newAlb.UpdateName(newName)
		require.NoError(t, err)

		t.Run("Success", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(validOwnerID)
			user.On("HasMemberRights").Return(true)

			repo.On("Update", mock.Anything, validAlbumID, mock.Anything).Return(validAlbum, nil)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			err := svc.UpdateAlbumName(ctx, validAlbumID, newName)
			require.NoError(t, err)

			repo.AssertExpectations(t)
			user.AssertExpectations(t)
		})

		t.Run("NotOwner", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(uuid.New())
			user.On("HasMemberRights").Return(true)

			repo.On("Update", mock.Anything, validAlbumID, mock.Anything).Return(validAlbum, nil)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			err := svc.UpdateAlbumName(ctx, validAlbumID, newName)
			assert.ErrorIs(t, err, albumservice.ErrNotOwner)
		})
	})

	t.Run("UpdateAlbumDescription", func(t *testing.T) {
		newDesc := "New Album Description"

		t.Run("Success", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(validOwnerID)
			user.On("HasMemberRights").Return(true)

			repo.On("Update", mock.Anything, validAlbumID, mock.Anything).Return(validAlbum, nil)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			err := svc.UpdateAlbumDescription(ctx, validAlbumID, newDesc)
			require.NoError(t, err)

			repo.AssertExpectations(t)
			user.AssertExpectations(t)
		})
	})

	t.Run("AddPlantToAlbum", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(validOwnerID)
			user.On("HasMemberRights").Return(true)

			repo.On("Update", mock.Anything, validAlbumID, mock.Anything).Return(validAlbum, nil)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			err := svc.AddPlantToAlbum(ctx, validAlbumID, validPlantID)
			require.NoError(t, err)

			repo.AssertExpectations(t)
			user.AssertExpectations(t)
		})

		t.Run("PlantAlreadyExists", func(t *testing.T) {
			albumWithPlant, err := album.NewAlbum("Test", "Desc", uuid.UUIDs{validPlantID}, validOwnerID)
			require.NoError(t, err)

			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(validOwnerID)
			user.On("HasMemberRights").Return(true)

			repo.On("Update", mock.Anything, validAlbumID, mock.Anything).Return(albumWithPlant, nil)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			err = svc.AddPlantToAlbum(ctx, validAlbumID, validPlantID)
			assert.ErrorIs(t, err, album.ErrPlantAlreadyInAlbum)
		})
	})

	t.Run("RemovePlantFromAlbum", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(validOwnerID)
			user.On("HasMemberRights").Return(true)

			repo.On("Update", mock.Anything, validAlbumID, mock.Anything).Return(validAlbum, nil)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			err := svc.RemovePlantFromAlbum(ctx, validAlbumID, validPlantID)
			require.NoError(t, err)

			repo.AssertExpectations(t)
			user.AssertExpectations(t)
		})

		t.Run("PlantNotInAlbum", func(t *testing.T) {
			emptyAlbum, err := album.NewAlbum("Test", "Desc", uuid.UUIDs{}, validOwnerID)
			require.NoError(t, err)

			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(validOwnerID)
			user.On("HasMemberRights").Return(true)

			repo.On("Update", mock.Anything, validAlbumID, mock.Anything).Return(emptyAlbum, nil)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			err = svc.RemovePlantFromAlbum(ctx, validAlbumID, validPlantID)
			assert.ErrorIs(t, err, album.ErrPlantNotFound)
		})
	})

	t.Run("DeleteAlbum", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(validOwnerID)
			user.On("HasMemberRights").Return(true)

			repo.On("Get", mock.Anything, validAlbumID).Return(validAlbum, nil)
			repo.On("Delete", mock.Anything, validAlbumID).Return(nil)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			err := svc.DeleteAlbum(ctx, validAlbumID)
			require.NoError(t, err)

			repo.AssertExpectations(t)
			user.AssertExpectations(t)
		})

		t.Run("NotOwner", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(uuid.New()) // Different owner
			user.On("HasMemberRights").Return(true)

			repo.On("Get", mock.Anything, validAlbumID).Return(validAlbum, nil)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			err := svc.DeleteAlbum(ctx, validAlbumID)
			assert.ErrorIs(t, err, albumservice.ErrNotOwner)
		})
	})

	t.Run("ListAlbums", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(validOwnerID)
			user.On("HasMemberRights").Return(true)

			expectedAlbums := []*album.Album{validAlbum}
			repo.On("List", mock.Anything, validOwnerID).Return(expectedAlbums, nil)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			result, err := svc.ListAlbums(ctx)
			require.NoError(t, err)
			assert.Equal(t, expectedAlbums, result)

			repo.AssertExpectations(t)
			user.AssertExpectations(t)
		})

		t.Run("EmptyList", func(t *testing.T) {
			repo := new(MockAlbumRepository)
			user := new(MockUser)
			user.On("ID").Return(validOwnerID)
			user.On("HasMemberRights").Return(true)

			repo.On("List", mock.Anything, validOwnerID).Return([]*album.Album{}, nil)

			svc := albumservice.NewAlbumService(repo)
			ctx := PutUserInContext(ctx, user)

			result, err := svc.ListAlbums(ctx)
			require.NoError(t, err)
			assert.Empty(t, result)
		})
	})
}
