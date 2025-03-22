package plant_test

import (
	"PlantSite/internal/models/plant"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPlantPhotoCreation(t *testing.T) {
	testID := 1
	t.Logf("Test %d: creating a new photo", testID)
	{
		descr := "some description"
		photo, err := plant.CreatePlantPhoto(uuid.New(), uuid.New(), descr)
		assert.NoError(t, err)
		assert.NotNil(t, photo)
	}
	testID++
	t.Logf("Test %d: creating a new photo with an empty file ID", testID)
	{
		photo, err := plant.CreatePlantPhoto(uuid.Nil, uuid.New(), "some description")
		assert.Error(t, err)
		assert.Nil(t, photo)
	}
	testID++
	t.Logf("Test %d: creating a new photo with an empty description", testID)
	{
		photo, err := plant.CreatePlantPhoto(uuid.New(), uuid.New(), "")
		assert.NoError(t, err)
		assert.NotNil(t, photo)
	}
	testID++
	t.Logf("Test %d: creating a new photo with an invalid file ID", testID)
	{
		photo, err := plant.CreatePlantPhoto(uuid.New(), uuid.Nil, "some description")
		assert.Error(t, err)
		assert.Nil(t, photo)
	}
}
