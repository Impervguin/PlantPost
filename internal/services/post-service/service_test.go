//go:build unit

package postservice_test

import (
	"context"
	"time"

	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	authservice "PlantSite/internal/services/auth-service"
	authmock "PlantSite/internal/services/auth-service/auth-mock"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
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
	user.On("HasAuthorRights").Return(hasAuthorRights)
	user.On("ID").Return(userID)

	sessions.On("Get", mock.Anything, sessionID).Return(validSession, nil)
	authCtx := asvc.Authenticate(context.Background(), sessionID)
	arepo.On("Get", authCtx, userID).Return(user, nil)

	return asvc, authCtx
}

type PostBuilder struct {
	title    string
	content  *post.Content
	tags     []string
	authorID uuid.UUID
	photos   *post.PostPhotos
}

func NewPostBuilder() *PostBuilder {
	return &PostBuilder{
		title:    "Test Post",
		tags:     []string{"tag1", "tag2"},
		authorID: uuid.New(),
		photos:   post.NewPostPhotos(),
	}
}

func (b *PostBuilder) WithTitle(title string) *PostBuilder {
	b.title = title
	return b
}

func (b *PostBuilder) WithContent(content string, contentType post.ContentFormat) *PostBuilder {
	postContent, err := post.NewContent(content, contentType)
	if err == nil {
		b.content = postContent
	}
	return b
}

func (b *PostBuilder) WithTags(tags []string) *PostBuilder {
	b.tags = tags
	return b
}

func (b *PostBuilder) WithAuthorID(authorID uuid.UUID) *PostBuilder {
	b.authorID = authorID
	return b
}

func (b *PostBuilder) WithPhoto(fileID uuid.UUID, placeNumber int) *PostBuilder {
	photo, err := post.NewPostPhoto(fileID, placeNumber)
	if err == nil {
		b.photos.Add(photo)
	}
	return b
}

func (b *PostBuilder) WithNoPhotos() *PostBuilder {
	b.photos = post.NewPostPhotos()
	return b
}

func (b *PostBuilder) Build() (*post.Post, error) {
	if b.content == nil {
		content, err := post.NewContent("Test content", post.ContentTypePlainText)
		if err != nil {
			return nil, err
		}
		b.content = content
	}

	return post.NewPost(
		b.title,
		*b.content,
		b.tags,
		b.authorID,
		b.photos,
	)
}
