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

type ConiferousSpecificationTestSuite struct {
	suite.Suite
	validHeight          float64
	validDiameter        float64
	validSoilAcidity     SoilAcidity
	validSoilMoisture    SoilMoisture
	validLightRelation   LightRelation
	validSoilType        Soil
	validWinterHardiness WinterHardiness
}

func (s *ConiferousSpecificationTestSuite) BeforeEach(t provider.T) {
	t.Epic("Plant Specifications")
	t.Feature("Coniferous Specification")

	// Valid parameters for tests
	s.validHeight = 10.5
	s.validDiameter = 2.3
	s.validSoilAcidity = SoilAcidity(5)
	s.validSoilMoisture = MediumMoisture
	s.validLightRelation = HalfShadow
	s.validSoilType = MediumSoil
	s.validWinterHardiness = WinterHardiness(6)
}

func (s *ConiferousSpecificationTestSuite) TestConiferousImplementsPlantSpecification(t provider.T) {
	t.Tags("interface", "implementation")
	t.Description("Verify that ConiferousSpecification implements PlantSpecification interface")

	// Assert
	t.WithNewStep("Check interface implementation", func(ctx provider.StepCtx) {
		assert.Implements(t, (*PlantSpecification)(nil), new(ConiferousSpecification))
	})
}

func (s *ConiferousSpecificationTestSuite) TestNewConiferousSpecificationSuccess(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful creation of ConiferousSpecification with valid parameters")

	// Act
	var spec *ConiferousSpecification
	var err error

	t.WithNewStep("Create coniferous specification", func(ctx provider.StepCtx) {
		spec, err = NewConiferousSpecification(
			s.validHeight,
			s.validDiameter,
			s.validSoilAcidity,
			s.validSoilMoisture,
			s.validLightRelation,
			s.validSoilType,
			s.validWinterHardiness,
		)
	})

	// Assert
	t.WithNewStep("Verify specification properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, s.validHeight, spec.GetHeightM())
		assert.Equal(t, s.validDiameter, spec.GetDiameterM())
		assert.Equal(t, s.validSoilAcidity, spec.GetSoilAcidity())
		assert.Equal(t, s.validSoilMoisture, spec.GetSoilMoisture())
		assert.Equal(t, s.validLightRelation, spec.GetLightRelation())
		assert.Equal(t, s.validSoilType, spec.GetSoilType())
		assert.Equal(t, s.validWinterHardiness, spec.GetWinterHardiness())
	})
}

func (s *ConiferousSpecificationTestSuite) TestNewConiferousSpecificationValidationErrors(t provider.T) {
	t.Tags("negative", "validation")
	t.Description("Test validation errors during coniferous specification creation")

	testCases := []struct {
		name            string
		height          float64
		diameter        float64
		soilAcidity     SoilAcidity
		soilMoisture    SoilMoisture
		lightRelation   LightRelation
		soilType        Soil
		winterHardiness WinterHardiness
		expectError     bool
	}{
		{
			name:            "Invalid height (zero value)",
			height:          0,
			diameter:        s.validDiameter,
			soilAcidity:     s.validSoilAcidity,
			soilMoisture:    s.validSoilMoisture,
			lightRelation:   s.validLightRelation,
			soilType:        s.validSoilType,
			winterHardiness: s.validWinterHardiness,
			expectError:     true,
		},
		{
			name:            "Invalid diameter (zero value)",
			height:          s.validHeight,
			diameter:        0,
			soilAcidity:     s.validSoilAcidity,
			soilMoisture:    s.validSoilMoisture,
			lightRelation:   s.validLightRelation,
			soilType:        s.validSoilType,
			winterHardiness: s.validWinterHardiness,
			expectError:     true,
		},
		{
			name:            "Invalid height (negative value)",
			height:          -5.0,
			diameter:        s.validDiameter,
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

			// Arrange
			t.WithNewStep("Prepare test parameters", func(ctx provider.StepCtx) {
				ctx.WithParameters(
					allure.NewParameter("Height", tc.height),
					allure.NewParameter("Diameter", tc.diameter),
					allure.NewParameter("ExpectError", tc.expectError),
				)
			})

			// Act
			var err error
			t.WithNewStep("Attempt to create specification", func(ctx provider.StepCtx) {
				_, err = NewConiferousSpecification(
					tc.height,
					tc.diameter,
					tc.soilAcidity,
					tc.soilMoisture,
					tc.lightRelation,
					tc.soilType,
					tc.winterHardiness,
				)
			})

			// Assert
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

func (s *ConiferousSpecificationTestSuite) TestValidateSuccess(t provider.T) {
	t.Tags("positive", "validation")
	t.Description("Successful validation of ConiferousSpecification")

	// Arrange
	spec := &ConiferousSpecification{
		heightM:         s.validHeight,
		diameterM:       s.validDiameter,
		soilAcidity:     s.validSoilAcidity,
		soilMoisture:    s.validSoilMoisture,
		lightRelation:   s.validLightRelation,
		soilType:        s.validSoilType,
		winterHardiness: s.validWinterHardiness,
	}

	// Act
	var validationErr error
	t.WithNewStep("Validate specification", func(ctx provider.StepCtx) {
		validationErr = spec.Validate()
	})

	// Assert
	t.WithNewStep("Verify validation success", func(ctx provider.StepCtx) {
		assert.NoError(t, validationErr)
	})
}

func (s *ConiferousSpecificationTestSuite) TestGettersReturnCorrectValues(t provider.T) {
	t.Tags("positive", "getters")
	t.Description("Test that ConiferousSpecification getters return correct values")

	// Arrange
	spec := &ConiferousSpecification{
		heightM:         s.validHeight,
		diameterM:       s.validDiameter,
		soilAcidity:     s.validSoilAcidity,
		soilMoisture:    s.validSoilMoisture,
		lightRelation:   s.validLightRelation,
		soilType:        s.validSoilType,
		winterHardiness: s.validWinterHardiness,
	}

	// Act & Assert
	t.WithNewStep("Verify all getter methods", func(ctx provider.StepCtx) {
		assert.Equal(t, s.validHeight, spec.GetHeightM(), "Height getter should return correct value")
		assert.Equal(t, s.validDiameter, spec.GetDiameterM(), "Diameter getter should return correct value")
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

func (s *ConiferousSpecificationTestSuite) TestPlantWithConiferousSpecification(t provider.T) {
	t.Tags("integration", "creation")
	t.Description("Integration test: creating a plant with coniferous specification")

	// Arrange
	spec, err := NewConiferousSpecification(
		15.0,
		3.0,
		6,
		MediumMoisture,
		HalfShadow,
		MediumSoil,
		5,
	)
	require.NoError(t, err)

	// Act
	var plant *Plant
	t.WithNewStep("Create plant with coniferous specification", func(ctx provider.StepCtx) {
		plant, err = NewPlant(
			"Pine",
			"Pinus sylvestris",
			"Scots pine",
			uuid.New(),
			*NewPlantPhotos(),
			ConiferousCategory,
			spec,
		)
	})

	// Assert
	t.WithNewStep("Verify plant properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, "Pine", plant.GetName())
		assert.Equal(t, "Pinus sylvestris", plant.GetLatinName())
		assert.Equal(t, ConiferousCategory, plant.GetCategory())

		plantSpec, ok := plant.GetSpecification().(*ConiferousSpecification)
		require.True(t, ok, "Specification should be of type ConiferousSpecification")
		assert.Equal(t, 15.0, plantSpec.GetHeightM(), "Height should match")
		assert.Equal(t, 3.0, plantSpec.GetDiameterM(), "Diameter should match")
		assert.Equal(t, WinterHardiness(5), plantSpec.GetWinterHardiness(), "Winter hardiness should match")
	})
}

func TestConiferousSpecification(t *testing.T) {
	suite.RunSuite(t, new(ConiferousSpecificationTestSuite))
}
