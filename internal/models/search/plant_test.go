package search

import (
	"PlantSite/internal/models/plant"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPlant создает тестовое растение с заданными параметрами
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

func TestPlantFilters(t *testing.T) {
	// Создаем тестовые спецификации
	coniferousSpec, err := plant.NewConiferousSpecification(10.5, 2.3, 5, plant.MediumMoisture, plant.HalfShadow, plant.MediumSoil, 6)
	require.NoError(t, err)

	deciduousSpec, err := plant.NewDeciduousSpecification(8.2, 1.8, plant.Spring, 6, plant.MediumMoisture, plant.HalfShadow, plant.MediumSoil, 5)
	require.NoError(t, err)
	// Создаем тестовые растения
	coniferousPlant, err := mockPlant("Pine", "Pinus sylvestris", "coniferous", coniferousSpec)
	require.NoError(t, err)
	deciduousPlant, err := mockPlant("Oak", "Quercus robur", "deciduous", deciduousSpec)
	require.NoError(t, err)

	t.Run("PlantNameFilter", func(t *testing.T) {
		filter := NewPlantNameFilter("Pine")
		assert.True(t, filter.Filter(coniferousPlant))
		assert.False(t, filter.Filter(deciduousPlant))

		filter = NewPlantNameFilter("Oak")
		assert.False(t, filter.Filter(coniferousPlant))
		assert.True(t, filter.Filter(deciduousPlant))
	})

	t.Run("PlantCategoryFilter", func(t *testing.T) {
		filter := NewPlantCategoryFilter("coniferous")
		assert.True(t, filter.Filter(coniferousPlant))
		assert.False(t, filter.Filter(deciduousPlant))

		filter = NewPlantCategoryFilter("deciduous")
		assert.False(t, filter.Filter(coniferousPlant))
		assert.True(t, filter.Filter(deciduousPlant))
	})

	t.Run("PlantLatinNameFilter", func(t *testing.T) {
		filter := NewPlantLatinNameFilter("Pinus sylvestris")
		assert.True(t, filter.Filter(coniferousPlant))
		assert.False(t, filter.Filter(deciduousPlant))

		filter = NewPlantLatinNameFilter("Quercus robur")
		assert.False(t, filter.Filter(coniferousPlant))
		assert.True(t, filter.Filter(deciduousPlant))
	})

	t.Run("PlantHeightFilter", func(t *testing.T) {
		tests := []struct {
			name  string
			min   float64
			max   float64
			conif bool
			decid bool
		}{
			{"Both plants match", 8.0, 12.0, true, true},
			{"Only coniferous", 10.0, 11.0, true, false},
			{"Only deciduous", 8.0, 9.0, false, true},
			{"None match", 15.0, 20.0, false, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				filter := NewPlantHeightFilter(tt.min, tt.max)
				assert.Equal(t, tt.conif, filter.Filter(coniferousPlant))
				assert.Equal(t, tt.decid, filter.Filter(deciduousPlant))
			})
		}
	})

	t.Run("PlantDiameterFilter", func(t *testing.T) {
		tests := []struct {
			name  string
			min   float64
			max   float64
			conif bool
			decid bool
		}{
			{"Both plants match", 1.0, 3.0, true, true},
			{"Only coniferous", 2.0, 3.0, true, false},
			{"Only deciduous", 1.0, 2.0, false, true},
			{"None match", 3.0, 5.0, false, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				filter := NewPlantDiameterFilter(tt.min, tt.max)
				assert.Equal(t, tt.conif, filter.Filter(coniferousPlant))
				assert.Equal(t, tt.decid, filter.Filter(deciduousPlant))
			})
		}
	})

	t.Run("PlantSoilAcidityFilter", func(t *testing.T) {
		filter := NewSoilAcidityFilter(4, 7)
		assert.True(t, filter.Filter(coniferousPlant))
		assert.True(t, filter.Filter(deciduousPlant))

		filter = NewSoilAcidityFilter(6, 7)
		assert.False(t, filter.Filter(coniferousPlant))
		assert.True(t, filter.Filter(deciduousPlant))
	})

	t.Run("PlantSoilMoistureFilter", func(t *testing.T) {
		filter := NewSoilMoistureFilter([]plant.SoilMoisture{plant.MediumMoisture})
		assert.True(t, filter.Filter(coniferousPlant))
		assert.True(t, filter.Filter(deciduousPlant))

		filter = NewSoilMoistureFilter([]plant.SoilMoisture{plant.DryMoisture})
		assert.False(t, filter.Filter(coniferousPlant))
		assert.False(t, filter.Filter(deciduousPlant))
	})

	t.Run("PlantLightRelationFilter", func(t *testing.T) {
		filter := NewLightRelationFilter([]plant.LightRelation{plant.HalfShadow})
		assert.True(t, filter.Filter(coniferousPlant))
		assert.True(t, filter.Filter(deciduousPlant))

		filter = NewLightRelationFilter([]plant.LightRelation{plant.Shadow})
		assert.False(t, filter.Filter(coniferousPlant))
		assert.False(t, filter.Filter(deciduousPlant))
	})

	t.Run("PlantHardinessFilter", func(t *testing.T) {
		filter := NewWinterHardinessFilter(4, 7)
		assert.True(t, filter.Filter(coniferousPlant))
		assert.True(t, filter.Filter(deciduousPlant))

		filter = NewWinterHardinessFilter(7, 9)
		assert.False(t, filter.Filter(coniferousPlant))
		assert.False(t, filter.Filter(deciduousPlant))
	})

	t.Run("PlantSoilTypeFilter", func(t *testing.T) {
		filter := NewSoilTypeFilter([]plant.Soil{plant.MediumSoil})
		assert.True(t, filter.Filter(coniferousPlant))
		assert.True(t, filter.Filter(deciduousPlant))

		filter = NewSoilTypeFilter([]plant.Soil{plant.LightSoil})
		assert.False(t, filter.Filter(coniferousPlant))
		assert.False(t, filter.Filter(deciduousPlant))
	})

	t.Run("PlantFloweringPeriodFilter", func(t *testing.T) {
		filter := NewFloweringPeriodFilter([]plant.FloweringPeriod{plant.Spring})
		assert.False(t, filter.Filter(coniferousPlant)) // Не должно работать для хвойных
		assert.True(t, filter.Filter(deciduousPlant))

		filter = NewFloweringPeriodFilter([]plant.FloweringPeriod{plant.Summer})
		assert.False(t, filter.Filter(deciduousPlant))
	})
}
