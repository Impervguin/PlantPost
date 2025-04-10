package searchservice_test

import (
	"context"
	"testing"

	"PlantSite/internal/models"
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/post"
	"PlantSite/internal/models/search"
	searchservice "PlantSite/internal/services/search-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

func TestNewSearchService(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)
		assert.NotNil(t, svc)
	})

	t.Run("NilDependencies", func(t *testing.T) {
		assert.Panics(t, func() {
			searchservice.NewSearchService(nil, new(MockFileRepository), new(MockFileRepository))
		})
		assert.Panics(t, func() {
			searchservice.NewSearchService(new(MockSearchRepository), nil, new(MockFileRepository))
		})
		assert.Panics(t, func() {
			searchservice.NewSearchService(new(MockSearchRepository), new(MockFileRepository), nil)
		})
	})
}
