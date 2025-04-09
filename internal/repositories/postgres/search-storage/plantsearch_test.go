//go:build integration

package searchstorage_test

import (
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/search"
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *SearchRepositoryTestSuite) TestSearchPlantsByType() {
	ctx := context.Background()

	// Create test plants - one coniferous and one deciduous
	coniferousPlant := s.createConiferousPlant(ctx, "Pine Tree", 1.5, 0.5, plant.MediumMoisture, 10, plant.Light, plant.MediumSoil, plant.WinterHardiness(10))

	deciduousPlant := s.createDeciduousPlant(ctx, "Oak Tree", 5.0, 1.0, plant.HighMoisture, 10, plant.Light, plant.MediumSoil, plant.WinterHardiness(10), plant.Spring)

	_, err := s.plantRepo.Create(ctx, coniferousPlant)
	require.NoError(s.T(), err)
	_, err = s.plantRepo.Create(ctx, deciduousPlant)
	require.NoError(s.T(), err)

	// Test searching for coniferous plants
	coniferousSearch := search.NewPlantSearch()
	coniferousSearch.AddFilter(search.NewPlantCategoryFilter(plant.ConiferousCategory))
	coniferousPlants, err := s.searchRepo.SearchPlants(ctx, coniferousSearch)
	require.NoError(s.T(), err)
	assert.Len(s.T(), coniferousPlants, 1)
	assert.Equal(s.T(), "Pine Tree", coniferousPlants[0].GetName())

	// Test searching for deciduous plants
	deciduousSearch := search.NewPlantSearch()
	deciduousSearch.AddFilter(search.NewPlantCategoryFilter(plant.DeciduousCategory))
	deciduousPlants, err := s.searchRepo.SearchPlants(ctx, deciduousSearch)
	require.NoError(s.T(), err)
	assert.Len(s.T(), deciduousPlants, 1)
	assert.Equal(s.T(), "Oak Tree", deciduousPlants[0].GetName())
}

func (s *SearchRepositoryTestSuite) TestSearchPlantsByHeightWithDifferentTypes() {
	ctx := context.Background()

	// Create plants of different types with varying heights
	tallConifer := s.createConiferousPlant(ctx, "Tall Pine", 10.0, 2.0, plant.MediumMoisture, 10, plant.Light, plant.MediumSoil, plant.WinterHardiness(10))
	shortConifer := s.createConiferousPlant(ctx, "Short Pine", 1.5, 0.5, plant.MediumMoisture, 10, plant.Light, plant.MediumSoil, plant.WinterHardiness(10))
	tallDeciduous := s.createDeciduousPlant(ctx, "Tall Oak", 12.0, 3.0, plant.HighMoisture, 10, plant.Light, plant.MediumSoil, plant.WinterHardiness(10), plant.Spring)

	_, err := s.plantRepo.Create(ctx, tallConifer)
	require.NoError(s.T(), err)
	_, err = s.plantRepo.Create(ctx, shortConifer)
	require.NoError(s.T(), err)
	_, err = s.plantRepo.Create(ctx, tallDeciduous)
	require.NoError(s.T(), err)

	// Search for tall plants (both types)
	tallSearch := search.NewPlantSearch()
	tallSearch.AddFilter(search.NewPlantHeightFilter(8.0, 15.0))
	tallPlants, err := s.searchRepo.SearchPlants(ctx, tallSearch)
	require.NoError(s.T(), err)
	assert.Len(s.T(), tallPlants, 2)

	// Verify both types are included
	var foundConifer, foundDeciduous bool
	for _, p := range tallPlants {
		if p.GetName() == "Tall Pine" {
			foundConifer = true
			assert.IsType(s.T(), &plant.ConiferousSpecification{}, p.GetSpecification())
		} else if p.GetName() == "Tall Oak" {
			foundDeciduous = true
			assert.IsType(s.T(), &plant.DeciduousSpecification{}, p.GetSpecification())
		}
	}
	assert.True(s.T(), foundConifer)
	assert.True(s.T(), foundDeciduous)
}

func (s *SearchRepositoryTestSuite) TestSearchPlantsByDiameterWithDifferentTypes() {
	ctx := context.Background()
	// Create plants of different types with varying diameters
	tallConifer := s.createConiferousPlant(ctx, "Wide Pine", 1.5, 2.0, plant.MediumMoisture, 10, plant.Light, plant.MediumSoil, plant.WinterHardiness(10))
	shortConifer := s.createConiferousPlant(ctx, "Not So Wide Pine", 1.5, 0.5, plant.MediumMoisture, 10, plant.Light, plant.MediumSoil, plant.WinterHardiness(10))

	tallDeciduous := s.createDeciduousPlant(ctx, "Wide Oak", 2.0, 3., plant.HighMoisture, 10, plant.Light, plant.MediumSoil, plant.WinterHardiness(10), plant.Spring)

	_, err := s.plantRepo.Create(ctx, tallConifer)
	require.NoError(s.T(), err)
	_, err = s.plantRepo.Create(ctx, shortConifer)
	require.NoError(s.T(), err)
	_, err = s.plantRepo.Create(ctx, tallDeciduous)
	require.NoError(s.T(), err)
	// Search for tall plants (both types)
	tallSearch := search.NewPlantSearch()
	tallSearch.AddFilter(search.NewPlantDiameterFilter(1.0, 3.1))
	tallPlants, err := s.searchRepo.SearchPlants(ctx, tallSearch)
	require.NoError(s.T(), err)
	assert.Len(s.T(), tallPlants, 2)
	// Verify both types are included
	var foundConifer, foundDeciduous bool
	for _, p := range tallPlants {
		if p.GetName() == "Wide Pine" {
			foundConifer = true
			assert.IsType(s.T(), &plant.ConiferousSpecification{}, p.GetSpecification())
		} else if p.GetName() == "Wide Oak" {
			foundDeciduous = true
			assert.IsType(s.T(), &plant.DeciduousSpecification{}, p.GetSpecification())
		}
	}
	assert.True(s.T(), foundConifer)
	assert.True(s.T(), foundDeciduous)
}

func (s *SearchRepositoryTestSuite) TestSearchPlantsByLightRelation() {
	ctx := context.Background()
	// Create plants with different light relations
	lightConifer := s.createConiferousPlant(ctx, "Light Pine", 1.5, 0.5, plant.MediumMoisture, 10, plant.Light, plant.MediumSoil, plant.WinterHardiness(10))
	shadowConifer := s.createConiferousPlant(ctx, "Dark Pine", 1.5, 0.5, plant.MediumMoisture, 10, plant.Shadow, plant.MediumSoil, plant.WinterHardiness(10))
	HalfShadowDeciduous := s.createDeciduousPlant(ctx, "Half Shadow Oak", 2.0, 1.0, plant.HighMoisture, 10, plant.HalfShadow, plant.MediumSoil, plant.WinterHardiness(10), plant.Spring)

	_, err := s.plantRepo.Create(ctx, lightConifer)
	require.NoError(s.T(), err)
	_, err = s.plantRepo.Create(ctx, shadowConifer)
	require.NoError(s.T(), err)
	_, err = s.plantRepo.Create(ctx, HalfShadowDeciduous)
	require.NoError(s.T(), err)
	// Search for plants with light relation
	lightSearch := search.NewPlantSearch()
	lightSearch.AddFilter(search.NewLightRelationFilter([]plant.LightRelation{plant.Light, plant.HalfShadow}))
	lightPlants, err := s.searchRepo.SearchPlants(ctx, lightSearch)
	require.NoError(s.T(), err)
	assert.Len(s.T(), lightPlants, 2)
	// Verify both types are included
	var foundConifer, foundDeciduous bool
	for _, p := range lightPlants {
		if p.GetName() == "Light Pine" {
			foundConifer = true
			assert.IsType(s.T(), &plant.ConiferousSpecification{}, p.GetSpecification())
		} else if p.GetName() == "Half Shadow Oak" {
			foundDeciduous = true
			assert.IsType(s.T(), &plant.DeciduousSpecification{}, p.GetSpecification())
		}
	}
	assert.True(s.T(), foundConifer)
	assert.True(s.T(), foundDeciduous)
}

func (s *SearchRepositoryTestSuite) TestSearchPlantsByFloweringPeriod() {
	ctx := context.Background()
	// Create deciduous plants with different flowering periods
	springFlowering := s.createDeciduousPlant(ctx, "Spring Bloomer", 3.0, 1.0, plant.MediumMoisture, 10, plant.HalfShadow, plant.MediumSoil, plant.WinterHardiness(10), plant.Spring)
	summerFlowering := s.createDeciduousPlant(ctx, "Summer Bloomer", 4.0, 1.2, plant.MediumMoisture, 10, plant.HalfShadow, plant.MediumSoil, plant.WinterHardiness(10), plant.Summer)
	fallFlowering := s.createDeciduousPlant(ctx, "Fall Bloomer", 5.0, 1.5, plant.MediumMoisture, 10, plant.HalfShadow, plant.MediumSoil, plant.WinterHardiness(10), plant.Autumn)
	conif := s.createConiferousPlant(ctx, "Spring Conifer", 1.0, 0.5, plant.MediumMoisture, 10, plant.Light, plant.MediumSoil, plant.WinterHardiness(10))

	_, err := s.plantRepo.Create(ctx, springFlowering)
	require.NoError(s.T(), err)
	_, err = s.plantRepo.Create(ctx, summerFlowering)
	require.NoError(s.T(), err)
	_, err = s.plantRepo.Create(ctx, fallFlowering)
	require.NoError(s.T(), err)
	_, err = s.plantRepo.Create(ctx, conif)
	require.NoError(s.T(), err)
	// Search for spring flowering plants
	springSearch := search.NewPlantSearch()
	springSearch.AddFilter(search.NewFloweringPeriodFilter([]plant.FloweringPeriod{plant.Spring}))
	springPlants, err := s.searchRepo.SearchPlants(ctx, springSearch)
	require.NoError(s.T(), err)
	assert.Len(s.T(), springPlants, 1)
	assert.Equal(s.T(), "Spring Bloomer", springPlants[0].GetName())
}
