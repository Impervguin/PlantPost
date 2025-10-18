//go:build unit

package searchservice_test

import (
	"context"
	"testing"

	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	"PlantSite/internal/models/search"
	searchservice "PlantSite/internal/services/search-service"
	"PlantSite/internal/utils/logs"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type SearchServiceTestSearchTestSuite struct {
	suite.Suite
}

func (s *SearchServiceTestSearchTestSuite) BeforeAll(t provider.T) {
	logs.InitNoopLogger()
}

func (s *SearchServiceTestSearchTestSuite) BeforeEach(t provider.T) {
	t.Epic("Search Service")
	t.Feature("Search Management")
}

type SearchPostResultBuilder struct {
	post       *post.Post
	photoFiles map[uuid.UUID]*models.File
}

func NewSearchPostResultBuilder(post *post.Post) *SearchPostResultBuilder {
	return &SearchPostResultBuilder{
		post:       post,
		photoFiles: make(map[uuid.UUID]*models.File),
	}
}

func (b *SearchPostResultBuilder) WithPhotoFile(fileID uuid.UUID, fileName string) *SearchPostResultBuilder {
	b.photoFiles[fileID] = &models.File{ID: fileID, Name: fileName}
	return b
}

func (b *SearchPostResultBuilder) Build() *searchservice.SearchPost {
	photos := make([]searchservice.SearchPostPhoto, 0)

	for _, photo := range b.post.Photos().List() {
		if file, exists := b.photoFiles[photo.FileID()]; exists {
			photos = append(photos, searchservice.SearchPostPhoto{
				ID:          photo.ID(),
				PlaceNumber: photo.PlaceNumber(),
				File:        *file,
			})
		}
	}

	return &searchservice.SearchPost{
		ID:        b.post.ID(),
		Title:     b.post.Title(),
		Content:   b.post.Content(),
		Tags:      b.post.Tags(),
		AuthorID:  b.post.AuthorID(),
		CreatedAt: b.post.CreatedAt(),
		UpdatedAt: b.post.UpdatedAt(),
		Photos:    photos,
	}
}

func (s *SearchServiceTestSearchTestSuite) TestSearchPosts(t provider.T) {
	t.Tags("search", "posts")
	t.Description("Test post search functionality")
	t.Parallel()

	ctx := context.Background()

	t.Run("Successful post search", func(t provider.T) {
		t.Parallel()

		validUserID := uuid.New()
		validPost, err := NewPostBuilder().
			WithAuthorID(validUserID).
			WithPhoto(uuid.New(), 1).
			WithPhoto(uuid.New(), 2).
			Build()
		require.NoError(t, err)

		photo1 := validPost.Photos().List()[0]
		photo2 := validPost.Photos().List()[1]
		photoFile1 := &models.File{ID: photo1.FileID(), Name: "photo1.jpg"}
		photoFile2 := &models.File{ID: photo2.FileID(), Name: "photo2.jpg"}

		resultBuilder := NewSearchPostResultBuilder(validPost).
			WithPhotoFile(photo1.FileID(), "photo1.jpg").
			WithPhotoFile(photo2.FileID(), "photo2.jpg")

		expectedResult := resultBuilder.Build()

		srch := search.NewPostSearch()
		srch.AddFilter(search.NewPostAuthorFilter(validUserID))

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup post search", func(pctx provider.StepCtx) {
			srepo.On("SearchPosts", mock.Anything, srch).Return([]*post.Post{validPost}, nil)
		})

		ptfrepo := new(MockFileRepository)
		t.WithNewStep("Setup file repository", func(pctx provider.StepCtx) {
			ptfrepo.On("Get", mock.Anything, photo1.FileID()).Return(photoFile1, nil)
			ptfrepo.On("Get", mock.Anything, photo2.FileID()).Return(photoFile2, nil)
		})

		pfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		var results []*searchservice.SearchPost
		t.WithNewStep("Search posts", func(pctx provider.StepCtx) {
			var err error
			results, err = svc.SearchPosts(ctx, srch)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify results", func(pctx provider.StepCtx) {
			require.Len(t, results, 1)
			assert.Equal(t, expectedResult, results[0])
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			srepo.AssertExpectations(t)
			ptfrepo.AssertExpectations(t)
		})
	})

	t.Run("Empty search results", func(t provider.T) {
		t.Parallel()

		srch := search.NewPostSearch()
		srch.AddFilter(search.NewPostAuthorFilter(uuid.New()))

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup empty search results", func(pctx provider.StepCtx) {
			srepo.On("SearchPosts", mock.Anything, srch).Return([]*post.Post{}, nil)
		})

		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		var results []*searchservice.SearchPost
		t.WithNewStep("Search posts with no results", func(pctx provider.StepCtx) {
			var err error
			results, err = svc.SearchPosts(ctx, srch)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify empty results", func(pctx provider.StepCtx) {
			assert.Empty(t, results)
		})
	})

	t.Run("Repository error during search", func(t provider.T) {
		t.Parallel()

		srch := search.NewPostSearch()

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup repository error", func(pctx provider.StepCtx) {
			srepo.On("SearchPosts", ctx, srch).Return(nil, assert.AnError)
		})

		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		t.WithNewStep("Attempt search with repository error", func(pctx provider.StepCtx) {
			_, err := svc.SearchPosts(ctx, srch)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Photo file not found during search", func(t provider.T) {
		t.Parallel()

		validPost, err := NewPostBuilder().
			WithPhoto(uuid.New(), 1).
			Build()
		require.NoError(t, err)

		photo := validPost.Photos().List()[0]

		srch := search.NewPostSearch()

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup post search", func(pctx provider.StepCtx) {
			srepo.On("SearchPosts", ctx, srch).Return([]*post.Post{validPost}, nil)
		})

		ptfrepo := new(MockFileRepository)
		t.WithNewStep("Setup photo not found", func(pctx provider.StepCtx) {
			ptfrepo.On("Get", ctx, photo.FileID()).Return(nil, assert.AnError)
		})

		pfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		t.WithNewStep("Attempt search with missing photo file", func(pctx provider.StepCtx) {
			_, err := svc.SearchPosts(ctx, srch)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})
}

func TestSearchServiceSearchPosts(t *testing.T) {
	suite.RunSuite(t, new(SearchServiceTestSearchTestSuite))
}
