//go:build unit

package searchservice_test

import (
	"context"
	"testing"

	"PlantSite/internal/models"
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/search"
	searchservice "PlantSite/internal/services/search-service"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SearchServiceTestSuite struct {
	suite.Suite
}

func (s *SearchServiceTestSuite) BeforeEach(t provider.T) {
	t.Epic("Search Service")
	t.Feature("Search Management")
}

func (s *SearchServiceTestSuite) TestSearchPlants(t provider.T) {
	t.Tags("search", "plants")
	t.Description("Test plant search functionality")
	t.Parallel()

	ctx := context.Background()

	t.Run("Successful search without filters", func(t provider.T) {
		t.Parallel()

		coniferousPlant, err := NewPlantBuilder().
			WithName("Pine").
			WithLatinName("Pinus sylvestris").
			WithCategory("coniferous").
			Build()
		require.NoError(t, err)

		deciduousPlant, err := NewPlantBuilder().
			WithName("Oak").
			WithLatinName("Quercus robur").
			WithCategory("deciduous").
			Build()
		require.NoError(t, err)

		coniferousMainPhoto := &models.File{ID: coniferousPlant.MainPhotoID(), Name: "pine.jpg"}
		deciduousMainPhoto := &models.File{ID: deciduousPlant.MainPhotoID(), Name: "oak.jpg"}

		searchQuery := search.NewPlantSearch()

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup plant search", func(pctx provider.StepCtx) {
			srepo.On("SearchPlants", ctx, searchQuery).Return([]*plant.Plant{coniferousPlant, deciduousPlant}, nil)
		})

		pfrepo := new(MockFileRepository)
		t.WithNewStep("Setup file repository", func(pctx provider.StepCtx) {
			pfrepo.On("Get", ctx, coniferousPlant.MainPhotoID()).Return(coniferousMainPhoto, nil)
			pfrepo.On("Get", ctx, deciduousPlant.MainPhotoID()).Return(deciduousMainPhoto, nil)
		})

		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		var results []*searchservice.SearchPlant
		t.WithNewStep("Search plants", func(pctx provider.StepCtx) {
			var err error
			results, err = svc.SearchPlants(ctx, searchQuery)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify results", func(pctx provider.StepCtx) {
			assert.Len(t, results, 2)
			assert.Equal(t, coniferousPlant.ID(), results[0].ID)
			assert.Equal(t, coniferousPlant.GetName(), results[0].Name)
			assert.Equal(t, *coniferousMainPhoto, results[0].MainPhoto)
		})

		t.WithNewStep("Verify expectations", func(pctx provider.StepCtx) {
			srepo.AssertExpectations(t)
			pfrepo.AssertExpectations(t)
		})
	})

	t.Run("Successful search with name filter", func(t provider.T) {
		t.Parallel()

		coniferousPlant, err := NewPlantBuilder().
			WithName("Pine").
			WithLatinName("Pinus sylvestris").
			WithCategory("coniferous").
			Build()
		require.NoError(t, err)

		mainPhoto := &models.File{ID: coniferousPlant.MainPhotoID(), Name: "pine.jpg"}

		searchQuery := search.NewPlantSearch()
		searchQuery.AddFilter(search.NewPlantNameFilter("Pine"))

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup plant search with name filter", func(pctx provider.StepCtx) {
			srepo.On("SearchPlants", ctx, searchQuery).Return([]*plant.Plant{coniferousPlant}, nil)
		})

		pfrepo := new(MockFileRepository)
		t.WithNewStep("Setup file repository", func(pctx provider.StepCtx) {
			pfrepo.On("Get", ctx, coniferousPlant.MainPhotoID()).Return(mainPhoto, nil)
		})

		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		var results []*searchservice.SearchPlant
		t.WithNewStep("Search plants with name filter", func(pctx provider.StepCtx) {
			var err error
			results, err = svc.SearchPlants(ctx, searchQuery)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify filtered results", func(pctx provider.StepCtx) {
			assert.Len(t, results, 1)
			assert.Equal(t, "Pine", results[0].Name)
		})
	})

	t.Run("Successful search with category filter", func(t provider.T) {
		t.Parallel()

		coniferousPlant, err := NewPlantBuilder().
			WithName("Pine").
			WithLatinName("Pinus sylvestris").
			WithCategory("coniferous").
			Build()
		require.NoError(t, err)

		mainPhoto := &models.File{ID: coniferousPlant.MainPhotoID(), Name: "pine.jpg"}

		searchQuery := search.NewPlantSearch()
		searchQuery.AddFilter(search.NewPlantCategoryFilter("coniferous"))

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup plant search with category filter", func(pctx provider.StepCtx) {
			srepo.On("SearchPlants", ctx, searchQuery).Return([]*plant.Plant{coniferousPlant}, nil)
		})

		pfrepo := new(MockFileRepository)
		t.WithNewStep("Setup file repository", func(pctx provider.StepCtx) {
			pfrepo.On("Get", ctx, coniferousPlant.MainPhotoID()).Return(mainPhoto, nil)
		})

		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		var results []*searchservice.SearchPlant
		t.WithNewStep("Search plants with category filter", func(pctx provider.StepCtx) {
			var err error
			results, err = svc.SearchPlants(ctx, searchQuery)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify filtered results", func(pctx provider.StepCtx) {
			assert.Len(t, results, 1)
			assert.Equal(t, "coniferous", results[0].Category)
		})
	})

	t.Run("Empty search results", func(t provider.T) {
		t.Parallel()

		searchQuery := search.NewPlantSearch()
		searchQuery.AddFilter(search.NewPlantNameFilter("Nonexistent Plant"))

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup empty search results", func(pctx provider.StepCtx) {
			srepo.On("SearchPlants", ctx, searchQuery).Return([]*plant.Plant{}, nil)
		})

		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		var results []*searchservice.SearchPlant
		t.WithNewStep("Search non-existent plants", func(pctx provider.StepCtx) {
			var err error
			results, err = svc.SearchPlants(ctx, searchQuery)
			require.NoError(t, err)
		})

		t.WithNewStep("Verify empty results", func(pctx provider.StepCtx) {
			assert.Empty(t, results)
		})
	})

	t.Run("Repository error during search", func(t provider.T) {
		t.Parallel()

		searchQuery := search.NewPlantSearch()

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup repository error", func(pctx provider.StepCtx) {
			srepo.On("SearchPlants", ctx, searchQuery).Return(nil, assert.AnError)
		})

		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		t.WithNewStep("Attempt search with repository error", func(pctx provider.StepCtx) {
			_, err := svc.SearchPlants(ctx, searchQuery)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})

	t.Run("Photo not found during search", func(t provider.T) {
		t.Parallel()

		coniferousPlant, err := NewPlantBuilder().
			WithName("Pine").
			WithLatinName("Pinus sylvestris").
			WithCategory("coniferous").
			Build()
		require.NoError(t, err)

		searchQuery := search.NewPlantSearch()

		srepo := new(MockSearchRepository)
		t.WithNewStep("Setup plant search", func(pctx provider.StepCtx) {
			srepo.On("SearchPlants", ctx, searchQuery).Return([]*plant.Plant{coniferousPlant}, nil)
		})

		pfrepo := new(MockFileRepository)
		t.WithNewStep("Setup photo not found", func(pctx provider.StepCtx) {
			pfrepo.On("Get", ctx, coniferousPlant.MainPhotoID()).Return(nil, assert.AnError)
		})

		ptfrepo := new(MockFileRepository)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		t.WithNewStep("Attempt search with missing photo", func(pctx provider.StepCtx) {
			_, err := svc.SearchPlants(ctx, searchQuery)
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
		})
	})
}

func TestSearchService(t *testing.T) {
	suite.RunSuite(t, new(SearchServiceTestSuite))
}
