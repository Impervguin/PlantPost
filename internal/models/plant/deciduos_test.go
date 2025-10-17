//go:build unit

package plant

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DeciduousSpecificationTestSuite struct {
	suite.Suite
	validHeight          float64
	validDiameter        float64
	validFloweringPeriod FloweringPeriod
	validSoilAcidity     SoilAcidity
	validSoilMoisture    SoilMoisture
	validLightRelation   LightRelation
	validSoilType        Soil
	validWinterHardiness WinterHardiness
}

func (s *DeciduousSpecificationTestSuite) BeforeEach(t provider.T) {
	t.Epic("Plant Specifications")
	t.Feature("Deciduous Specification")

	s.validHeight = 8.2
	s.validDiameter = 1.8
	s.validFloweringPeriod = Spring
	s.validSoilAcidity = SoilAcidity(6)
	s.validSoilMoisture = MediumMoisture
	s.validLightRelation = HalfShadow
	s.validSoilType = MediumSoil
	s.validWinterHardiness = WinterHardiness(5)
}

func (s *DeciduousSpecificationTestSuite) TestDeciduosImplementsPlantSpecification(t provider.T) {
	t.Tags("interface", "implementation")
	t.Description("Verify that DeciduousSpecification implements PlantSpecification interface")

	assert.Implements(t, (*PlantSpecification)(nil), new(DeciduousSpecification))
}

func (s *DeciduousSpecificationTestSuite) TestNewDeciduousSpecificationSuccess(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful creation of DeciduousSpecification with valid parameters")

	var spec *DeciduousSpecification
	var err error

	t.WithNewStep("Create deciduous specification", func(ctx provider.StepCtx) {
		spec, err = NewDeciduousSpecification(
			s.validHeight,
			s.validDiameter,
			s.validFloweringPeriod,
			s.validSoilAcidity,
			s.validSoilMoisture,
			s.validLightRelation,
			s.validSoilType,
			s.validWinterHardiness,
		)
	})

	t.WithNewStep("Verify specification properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, s.validHeight, spec.GetHeightM())
		assert.Equal(t, s.validDiameter, spec.GetDiameterM())
		assert.Equal(t, s.validFloweringPeriod, spec.GetFloweringPeriod())
		assert.Equal(t, s.validSoilAcidity, spec.GetSoilAcidity())
		assert.Equal(t, s.validSoilMoisture, spec.GetSoilMoisture())
		assert.Equal(t, s.validLightRelation, spec.GetLightRelation())
		assert.Equal(t, s.validSoilType, spec.GetSoilType())
		assert.Equal(t, s.validWinterHardiness, spec.GetWinterHardiness())
	})
}

func (s *DeciduousSpecificationTestSuite) TestNewDeciduousSpecificationValidationErrors(t provider.T) {
	t.Tags("negative", "validation")
	t.Description("Test validation errors during deciduous specification creation")

	testCases := []struct {
		name            string
		height          float64
		diameter        float64
		floweringPeriod FloweringPeriod
		soilAcidity     SoilAcidity
		soilMoisture    SoilMoisture
		lightRelation   LightRelation
		soilType        Soil
		winterHardiness WinterHardiness
		expectError     bool
	}{
		{
			name:            "Invalid height",
			height:          0,
			diameter:        s.validDiameter,
			floweringPeriod: s.validFloweringPeriod,
			soilAcidity:     s.validSoilAcidity,
			soilMoisture:    s.validSoilMoisture,
			lightRelation:   s.validLightRelation,
			soilType:        s.validSoilType,
			winterHardiness: s.validWinterHardiness,
			expectError:     true,
		},
		{
			name:            "Invalid diameter",
			height:          s.validHeight,
			diameter:        0,
			floweringPeriod: s.validFloweringPeriod,
			soilAcidity:     s.validSoilAcidity,
			soilMoisture:    s.validSoilMoisture,
			lightRelation:   s.validLightRelation,
			soilType:        s.validSoilType,
			winterHardiness: s.validWinterHardiness,
			expectError:     true,
		},
		{
			name:            "Valid parameters",
			height:          s.validHeight,
			diameter:        s.validDiameter,
			floweringPeriod: s.validFloweringPeriod,
			soilAcidity:     s.validSoilAcidity,
			soilMoisture:    s.validSoilMoisture,
			lightRelation:   s.validLightRelation,
			soilType:        s.validSoilType,
			winterHardiness: s.validWinterHardiness,
			expectError:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t provider.T) {
			t.Tags("negative", "validation")

			t.WithNewStep("Prepare test parameters", func(ctx provider.StepCtx) {
				ctx.WithParameters(
					allure.NewParameter("Height", tc.height),
					allure.NewParameter("Diameter", tc.diameter),
					allure.NewParameter("ExpectError", tc.expectError),
				)
			})

			var err error
			t.WithNewStep("Attempt to create specification", func(ctx provider.StepCtx) {
				_, err = NewDeciduousSpecification(
					tc.height,
					tc.diameter,
					tc.floweringPeriod,
					tc.soilAcidity,
					tc.soilMoisture,
					tc.lightRelation,
					tc.soilType,
					tc.winterHardiness,
				)
			})

			t.WithNewStep("Verify validation result", func(ctx provider.StepCtx) {
				if tc.expectError {
					assert.Error(t, err, "Expected validation error for: %s", tc.name)
				} else {
					assert.NoError(t, err, "Unexpected error for valid parameters: %s", tc.name)
				}
			})
		})
	}
}

func (s *DeciduousSpecificationTestSuite) TestValidateSuccess(t provider.T) {
	t.Tags("positive", "validation")
	t.Description("Successful validation of DeciduousSpecification")

	spec := &DeciduousSpecification{
		heightM:         s.validHeight,
		diameterM:       s.validDiameter,
		floweringPeriod: s.validFloweringPeriod,
		soilAcidity:     s.validSoilAcidity,
		soilMoisture:    s.validSoilMoisture,
		lightRelation:   s.validLightRelation,
		soilType:        s.validSoilType,
		winterHardiness: s.validWinterHardiness,
	}

	var validationErr error
	t.WithNewStep("Validate specification", func(ctx provider.StepCtx) {
		validationErr = spec.Validate()
	})

	t.WithNewStep("Verify validation success", func(ctx provider.StepCtx) {
		assert.NoError(t, validationErr)
	})
}

func (s *DeciduousSpecificationTestSuite) TestGettersReturnCorrectValues(t provider.T) {
	t.Tags("positive", "getters")
	t.Description("Test that DeciduousSpecification getters return correct values")

	spec := &DeciduousSpecification{
		heightM:         s.validHeight,
		diameterM:       s.validDiameter,
		floweringPeriod: s.validFloweringPeriod,
		soilAcidity:     s.validSoilAcidity,
		soilMoisture:    s.validSoilMoisture,
		lightRelation:   s.validLightRelation,
		soilType:        s.validSoilType,
		winterHardiness: s.validWinterHardiness,
	}

	t.WithNewStep("Verify all getter methods", func(ctx provider.StepCtx) {
		assert.Equal(t, s.validHeight, spec.GetHeightM(), "Height getter should return correct value")
		assert.Equal(t, s.validDiameter, spec.GetDiameterM(), "Diameter getter should return correct value")
		assert.Equal(
			t,
			s.validFloweringPeriod,
			spec.GetFloweringPeriod(),
			"Flowering period getter should return correct value",
		)
		assert.Equal(t, s.validSoilAcidity, spec.GetSoilAcidity(), "Soil acidity getter should return correct value")
		assert.Equal(t, s.validSoilMoisture, spec.GetSoilMoisture(), "Soil moisture getter should return correct value")
		assert.Equal(
			t,
			s.validLightRelation,
			spec.GetLightRelation(),
			"Light relation getter should return correct value",
		)
		assert.Equal(t, s.validSoilType, spec.GetSoilType(), "Soil type getter should return correct value")
		assert.Equal(
			t,
			s.validWinterHardiness,
			spec.GetWinterHardiness(),
			"Winter hardiness getter should return correct value",
		)
	})
}

func (s *DeciduousSpecificationTestSuite) TestPlantWithDeciduousSpecification(t provider.T) {
	t.Tags("integration", "creation")
	t.Description("Integration test: creating a plant with deciduous specification")

	spec, err := NewDeciduousSpecification(
		12.0,
		2.5,
		Spring,
		6,
		MediumMoisture,
		HalfShadow,
		MediumSoil,
		5,
	)
	require.NoError(t, err)

	var plant *Plant
	t.WithNewStep("Create plant with deciduous specification", func(ctx provider.StepCtx) {
		plant, err = NewPlant(
			"Oak",
			"Quercus robur",
			"Common oak",
			uuid.New(),
			*NewPlantPhotos(),
			DeciduousCategory,
			spec,
		)
	})

	t.WithNewStep("Verify plant properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, "Oak", plant.GetName())
		assert.Equal(t, "Quercus robur", plant.GetLatinName())
		assert.Equal(t, DeciduousCategory, plant.GetCategory())

		plantSpec, ok := plant.GetSpecification().(*DeciduousSpecification)
		require.True(t, ok, "Specification should be of type DeciduousSpecification")
		assert.Equal(t, 12.0, plantSpec.GetHeightM(), "Height should match")
		assert.Equal(t, Spring, plantSpec.GetFloweringPeriod(), "Flowering period should match")
		assert.Equal(t, WinterHardiness(5), plantSpec.GetWinterHardiness(), "Winter hardiness should match")
	})
}

func TestDeciduousSpecification(t *testing.T) {
	suite.RunSuite(t, new(DeciduousSpecificationTestSuite))
}
