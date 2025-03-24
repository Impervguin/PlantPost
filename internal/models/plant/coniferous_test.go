package plant

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConiferousImplementing(t *testing.T) {
	assert.Implements(t, (*PlantSpecification)(nil), new(ConiferousSpecification))
}

func TestConiferousSpecification(t *testing.T) {
	// Валидные параметры для тестов
	validHeight := 10.5
	validDiameter := 2.3
	validSoilAcidity := SoilAcidity(5)
	validSoilMoisture := MediumMoisture
	validLightRelation := HalfShadow
	validSoilType := MediumSoil
	validWinterHardiness := WinterHardiness(6)

	t.Run("NewConiferousSpecification - успешное создание", func(t *testing.T) {
		spec, err := NewConiferousSpecification(
			validHeight,
			validDiameter,
			validSoilAcidity,
			validSoilMoisture,
			validLightRelation,
			validSoilType,
			validWinterHardiness,
		)

		require.NoError(t, err)
		assert.Equal(t, validHeight, spec.GetHeightM())
		assert.Equal(t, validDiameter, spec.GetDiameterM())
		assert.Equal(t, validSoilAcidity, spec.GetSoilAcidity())
		assert.Equal(t, validSoilMoisture, spec.GetSoilMoisture())
		assert.Equal(t, validLightRelation, spec.GetLightRelation())
		assert.Equal(t, validSoilType, spec.GetSoilType())
		assert.Equal(t, validWinterHardiness, spec.GetWinterHardiness())
	})

	t.Run("NewConiferousSpecification - ошибки валидации", func(t *testing.T) {
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
				name:            "невалидная высота",
				height:          0,
				diameter:        validDiameter,
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
				_, err := NewConiferousSpecification(
					tc.height,
					tc.diameter,
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
		spec := &ConiferousSpecification{
			heightM:         validHeight,
			diameterM:       validDiameter,
			soilAcidity:     validSoilAcidity,
			soilMoisture:    validSoilMoisture,
			lightRelation:   validLightRelation,
			soilType:        validSoilType,
			winterHardiness: validWinterHardiness,
		}

		assert.NoError(t, spec.Validate())
	})

	t.Run("Getters - проверка геттеров", func(t *testing.T) {
		spec := &ConiferousSpecification{
			heightM:         validHeight,
			diameterM:       validDiameter,
			soilAcidity:     validSoilAcidity,
			soilMoisture:    validSoilMoisture,
			lightRelation:   validLightRelation,
			soilType:        validSoilType,
			winterHardiness: validWinterHardiness,
		}

		assert.Equal(t, validHeight, spec.GetHeightM())
		assert.Equal(t, validDiameter, spec.GetDiameterM())
		assert.Equal(t, validSoilAcidity, spec.GetSoilAcidity())
		assert.Equal(t, validSoilMoisture, spec.GetSoilMoisture())
		assert.Equal(t, validLightRelation, spec.GetLightRelation())
		assert.Equal(t, validSoilType, spec.GetSoilType())
		assert.Equal(t, validWinterHardiness, spec.GetWinterHardiness())
	})
}

func TestPlantWithConiferousSpecification(t *testing.T) {
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

	t.Run("Создание растения с хвойной спецификацией", func(t *testing.T) {
		plant, err := NewPlant(
			"Сосна",
			"Pinus sylvestris",
			"Обыкновенная сосна",
			uuid.New(),
			*NewPlantPhotos(),
			ConiferousCategory,
			spec,
		)

		require.NoError(t, err)
		assert.Equal(t, "Сосна", plant.GetName())
		assert.Equal(t, ConiferousCategory, plant.GetCategory())

		plantSpec, ok := plant.GetSpecification().(*ConiferousSpecification)
		require.True(t, ok)
		assert.Equal(t, 15.0, plantSpec.GetHeightM())
		assert.Equal(t, 5, int(plantSpec.GetWinterHardiness()))
	})
}
