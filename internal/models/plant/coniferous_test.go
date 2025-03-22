package plant_test

import (
	"PlantSite/internal/models/plant"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConiferousImplementing(t *testing.T) {
	assert.Implements(t, (*plant.PlantSpecification)(nil), new(plant.ConiferousSpecification))
}

func TestConiferous(t *testing.T) {
	defaultHeightM := 10.
	defaultDiameterM := 20.
	defaultHardiness := plant.WinterHardiness(5)
	defaultLightRelation := plant.HalfShadow
	defaultSoilAcidity := plant.SoilAcidity(5)
	defaultSoilMoisture := plant.MediumMoisture
	defaultSoilType := plant.MediumSoil
	testID := 0

	t.Logf("Test %d: default values", testID)
	{
		spec, err := plant.NewConiferousSpecification(
			defaultHeightM,
			defaultDiameterM,
			defaultSoilAcidity,
			defaultSoilMoisture,
			defaultLightRelation,
			defaultSoilType,
			defaultHardiness,
		)
		assert.Nil(t, err)
		assert.NotNil(t, spec)
		assert.Equal(t, defaultHeightM, spec.GetHeightM())
		assert.Equal(t, defaultDiameterM, spec.GetDiameterM())
		assert.Equal(t, defaultSoilAcidity, spec.GetSoilAcidity())
		assert.Equal(t, defaultSoilMoisture, spec.GetSoilMoisture())
		assert.Equal(t, defaultLightRelation, spec.GetLightRelation())
		assert.Equal(t, defaultSoilType, spec.GetSoilType())
		assert.Equal(t, defaultHardiness, spec.GetWinterHardiness())
		assert.Nil(t, spec.Validate())
	}
	testID++
	t.Logf("Test %d: soil types", testID)
	{
		values := []plant.Soil{
			plant.LightSoil,
			plant.MediumSoil,
			plant.HeavySoil,
		}
		for _, soil := range values {
			spec, err := plant.NewConiferousSpecification(
				defaultHeightM,
				defaultDiameterM,
				defaultSoilAcidity,
				defaultSoilMoisture,
				defaultLightRelation,
				soil,
				defaultHardiness,
			)
			assert.Nil(t, err)
			assert.Equal(t, soil, spec.GetSoilType())
		}
		spec, err := plant.NewConiferousSpecification(
			defaultHeightM,
			defaultDiameterM,
			defaultSoilAcidity,
			defaultSoilMoisture,
			defaultLightRelation,
			"invalid",
			defaultHardiness,
		)
		assert.NotNil(t, err)
		assert.Nil(t, spec)
	}
	testID++
	t.Logf("Test %d: soil moistures", testID)
	{
		values := []plant.SoilMoisture{
			plant.DryMoisture,
			plant.LowMoisture,
			plant.MediumMoisture,
			plant.HighMoisture,
		}
		for _, moisture := range values {
			spec, err := plant.NewConiferousSpecification(
				defaultHeightM,
				defaultDiameterM,
				defaultSoilAcidity,
				moisture,
				defaultLightRelation,
				defaultSoilType,
				defaultHardiness,
			)
			assert.Nil(t, err)
			assert.Equal(t, moisture, spec.GetSoilMoisture())
		}
	}
	testID++
	t.Logf("Test %d: soil acidities", testID)
	{
		for acidity := range 10 {
			spec, err := plant.NewConiferousSpecification(
				defaultHeightM,
				defaultDiameterM,
				plant.SoilAcidity(acidity+1),
				defaultSoilMoisture,
				defaultLightRelation,
				defaultSoilType,
				defaultHardiness,
			)
			assert.Nil(t, err)
			assert.Equal(t, plant.SoilAcidity(acidity+1), spec.GetSoilAcidity())
		}
		spec, err := plant.NewConiferousSpecification(
			defaultHeightM,
			defaultDiameterM,
			plant.SoilAcidity(0),
			defaultSoilMoisture,
			defaultLightRelation,
			defaultSoilType,
			defaultHardiness,
		)
		assert.NotNil(t, err)
		assert.Nil(t, spec)
	}
	testID++
	t.Logf("Test %d: light relations", testID)
	{
		values := []plant.LightRelation{
			plant.Shadow,
			plant.HalfShadow,
			plant.Light,
		}
		for _, relation := range values {
			spec, err := plant.NewConiferousSpecification(
				defaultHeightM,
				defaultDiameterM,
				defaultSoilAcidity,
				defaultSoilMoisture,
				relation,
				defaultSoilType,
				defaultHardiness,
			)
			assert.Nil(t, err)
			assert.Equal(t, relation, spec.GetLightRelation())
		}
		spec, err := plant.NewConiferousSpecification(
			defaultHeightM,
			defaultDiameterM,
			defaultSoilAcidity,
			defaultSoilMoisture,
			plant.LightRelation("invalid"),
			defaultSoilType,
			defaultHardiness,
		)
		assert.NotNil(t, err)
		assert.Nil(t, spec)
	}
	testID++
	t.Logf("Test %d: winter hardiness", testID)
	{
		for hardiness := range 10 {
			spec, err := plant.NewConiferousSpecification(
				defaultHeightM,
				defaultDiameterM,
				defaultSoilAcidity,
				defaultSoilMoisture,
				defaultLightRelation,
				defaultSoilType,
				plant.WinterHardiness(hardiness+1),
			)
			assert.Nil(t, err)
			assert.Equal(t, plant.WinterHardiness(hardiness+1), spec.GetWinterHardiness())
		}
		spec, err := plant.NewConiferousSpecification(
			defaultHeightM,
			defaultDiameterM,
			defaultSoilAcidity,
			defaultSoilMoisture,
			defaultLightRelation,
			defaultSoilType,
			plant.WinterHardiness(-1),
		)
		assert.NotNil(t, err)
		assert.Nil(t, spec)
	}
	testID++
	t.Logf("Test %d: height negative", testID)
	{
		spec, err := plant.NewConiferousSpecification(-1, defaultDiameterM, defaultSoilAcidity, defaultSoilMoisture, defaultLightRelation, defaultSoilType, defaultHardiness)
		assert.NotNil(t, err)
		assert.Nil(t, spec)
	}

	testID++
	t.Logf("Test %d: diameter negative", testID)
	{
		spec, err := plant.NewConiferousSpecification(defaultHeightM, -1, defaultSoilAcidity, defaultSoilMoisture, defaultLightRelation, defaultSoilType, defaultHardiness)
		assert.NotNil(t, err)
		assert.Nil(t, spec)
	}
}
