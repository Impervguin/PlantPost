//go:build unit

package searchservice_test

import (
	"context"
	"testing"

	"PlantSite/internal/models"
	"PlantSite/internal/models/plant"
	searchservice "PlantSite/internal/services/search-service"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SearchServiceGetPlantTestSuite struct {
	suite.Suite
}

func (s *SearchServiceGetPlantTestSuite) BeforeEach(t provider.T) {
	t.Epic("Search Service")
	t.Feature("Search Management")
}

type GetPlantResultBuilder struct {
	plant      *plant.Plant
	photoFiles map[uuid.UUID]*models.File
}

func NewGetPlantResultBuilder(plant *plant.Plant) *GetPlantResultBuilder {
	return &GetPlantResultBuilder{
		plant:      plant,
		photoFiles: make(map[uuid.UUID]*models.File),
	}
}

func (b *GetPlantResultBuilder) WithPhotoFile(fileID uuid.UUID, fileName string) *GetPlantResultBuilder {
	b.photoFiles[fileID] = &models.File{ID: fileID, Name: fileName}
	return b
}

func (b *GetPlantResultBuilder) Build() *searchservice.GetPlant {
	photos := make([]searchservice.GetPlantPhoto, 0)

	b.plant.GetPhotos().Iterate(func(photo plant.PlantPhoto) error {
		if file, exists := b.photoFiles[photo.FileID()]; exists {
			photos = append(photos, searchservice.GetPlantPhoto{
				ID:          photo.ID(),
				File:        *file,
				Description: photo.Description(),
			})
		}
		return nil
	})

	mainPhotoFile := b.photoFiles[b.plant.MainPhotoID()]

	return &searchservice.GetPlant{
		ID:            b.plant.ID(),
		Name:          b.plant.GetName(),
		LatinName:     b.plant.GetLatinName(),
		Description:   b.plant.GetDescription(),
		MainPhoto:     *mainPhotoFile,
		Photos:        photos,
		Category:      b.plant.GetCategory(),
		Specification: b.plant.GetSpecification(),
		CreatedAt:     b.plant.CreatedAt(),
	}
}

func (s *SearchServiceGetPlantTestSuite) TestGetPlantByID(t provider.T) {
	t.Tags("get", "plant")
	t.Description("Test plant retrieval by ID functionality")
	t.Parallel()

	t.Run("Successful plant retrieval", func(t provider.T) {
		t.Parallel()

		ctx := context.Background()
		validPlantID := uuid.New()

		plantBuilder := NewPlantBuilder().
			WithPhoto(uuid.New(), "additional photo")

		validPlant, _ := plantBuilder.Build()

		mainPhotoFile := &models.File{ID: validPlant.MainPhotoID(), Name: "main.jpg"}
		var photoFile *models.File
		validPlant.GetPhotos().Iterate(func(photo plant.PlantPhoto) error {
			if photo.FileID() == mainPhotoFile.ID {
				return nil
			}
			photoFile = &models.File{ID: photo.FileID(), Name: "photo.jpg"}
			return nil
		})

		resultBuilder := NewGetPlantResultBuilder(validPlant).
			WithPhotoFile(mainPhotoFile.ID, mainPhotoFile.Name).
			WithPhotoFile(photoFile.ID, photoFile.Name)

		expectedResult := resultBuilder.Build()

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup plant repository", func(pctx provider.StepCtx) {
			srepo.On("GetPlantByID", ctx, validPlantID).Return(validPlant, nil)
		})

		pfrepo := new(MockFileRepository)
		t.WithNewStep("Setup file repository", func(pctx provider.StepCtx) {
			pfrepo.On("Get", ctx, mainPhotoFile.ID).Return(mainPhotoFile, nil)
			pfrepo.On("Get", ctx, photoFile.ID).Return(photoFile, nil)
		})

		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		var result *searchservice.GetPlant
		t.WithNewStep("Get plant by ID", func(pctx provider.StepCtx) {
			var err error
			result, err = svc.GetPlantByID(ctx, validPlantID)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify result", func(pctx provider.StepCtx) {
			assert.Equal(t, expectedResult, result)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			srepo.AssertExpectations(t)
			pfrepo.AssertExpectations(t)
		})
	})

	t.Run("Plant not found", func(t provider.T) {
		t.Parallel()

		ctx := context.Background()
		validPlantID := uuid.New()

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup plant not found", func(pctx provider.StepCtx) {
			srepo.On("GetPlantByID", ctx, validPlantID).Return(nil, assert.AnError)
		})

		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		t.WithNewStep("Attempt get non-existent plant", func(pctx provider.StepCtx) {
			_, err := svc.GetPlantByID(ctx, validPlantID)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Main photo not found", func(t provider.T) {
		t.Parallel()

		ctx := context.Background()
		validPlantID := uuid.New()

		validPlant, err := NewPlantBuilder().Build()
		require.NoError(t, err)

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup plant repository", func(pctx provider.StepCtx) {
			srepo.On("GetPlantByID", ctx, validPlantID).Return(validPlant, nil)
		})

		pfrepo := new(MockFileRepository)
		t.WithNewStep("Setup main photo not found", func(pctx provider.StepCtx) {
			pfrepo.On("Get", ctx, validPlant.MainPhotoID()).Return(nil, assert.AnError)
		})

		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		t.WithNewStep("Attempt get plant with missing main photo", func(pctx provider.StepCtx) {
			_, err := svc.GetPlantByID(ctx, validPlantID)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Plant with no additional photos", func(t provider.T) {
		t.Parallel()

		ctx := context.Background()
		validPlantID := uuid.New()

		validPlant, err := NewPlantBuilder().WithNoPhotos().Build()
		require.NoError(t, err)

		mainPhotoFile := &models.File{ID: validPlant.MainPhotoID(), Name: "main.jpg"}

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup plant repository", func(pctx provider.StepCtx) {
			srepo.On("GetPlantByID", ctx, validPlantID).Return(validPlant, nil)
		})

		pfrepo := new(MockFileRepository)
		t.WithNewStep("Setup file repository", func(pctx provider.StepCtx) {
			pfrepo.On("Get", ctx, mainPhotoFile.ID).Return(mainPhotoFile, nil)
		})

		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		var result *searchservice.GetPlant
		t.WithNewStep("Get plant with no additional photos", func(pctx provider.StepCtx) {
			var err error
			result, err = svc.GetPlantByID(ctx, validPlantID)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify empty photos", func(pctx provider.StepCtx) {
			assert.Empty(t, result.Photos)
		})
	})

	t.Run("Nil plant ID", func(t provider.T) {
		t.Parallel()

		ctx := context.Background()

		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		t.WithNewStep("Attempt get plant with nil ID", func(pctx provider.StepCtx) {
			_, err := svc.GetPlantByID(ctx, uuid.Nil)
			require.Error(t, err)
		})
	})
}

func TestSearchServiceGetPlant(t *testing.T) {
	suite.RunSuite(t, new(SearchServiceGetPlantTestSuite))
}
