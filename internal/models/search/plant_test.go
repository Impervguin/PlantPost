//go:build unit

package search

import (
	"PlantSite/internal/models/plant"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
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

type PlantFiltersTestSuite struct {
	suite.Suite
	coniferousPlant *plant.Plant
	deciduousPlant  *plant.Plant
	coniferousSpec  *plant.ConiferousSpecification
	deciduousSpec   *plant.DeciduousSpecification
}

func (s *PlantFiltersTestSuite) BeforeEach(t provider.T) {
	t.Epic("Search")
	t.Feature("Plant Filters")

	var err error
	s.coniferousSpec, err = plant.NewConiferousSpecification(
		10.5,
		2.3,
		5,
		plant.MediumMoisture,
		plant.HalfShadow,
		plant.MediumSoil,
		6,
	)
	require.NoError(t, err)

	s.deciduousSpec, err = plant.NewDeciduousSpecification(
		8.2,
		1.8,
		plant.Spring,
		6,
		plant.MediumMoisture,
		plant.HalfShadow,
		plant.MediumSoil,
		5,
	)
	require.NoError(t, err)

	s.coniferousPlant, err = mockPlant("Pine", "Pinus sylvestris", "coniferous", s.coniferousSpec)
	require.NoError(t, err)

	s.deciduousPlant, err = mockPlant("Oak", "Quercus robur", "deciduous", s.deciduousSpec)
	require.NoError(t, err)
}

func (s *PlantFiltersTestSuite) TestPlantNameFilter(t provider.T) {
	t.Tags("filter", "name")
	t.Description("Test plant name filtering")

	filter := NewPlantNameFilter("Pine")
	t.WithNewStep("Filter by Pine name", func(ctx provider.StepCtx) {
		assert.True(t, filter.Filter(s.coniferousPlant))
		assert.False(t, filter.Filter(s.deciduousPlant))
	})

	filter = NewPlantNameFilter("Oak")
	t.WithNewStep("Filter by Oak name", func(ctx provider.StepCtx) {
		assert.False(t, filter.Filter(s.coniferousPlant))
		assert.True(t, filter.Filter(s.deciduousPlant))
	})
}

func (s *PlantFiltersTestSuite) TestPlantCategoryFilter(t provider.T) {
	t.Tags("filter", "category")
	t.Description("Test plant category filtering")

	filter := NewPlantCategoryFilter("coniferous")
	t.WithNewStep("Filter by coniferous category", func(ctx provider.StepCtx) {
		assert.True(t, filter.Filter(s.coniferousPlant))
		assert.False(t, filter.Filter(s.deciduousPlant))
	})

	filter = NewPlantCategoryFilter("deciduous")
	t.WithNewStep("Filter by deciduous category", func(ctx provider.StepCtx) {
		assert.False(t, filter.Filter(s.coniferousPlant))
		assert.True(t, filter.Filter(s.deciduousPlant))
	})
}

func (s *PlantFiltersTestSuite) TestPlantLatinNameFilter(t provider.T) {
	t.Tags("filter", "latin")
	t.Description("Test plant latin name filtering")

	filter := NewPlantLatinNameFilter("Pinus sylvestris")
	t.WithNewStep("Filter by Pinus sylvestris", func(ctx provider.StepCtx) {
		assert.True(t, filter.Filter(s.coniferousPlant))
		assert.False(t, filter.Filter(s.deciduousPlant))
	})

	filter = NewPlantLatinNameFilter("Quercus robur")
	t.WithNewStep("Filter by Quercus robur", func(ctx provider.StepCtx) {
		assert.False(t, filter.Filter(s.coniferousPlant))
		assert.True(t, filter.Filter(s.deciduousPlant))
	})
}

func (s *PlantFiltersTestSuite) TestPlantHeightFilter(t provider.T) {
	t.Tags("filter", "height")
	t.Description("Test plant height range filtering")

	tests := []struct {
		name  string
		min   float64
		max   float64
		conif bool
		decid bool
	}{
		{"Both plants match range", 8.0, 12.0, true, true},
		{"Only coniferous matches", 10.0, 11.0, true, false},
		{"Only deciduous matches", 8.0, 9.0, false, true},
		{"No plants match range", 15.0, 20.0, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t provider.T) {
			t.WithNewStep("Prepare height filter parameters", func(ctx provider.StepCtx) {
				ctx.WithParameters(
					allure.NewParameter("MinHeight", tt.min),
					allure.NewParameter("MaxHeight", tt.max),
				)
			})

			filter := NewPlantHeightFilter(tt.min, tt.max)
			t.WithNewStep("Apply height filter", func(ctx provider.StepCtx) {
				assert.Equal(t, tt.conif, filter.Filter(s.coniferousPlant))
				assert.Equal(t, tt.decid, filter.Filter(s.deciduousPlant))
			})
		})
	}
}

func (s *PlantFiltersTestSuite) TestPlantDiameterFilter(t provider.T) {
	t.Tags("filter", "diameter")
	t.Description("Test plant diameter range filtering")

	tests := []struct {
		name  string
		min   float64
		max   float64
		conif bool
		decid bool
	}{
		{"Both plants match range", 1.0, 3.0, true, true},
		{"Only coniferous matches", 2.0, 3.0, true, false},
		{"Only deciduous matches", 1.0, 2.0, false, true},
		{"No plants match range", 3.0, 5.0, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t provider.T) {
			t.WithNewStep("Prepare diameter filter parameters", func(ctx provider.StepCtx) {
				ctx.WithParameters(
					allure.NewParameter("MinDiameter", tt.min),
					allure.NewParameter("MaxDiameter", tt.max),
				)
			})

			filter := NewPlantDiameterFilter(tt.min, tt.max)
			t.WithNewStep("Apply diameter filter", func(ctx provider.StepCtx) {
				assert.Equal(t, tt.conif, filter.Filter(s.coniferousPlant))
				assert.Equal(t, tt.decid, filter.Filter(s.deciduousPlant))
			})
		})
	}
}

func (s *PlantFiltersTestSuite) TestPlantSoilAcidityFilter(t provider.T) {
	t.Tags("filter", "soil", "acidity")
	t.Description("Test plant soil acidity range filtering")

	filter := NewSoilAcidityFilter(4, 7)
	t.WithNewStep("Filter by acidity range 4-7", func(ctx provider.StepCtx) {
		assert.True(t, filter.Filter(s.coniferousPlant))
		assert.True(t, filter.Filter(s.deciduousPlant))
	})

	filter = NewSoilAcidityFilter(6, 7)
	t.WithNewStep("Filter by acidity range 6-7", func(ctx provider.StepCtx) {
		assert.False(t, filter.Filter(s.coniferousPlant))
		assert.True(t, filter.Filter(s.deciduousPlant))
	})
}

func (s *PlantFiltersTestSuite) TestPlantSoilMoistureFilter(t provider.T) {
	t.Tags("filter", "soil", "moisture")
	t.Description("Test plant soil moisture filtering")

	filter := NewSoilMoistureFilter([]plant.SoilMoisture{plant.MediumMoisture})
	t.WithNewStep("Filter by medium moisture", func(ctx provider.StepCtx) {
		assert.True(t, filter.Filter(s.coniferousPlant))
		assert.True(t, filter.Filter(s.deciduousPlant))
	})

	filter = NewSoilMoistureFilter([]plant.SoilMoisture{plant.DryMoisture})
	t.WithNewStep("Filter by dry moisture", func(ctx provider.StepCtx) {
		assert.False(t, filter.Filter(s.coniferousPlant))
		assert.False(t, filter.Filter(s.deciduousPlant))
	})
}

func (s *PlantFiltersTestSuite) TestPlantLightRelationFilter(t provider.T) {
	t.Tags("filter", "light")
	t.Description("Test plant light relation filtering")

	filter := NewLightRelationFilter([]plant.LightRelation{plant.HalfShadow})
	t.WithNewStep("Filter by half shadow", func(ctx provider.StepCtx) {
		assert.True(t, filter.Filter(s.coniferousPlant))
		assert.True(t, filter.Filter(s.deciduousPlant))
	})

	filter = NewLightRelationFilter([]plant.LightRelation{plant.Shadow})
	t.WithNewStep("Filter by shadow", func(ctx provider.StepCtx) {
		assert.False(t, filter.Filter(s.coniferousPlant))
		assert.False(t, filter.Filter(s.deciduousPlant))
	})
}

func (s *PlantFiltersTestSuite) TestPlantHardinessFilter(t provider.T) {
	t.Tags("filter", "hardiness")
	t.Description("Test plant winter hardiness filtering")

	filter := NewWinterHardinessFilter(4, 7)
	t.WithNewStep("Filter by hardiness range 4-7", func(ctx provider.StepCtx) {
		assert.True(t, filter.Filter(s.coniferousPlant))
		assert.True(t, filter.Filter(s.deciduousPlant))
	})

	filter = NewWinterHardinessFilter(7, 9)
	t.WithNewStep("Filter by hardiness range 7-9", func(ctx provider.StepCtx) {
		assert.False(t, filter.Filter(s.coniferousPlant))
		assert.False(t, filter.Filter(s.deciduousPlant))
	})
}

func (s *PlantFiltersTestSuite) TestPlantSoilTypeFilter(t provider.T) {
	t.Tags("filter", "soil", "type")
	t.Description("Test plant soil type filtering")

	filter := NewSoilTypeFilter([]plant.Soil{plant.MediumSoil})
	t.WithNewStep("Filter by medium soil", func(ctx provider.StepCtx) {
		assert.True(t, filter.Filter(s.coniferousPlant))
		assert.True(t, filter.Filter(s.deciduousPlant))
	})

	filter = NewSoilTypeFilter([]plant.Soil{plant.LightSoil})
	t.WithNewStep("Filter by light soil", func(ctx provider.StepCtx) {
		assert.False(t, filter.Filter(s.coniferousPlant))
		assert.False(t, filter.Filter(s.deciduousPlant))
	})
}

func (s *PlantFiltersTestSuite) TestPlantFloweringPeriodFilter(t provider.T) {
	t.Tags("filter", "flowering")
	t.Description("Test plant flowering period filtering")

	filter := NewFloweringPeriodFilter([]plant.FloweringPeriod{plant.Spring})
	t.WithNewStep("Filter by spring flowering", func(ctx provider.StepCtx) {
		assert.False(t, filter.Filter(s.coniferousPlant))
		assert.True(t, filter.Filter(s.deciduousPlant))
	})

	filter = NewFloweringPeriodFilter([]plant.FloweringPeriod{plant.Summer})
	t.WithNewStep("Filter by summer flowering", func(ctx provider.StepCtx) {
		assert.False(t, filter.Filter(s.deciduousPlant))
	})
}

func TestPlantFilters(t *testing.T) {
	suite.RunSuite(t, new(PlantFiltersTestSuite))
}
