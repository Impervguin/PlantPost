//go:build unit

package searchservice_test

import (
	"context"

	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/post"
	"PlantSite/internal/models/search"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockSearchRepository struct {
	mock.Mock
}

func (m *MockSearchRepository) SearchPosts(ctx context.Context, search *search.PostSearch) ([]*post.Post, error) {
	args := m.Called(ctx, search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	res := make([]*post.Post, 0)
	for _, post := range args.Get(0).([]*post.Post) {
		if search.Filter(post) {
			res = append(res, post)
		}
	}

	return res, args.Error(1)
}

func (m *MockSearchRepository) SearchPlants(ctx context.Context, search *search.PlantSearch) ([]*plant.Plant, error) {
	args := m.Called(ctx, search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	res := make([]*plant.Plant, 0)
	for _, plnt := range args.Get(0).([]*plant.Plant) {
		if search.Filter(plnt) {
			res = append(res, plnt)
		}
	}

	return res, args.Error(1)
}

func (m *MockSearchRepository) GetPostByID(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*post.Post), args.Error(1)
}

func (m *MockSearchRepository) GetPlantByID(ctx context.Context, id uuid.UUID) (*plant.Plant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*plant.Plant), args.Error(1)
}

func (m *MockSearchRepository) GetPostAuthors(ctx context.Context) ([]*auth.Author, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auth.Author), args.Error(1)
}

func (m *MockSearchRepository) GetPostTags(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
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

type PlantBuilder struct {
	name        string
	latinName   string
	description string
	mainPhotoID uuid.UUID
	photos      *plant.PlantPhotos
	category    string
	spec        *MockPlantSpecification
}

func NewPlantBuilder() *PlantBuilder {
	return &PlantBuilder{
		name:        "Rose",
		latinName:   "Rosa",
		description: "Beautiful flower",
		mainPhotoID: uuid.New(),
		photos:      plant.NewPlantPhotos(),
		category:    "flowers",
		spec:        &MockPlantSpecification{},
	}
}

func (b *PlantBuilder) WithName(name string) *PlantBuilder {
	b.name = name
	return b
}

func (b *PlantBuilder) WithLatinName(latinName string) *PlantBuilder {
	b.latinName = latinName
	return b
}

func (b *PlantBuilder) WithDescription(description string) *PlantBuilder {
	b.description = description
	return b
}

func (b *PlantBuilder) WithMainPhotoID(photoID uuid.UUID) *PlantBuilder {
	b.mainPhotoID = photoID
	return b
}

func (b *PlantBuilder) WithCategory(category string) *PlantBuilder {
	b.category = category
	return b
}

func (b *PlantBuilder) WithPhoto(fileID uuid.UUID, description string) *PlantBuilder {
	photo, err := plant.NewPlantPhoto(fileID, description)
	if err == nil {
		b.photos.Add(photo)
	}
	return b
}

func (b *PlantBuilder) WithNoPhotos() *PlantBuilder {
	b.photos = plant.NewPlantPhotos()
	return b
}

type PostBuilder struct {
	title    string
	content  *post.Content
	tags     []string
	authorID uuid.UUID
	photos   *post.PostPhotos
}

func NewPostBuilder() *PostBuilder {
	content, _ := post.NewContent("Test content", post.ContentTypePlainText)
	return &PostBuilder{
		title:    "Test Post",
		content:  content,
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
	return post.NewPost(
		b.title,
		*b.content,
		b.tags,
		b.authorID,
		b.photos,
	)
}

func (b *PlantBuilder) Build() (*plant.Plant, error) {
	b.spec.On("Validate").Return(nil)
	b.spec.On("Category").Return(b.category)

	return plant.NewPlant(
		b.name,
		b.latinName,
		b.description,
		b.mainPhotoID,
		*b.photos,
		b.category,
		b.spec,
	)
}
