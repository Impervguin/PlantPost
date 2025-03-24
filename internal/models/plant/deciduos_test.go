package plant

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeciduosImplementing(t *testing.T) {
	assert.Implements(t, (*PlantSpecification)(nil), new(DeciduousSpecification))
}

func TestDeciduousSpecification(t *testing.T) {
	// Валидные параметры для тестов
	validHeight := 8.2
	validDiameter := 1.8
	validFloweringPeriod := Spring
	validSoilAcidity := SoilAcidity(6)
	validSoilMoisture := MediumMoisture
	validLightRelation := HalfShadow
	validSoilType := MediumSoil
	validWinterHardiness := WinterHardiness(5)

	t.Run("NewDeciduousSpecification - успешное создание", func(t *testing.T) {
		spec, err := NewDeciduousSpecification(
			validHeight,
			validDiameter,
			validFloweringPeriod,
			validSoilAcidity,
			validSoilMoisture,
			validLightRelation,
			validSoilType,
			validWinterHardiness,
		)

		require.NoError(t, err)
		assert.Equal(t, validHeight, spec.GetHeightM())
		assert.Equal(t, validDiameter, spec.GetDiameterM())
		assert.Equal(t, validFloweringPeriod, spec.GetFloweringPeriod())
		assert.Equal(t, validSoilAcidity, spec.GetSoilAcidity())
		assert.Equal(t, validSoilMoisture, spec.GetSoilMoisture())
		assert.Equal(t, validLightRelation, spec.GetLightRelation())
		assert.Equal(t, validSoilType, spec.GetSoilType())
		assert.Equal(t, validWinterHardiness, spec.GetWinterHardiness())
	})

	t.Run("NewDeciduousSpecification - ошибки валидации", func(t *testing.T) {
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
				name:            "невалидная высота",
				height:          0,
				diameter:        validDiameter,
				floweringPeriod: validFloweringPeriod,
				soilAcidity:     validSoilAcidity,
				soilMoisture:    validSoilMoisture,
				lightRelation:   validLightRelation,
				soilType:        validSoilType,
				winterHardiness: validWinterHardiness,
				expectError:     true,
			},
			{
				name:            "невалидный диаметр",
				height:          validHeight,
				diameter:        0,
				floweringPeriod: validFloweringPeriod,
				soilAcidity:     validSoilAcidity,
				soilMoisture:    validSoilMoisture,
				lightRelation:   validLightRelation,
				soilType:        validSoilType,
				winterHardiness: validWinterHardiness,
				expectError:     true,
			},
			{
				name:            "валидные параметры",
				height:          validHeight,
				diameter:        validDiameter,
				floweringPeriod: validFloweringPeriod,
				soilAcidity:     validSoilAcidity,
				soilMoisture:    validSoilMoisture,
				lightRelation:   validLightRelation,
				soilType:        validSoilType,
				winterHardiness: validWinterHardiness,
				expectError:     false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := NewDeciduousSpecification(
					tc.height,
					tc.diameter,
					tc.floweringPeriod,
					tc.soilAcidity,
					tc.soilMoisture,
					tc.lightRelation,
					tc.soilType,
					tc.winterHardiness,
				)

				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Validate - проверка валидации", func(t *testing.T) {
		spec := &DeciduousSpecification{
			heightM:         validHeight,
			diameterM:       validDiameter,
			floweringPeriod: validFloweringPeriod,
			soilAcidity:     validSoilAcidity,
			soilMoisture:    validSoilMoisture,
			lightRelation:   validLightRelation,
			soilType:        validSoilType,
			winterHardiness: validWinterHardiness,
		}

		assert.NoError(t, spec.Validate())
	})

	t.Run("Getters - проверка геттеров", func(t *testing.T) {
		spec := &DeciduousSpecification{
			heightM:         validHeight,
			diameterM:       validDiameter,
			floweringPeriod: validFloweringPeriod,
			soilAcidity:     validSoilAcidity,
			soilMoisture:    validSoilMoisture,
			lightRelation:   validLightRelation,
			soilType:        validSoilType,
			winterHardiness: validWinterHardiness,
		}

		assert.Equal(t, validHeight, spec.GetHeightM())
		assert.Equal(t, validDiameter, spec.GetDiameterM())
		assert.Equal(t, validFloweringPeriod, spec.GetFloweringPeriod())
		assert.Equal(t, validSoilAcidity, spec.GetSoilAcidity())
		assert.Equal(t, validSoilMoisture, spec.GetSoilMoisture())
		assert.Equal(t, validLightRelation, spec.GetLightRelation())
		assert.Equal(t, validSoilType, spec.GetSoilType())
		assert.Equal(t, validWinterHardiness, spec.GetWinterHardiness())
	})
}

func TestPlantWithDeciduousSpecification(t *testing.T) {
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

	t.Run("Создание растения с лиственной спецификацией", func(t *testing.T) {
		plant, err := NewPlant(
			"Дуб",
			"Quercus robur",
			"Дуб обыкновенный",
			uuid.New(),
			*NewPlantPhotos(),
			DeciduousCategory,
			spec,
		)

		require.NoError(t, err)
		assert.Equal(t, "Дуб", plant.GetName())
		assert.Equal(t, DeciduousCategory, plant.GetCategory())

		plantSpec, ok := plant.GetSpecification().(*DeciduousSpecification)
		require.True(t, ok)
		assert.Equal(t, 12.0, plantSpec.GetHeightM())
		assert.Equal(t, Spring, plantSpec.GetFloweringPeriod())
		assert.Equal(t, 5, int(plantSpec.GetWinterHardiness()))
	})
}
