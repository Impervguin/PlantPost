//go:build unit

package searchservice_test

import (
	"context"
	"testing"

	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	searchservice "PlantSite/internal/services/search-service"
	"PlantSite/internal/utils/logs"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SearchServiceTestGetPostSuite struct {
	suite.Suite
}

func (s *SearchServiceTestGetPostSuite) BeforeAll(t provider.T) {
	logs.InitNoopLogger()
}

func (s *SearchServiceTestGetPostSuite) BeforeEach(t provider.T) {
	t.Epic("Search Service")
	t.Feature("Search Management")
}

type GetPostResultBuilder struct {
	post       *post.Post
	photoFiles map[uuid.UUID]*models.File
}

func NewGetPostResultBuilder(post *post.Post) *GetPostResultBuilder {
	return &GetPostResultBuilder{
		post:       post,
		photoFiles: make(map[uuid.UUID]*models.File),
	}
}

func (b *GetPostResultBuilder) WithPhotoFile(fileID uuid.UUID, fileName string) *GetPostResultBuilder {
	b.photoFiles[fileID] = &models.File{ID: fileID, Name: fileName}
	return b
}

func (b *GetPostResultBuilder) Build() *searchservice.GetPost {
	photos := make([]searchservice.GetPostPhoto, 0)

	for _, photo := range b.post.Photos().List() {
		if file, exists := b.photoFiles[photo.FileID()]; exists {
			photos = append(photos, searchservice.GetPostPhoto{
				ID:          photo.ID(),
				PlaceNumber: photo.PlaceNumber(),
				File:        *file,
			})
		}
	}

	return &searchservice.GetPost{
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

func (s *SearchServiceTestGetPostSuite) TestGetPost(t provider.T) {
	t.Tags("get", "post")
	t.Description("Test post retrieval by ID functionality")
	t.Parallel()

	ctx := context.Background()

	t.Run("Successful post retrieval", func(t provider.T) {
		t.Parallel()

		validPost, err := NewPostBuilder().
			WithPhoto(uuid.New(), 1).
			WithPhoto(uuid.New(), 2).
			Build()
		require.NoError(t, err)
		validPostID := validPost.ID()

		photo1 := validPost.Photos().List()[0]
		photo2 := validPost.Photos().List()[1]
		photoFile1 := &models.File{ID: photo1.FileID(), Name: "photo1.jpg"}
		photoFile2 := &models.File{ID: photo2.FileID(), Name: "photo2.jpg"}

		resultBuilder := NewGetPostResultBuilder(validPost).
			WithPhotoFile(photo1.FileID(), "photo1.jpg").
			WithPhotoFile(photo2.FileID(), "photo2.jpg")

		expectedResult := resultBuilder.Build()

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup post repository", func(pctx provider.StepCtx) {
			srepo.On("GetPostByID", ctx, validPostID).Return(validPost, nil)
		})

		ptfrepo := new(MockFileRepository)
		t.WithNewStep("Setup file repository", func(pctx provider.StepCtx) {
			ptfrepo.On("Get", ctx, photo1.FileID()).Return(photoFile1, nil)
			ptfrepo.On("Get", ctx, photo2.FileID()).Return(photoFile2, nil)
		})

		pfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		var result *searchservice.GetPost
		t.WithNewStep("Get post by ID", func(pctx provider.StepCtx) {
			var err error
			result, err = svc.GetPost(ctx, validPostID)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify result", func(pctx provider.StepCtx) {
			assert.Equal(t, expectedResult, result)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			srepo.AssertExpectations(t)
			ptfrepo.AssertExpectations(t)
		})
	})

	t.Run("Post not found", func(t provider.T) {
		t.Parallel()

		validPostID := uuid.New()

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup post not found", func(pctx provider.StepCtx) {
			srepo.On("GetPostByID", ctx, validPostID).Return(nil, assert.AnError)
		})

		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		t.WithNewStep("Attempt get non-existent post", func(pctx provider.StepCtx) {
			_, err := svc.GetPost(ctx, validPostID)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Photo file not found", func(t provider.T) {
		t.Parallel()

		validPost, err := NewPostBuilder().
			WithPhoto(uuid.New(), 1).
			Build()
		require.NoError(t, err)
		validPostID := validPost.ID()

		photo := validPost.Photos().List()[0]

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup post repository", func(pctx provider.StepCtx) {
			srepo.On("GetPostByID", ctx, validPostID).Return(validPost, nil)
		})

		ptfrepo := new(MockFileRepository)
		t.WithNewStep("Setup file not found", func(pctx provider.StepCtx) {
			ptfrepo.On("Get", ctx, photo.FileID()).Return(nil, assert.AnError)
		})

		pfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		t.WithNewStep("Attempt get post with missing photo file", func(pctx provider.StepCtx) {
			_, err := svc.GetPost(ctx, validPostID)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Post with no photos", func(t provider.T) {
		t.Parallel()

		validPost, err := NewPostBuilder().
			WithNoPhotos().
			Build()
		require.NoError(t, err)
		validPostID := validPost.ID()

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup post with no photos", func(pctx provider.StepCtx) {
			srepo.On("GetPostByID", ctx, validPostID).Return(validPost, nil)
		})

		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		var result *searchservice.GetPost
		t.WithNewStep("Get post with no photos", func(pctx provider.StepCtx) {
			var err error
			result, err = svc.GetPost(ctx, validPostID)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify empty photos", func(pctx provider.StepCtx) {
			assert.Empty(t, result.Photos)
		})
	})

	t.Run("Nil post ID", func(t provider.T) {
		t.Parallel()

		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		t.WithNewStep("Attempt get post with nil ID", func(pctx provider.StepCtx) {
			_, err := svc.GetPost(ctx, uuid.Nil)
			require.Error(t, err)
		})
	})
}

func TestSearchServiceGetPost(t *testing.T) {
	suite.RunSuite(t, new(SearchServiceTestGetPostSuite))
}
