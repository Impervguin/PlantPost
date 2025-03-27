package postservice_test

import (
	"context"

	authservice "PlantSite/internal/auth-service"
	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/post"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockPostRepository implements post.PostRepository interface
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(ctx context.Context, p *post.Post) (*post.Post, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*post.Post), args.Error(1)
}

func (m *MockPostRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(*post.Post) (*post.Post, error)) (*post.Post, error) {
	args := m.Called(ctx, id, updateFn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	pst, err := updateFn(args.Get(0).(*post.Post))
	return pst, err
}

func (m *MockPostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostRepository) Get(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*post.Post), args.Error(1)
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

// MockUser implements auth.User interface
type MockUser struct {
	mock.Mock
}

func (m *MockUser) ID() uuid.UUID {
	return m.Called().Get(0).(uuid.UUID)
}

func (m *MockUser) Name() string {
	return m.Called().String(0)
}

func (m *MockUser) Email() string {
	return m.Called().String(0)
}

func (m *MockUser) HashedPassword() []byte {
	return m.Called().Get(0).([]byte)
}

func (m *MockUser) Auth(password []byte, compareFn func(hashPasswd, plainPasswd []byte) (bool, error)) bool {
	args := m.Called(password, compareFn)
	return args.Bool(0)
}

func (m *MockUser) HasMemberRights() bool {
	return m.Called().Bool(0)
}

func (m *MockUser) HasAuthorRights() bool {
	return m.Called().Bool(0)
}

func PutUserInContext(ctx context.Context, user auth.User) context.Context {
	return context.WithValue(ctx, authservice.AuthContextKey, user)
}
