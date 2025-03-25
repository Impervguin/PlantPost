package searchservice_test

import (
	"context"
	"testing"

	"PlantSite/internal/models"
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/search"
	searchservice "PlantSite/internal/search-service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockPlant(name, latinName, category string, spec plant.PlantSpecification) (*plant.Plant, error) {

	return plant.NewPlant(
		name,
		latinName,
		"Test description",
		uuid.New(),
		*plant.NewPlantPhotos(),
		category,
		spec,
	)
}

func TestSearchPlants(t *testing.T) {
	ctx := context.Background()
	validFileID := uuid.New()

	coniferousSpec, err := plant.NewConiferousSpecification(10.5, 2.3, 5, plant.MediumMoisture, plant.HalfShadow, plant.MediumSoil, 6)
	require.NoError(t, err)

	deciduousSpec, err := plant.NewDeciduousSpecification(8.2, 1.8, plant.Spring, 6, plant.MediumMoisture, plant.HalfShadow, plant.MediumSoil, 5)
	require.NoError(t, err)
	// Создаем тестовые растения
	coniferousPlant, err := mockPlant("Pine", "Pinus sylvestris", "coniferous", coniferousSpec)
	require.NoError(t, err)
	deciduousPlant, err := mockPlant("Oak", "Quercus robur", "deciduous", deciduousSpec)
	require.NoError(t, err)

	mainPhotoFile := &models.File{ID: validFileID, Name: "pine.jpg"}

	t.Run("SuccessWithoutFilters", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		searchQuery := search.NewPlantSearch()

		srepo.On("SearchPlants", ctx, searchQuery).Return([]*plant.Plant{coniferousPlant, deciduousPlant}, nil)
		pfrepo.On("Get", ctx, coniferousPlant.MainPhotoID()).Return(mainPhotoFile, nil)
		pfrepo.On("Get", ctx, deciduousPlant.MainPhotoID()).Return(&models.File{ID: deciduousPlant.MainPhotoID()}, nil)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		results, err := svc.SearchPlants(ctx, searchQuery)
		require.NoError(t, err)
		assert.Len(t, results, 2)

		assert.Equal(t, coniferousPlant.ID(), results[0].ID)
		assert.Equal(t, coniferousPlant.GetName(), results[0].Name)
		assert.Equal(t, *mainPhotoFile, results[0].MainPhoto)

		srepo.AssertExpectations(t)
		pfrepo.AssertExpectations(t)
	})

	t.Run("SuccessWithNameFilter", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		searchQuery := search.NewPlantSearch()
		searchQuery.AddFilter(search.NewPlantNameFilter("Pine"))

		srepo.On("SearchPlants", ctx, searchQuery).Return([]*plant.Plant{coniferousPlant}, nil)
		pfrepo.On("Get", ctx, coniferousPlant.MainPhotoID()).Return(mainPhotoFile, nil)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		results, err := svc.SearchPlants(ctx, searchQuery)
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Pine", results[0].Name)
	})

	t.Run("SuccessWithCategoryFilter", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		searchQuery := search.NewPlantSearch()
		searchQuery.AddFilter(search.NewPlantCategoryFilter("coniferous"))

		srepo.On("SearchPlants", ctx, searchQuery).Return([]*plant.Plant{coniferousPlant}, nil)
		pfrepo.On("Get", ctx, coniferousPlant.MainPhotoID()).Return(mainPhotoFile, nil)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		results, err := svc.SearchPlants(ctx, searchQuery)
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "coniferous", results[0].Category)
	})

	t.Run("SuccessWithHeightFilter", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		searchQuery := search.NewPlantSearch()
		searchQuery.AddFilter(search.NewPlantHeightFilter(10.0, 11.0))

		srepo.On("SearchPlants", ctx, searchQuery).Return([]*plant.Plant{coniferousPlant}, nil)
		pfrepo.On("Get", ctx, coniferousPlant.MainPhotoID()).Return(mainPhotoFile, nil)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		results, err := svc.SearchPlants(ctx, searchQuery)
		require.NoError(t, err)
		assert.Len(t, results, 1)
	})

	t.Run("SuccessWithMultipleFilters", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		searchQuery := search.NewPlantSearch()
		searchQuery.AddFilter(search.NewPlantCategoryFilter("coniferous"))
		searchQuery.AddFilter(search.NewPlantHeightFilter(10.0, 11.0))
		searchQuery.AddFilter(search.NewSoilAcidityFilter(4, 5))

		srepo.On("SearchPlants", ctx, searchQuery).Return([]*plant.Plant{coniferousPlant}, nil)
		pfrepo.On("Get", ctx, coniferousPlant.MainPhotoID()).Return(mainPhotoFile, nil)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		results, err := svc.SearchPlants(ctx, searchQuery)
		require.NoError(t, err)
		assert.Len(t, results, 1)
	})

	t.Run("EmptyResults", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		searchQuery := search.NewPlantSearch()
		searchQuery.AddFilter(search.NewPlantNameFilter("Nonexistent Plant"))

		srepo.On("SearchPlants", ctx, searchQuery).Return([]*plant.Plant{}, nil)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		results, err := svc.SearchPlants(ctx, searchQuery)
		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		searchQuery := search.NewPlantSearch()
		srepo.On("SearchPlants", ctx, searchQuery).Return(nil, assert.AnError)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		_, err := svc.SearchPlants(ctx, searchQuery)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("PhotoNotFound", func(t *testing.T) {
		srepo := new(MockSearchRepository)
		pfrepo := new(MockFileRepository)
		ptfrepo := new(MockFileRepository)

		searchQuery := search.NewPlantSearch()
		srepo.On("SearchPlants", ctx, searchQuery).Return([]*plant.Plant{coniferousPlant}, nil)
		pfrepo.On("Get", ctx, coniferousPlant.MainPhotoID()).Return(nil, assert.AnError)

		svc := searchservice.NewSearchService(srepo, pfrepo, ptfrepo)

		_, err := svc.SearchPlants(ctx, searchQuery)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
