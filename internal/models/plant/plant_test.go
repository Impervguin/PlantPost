//go:build unit

package plant

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockPlantSpecification for testing
type MockPlantSpecification struct {
	Invalid bool
}

func (m *MockPlantSpecification) Validate() error {
	if !m.Invalid {
		return nil
	}
	return fmt.Errorf("invalid specification")
}

func (m *MockPlantSpecification) Category() string {
	return ""
}

type PlantTestSuite struct {
	suite.Suite
	validID     uuid.UUID
	validFileID uuid.UUID
	validTime   time.Time
	validSpec   *MockPlantSpecification
	validPhoto  *PlantPhoto
	validPhotos *PlantPhotos
}

func (s *PlantTestSuite) BeforeEach(t provider.T) {
	t.Epic("Plant Management")
	t.Feature("Plant Model Operations")

	// Prepare test data
	s.validID = uuid.New()
	s.validFileID = uuid.New()
	s.validTime = time.Now().Add(-time.Hour)
	s.validSpec = &MockPlantSpecification{}
	s.validPhoto, _ = CreatePlantPhoto(uuid.New(), s.validFileID, "test photo")
	s.validPhotos = NewPlantPhotos()
	_ = s.validPhotos.Add(s.validPhoto)
}

func (s *PlantTestSuite) TestCreatePlant_Success(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful plant creation with valid data")

	// Arrange
	t.WithNewStep("Prepare test data", func(ctx provider.StepCtx) {
		ctx.WithNewStep("Generate valid plant data", func(ctx provider.StepCtx) {})
	})

	// Act
	var plant *Plant
	var err error

	t.WithNewStep("Create plant instance", func(ctx provider.StepCtx) {
		plant, err = CreatePlant(
			s.validID,
			"Test Plant",
			"Testus Plantus",
			"Test description",
			s.validFileID,
			*s.validPhotos,
			"Test Category",
			s.validSpec,
			s.validTime,
			s.validTime,
		)
	})

	// Assert
	t.WithNewStep("Verify plant creation results", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.NotNil(t, plant)
		assert.Equal(t, s.validID, plant.ID())
		assert.Equal(t, s.validFileID, plant.MainPhotoID())
		assert.Equal(t, *s.validPhotos, plant.GetPhotos())
		assert.Equal(t, "Test Category", plant.GetCategory())
		assert.Equal(t, s.validSpec, plant.GetSpecification())
		assert.Equal(t, s.validTime, plant.CreatedAt())
		assert.Equal(t, s.validFileID, plant.MainPhotoID())
		assert.Equal(t, "Test Plant", plant.GetName())
		assert.Equal(t, "Testus Plantus", plant.GetLatinName())
		assert.Equal(t, "Test description", plant.GetDescription())
	})
}

func (s *PlantTestSuite) TestCreatePlant_ValidationErrors(t provider.T) {
	t.Tags("negative", "validation")
	t.Description("Test validation errors during plant creation")

	testCases := []struct {
		name        string
		id          uuid.UUID
		nameStr     string
		latinName   string
		description string
		mainPhotoID uuid.UUID
		photos      PlantPhotos
		category    string
		spec        PlantSpecification
		createdAt   time.Time
		updatedAt   time.Time
		expectError bool
	}{
		{
			name:        "Empty UUID should cause validation error",
			id:          uuid.Nil,
			nameStr:     "Test",
			latinName:   "Testus",
			description: "Desc",
			mainPhotoID: s.validFileID,
			photos:      *s.validPhotos,
			category:    "Cat",
			spec:        s.validSpec,
			createdAt:   s.validTime,
			updatedAt:   s.validTime,
			expectError: true,
		},
		{
			name:        "Empty plant name should cause validation error",
			id:          s.validID,
			nameStr:     "",
			latinName:   "Testus",
			description: "Desc",
			mainPhotoID: s.validFileID,
			photos:      *s.validPhotos,
			category:    "Cat",
			spec:        s.validSpec,
			createdAt:   s.validTime,
			updatedAt:   s.validTime,
			expectError: true,
		},
		{
			name:        "Future creation date should cause validation error",
			id:          s.validID,
			nameStr:     "Test",
			latinName:   "Testus",
			description: "Desc",
			mainPhotoID: s.validFileID,
			photos:      *s.validPhotos,
			category:    "Cat",
			spec:        s.validSpec,
			createdAt:   time.Now().Add(time.Hour),
			updatedAt:   time.Now().Add(time.Hour),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t provider.T) {
			t.Tags("negative", "validation")

			// Arrange
			t.WithNewStep("Prepare invalid test data", func(ctx provider.StepCtx) {
				ctx.WithParameters(
					allure.NewParameter("ID", tc.id.String()),
					allure.NewParameter("Name", tc.nameStr),
					allure.NewParameter("CreatedAt", tc.createdAt.String()),
				)
			})

			// Act
			var err error
			t.WithNewStep("Attempt to create plant with invalid data", func(ctx provider.StepCtx) {
				_, err = CreatePlant(
					tc.id,
					tc.nameStr,
					tc.latinName,
					tc.description,
					tc.mainPhotoID,
					tc.photos,
					tc.category,
					tc.spec,
					tc.createdAt,
					tc.updatedAt,
				)
			})

			// Assert
			t.WithNewStep("Verify validation error occurred", func(ctx provider.StepCtx) {
				if tc.expectError {
					require.Error(t, err, "Expected validation error for: %s", tc.name)
				} else {
					require.NoError(t, err, "Unexpected error for: %s", tc.name)
				}
			})
		})
	}
}

func (s *PlantTestSuite) TestNewPlant_Success(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful plant creation via NewPlant with auto-generated ID")

	// Act
	var plant *Plant
	var err error

	t.WithNewStep("Create plant using NewPlant factory", func(ctx provider.StepCtx) {
		plant, err = NewPlant(
			"Test Plant",
			"Testus Plantus",
			"Test description",
			s.validFileID,
			*s.validPhotos,
			"Test Category",
			s.validSpec,
		)
	})

	// Assert
	t.WithNewStep("Verify auto-generated plant properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.NotNil(t, plant)
		assert.NotEqual(t, uuid.Nil, plant.ID(), "Plant should have non-nil UUID")
		assert.True(t, plant.CreatedAt().Before(time.Now().Add(time.Second)), "Creation time should be in the past")
		assert.True(t, plant.UpdatedAt().Before(time.Now().Add(time.Second)), "Update time should be in the past")
	})
}

func (s *PlantTestSuite) TestUpdateSpecSuccess(t provider.T) {
	t.Tags("positive", "update")
	t.Description("Successful update of plant specification")

	// Arrange
	plant, _ := CreatePlant(
		s.validID,
		"Test Plant",
		"Testus Plantus",
		"Test description",
		s.validFileID,
		*s.validPhotos,
		"Test Category",
		s.validSpec,
		s.validTime,
		s.validTime,
	)

	newSpec := &MockPlantSpecification{}

	// Act
	var err error
	t.WithNewStep("Update plant specification", func(ctx provider.StepCtx) {
		err = plant.UpdateSpec(newSpec)
	})

	// Assert
	t.WithNewStep("Verify specification update success", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, newSpec, plant.GetSpecification(), "Specification should be updated")
	})
}

func (s *PlantTestSuite) TestUpdateSpecValidationError(t provider.T) {
	t.Tags("negative", "update")
	t.Description("Error when updating plant specification with invalid data")

	// Arrange
	plant, _ := CreatePlant(
		s.validID,
		"Test Plant",
		"Testus Plantus",
		"Test description",
		s.validFileID,
		*s.validPhotos,
		"Test Category",
		s.validSpec,
		s.validTime,
		s.validTime,
	)

	newSpec := &MockPlantSpecification{}
	newSpec.Invalid = true

	// Act
	var err error
	t.WithNewStep("Update plant specification with invalid data", func(ctx provider.StepCtx) {
		err = plant.UpdateSpec(newSpec)
	})

	// Assert
	t.WithNewStep("Verify specification update error", func(ctx provider.StepCtx) {
		require.Error(t, err, "Expected error for invalid specification")
	})
}

func (s *PlantTestSuite) TestAddPhotoSuccess(t provider.T) {
	t.Tags("positive", "photos")
	t.Description("Successful addition of photo to plant")

	// Arrange
	emptyPhotos := NewPlantPhotos()
	plant, _ := CreatePlant(
		s.validID,
		"Test Plant",
		"Testus Plantus",
		"Test description",
		s.validFileID,
		*emptyPhotos,
		"Test Category",
		s.validSpec,
		s.validTime,
		s.validTime,
	)

	newPhoto, _ := CreatePlantPhoto(uuid.New(), uuid.New(), "new photo")

	// Act
	var err error
	t.WithNewStep("Add photo to plant", func(ctx provider.StepCtx) {
		err = plant.AddPhoto(newPhoto)
	})

	// Assert
	t.WithNewStep("Verify photo addition success", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		photos := plant.GetPhotos()
		photoFound := false
		photos.Iterate(func(e PlantPhoto) error {
			assert.True(t, e.Compare(newPhoto), "Plant should contain the new photo")
			assert.Equal(t, newPhoto.ID(), e.ID(), "Photo ID should match")
			assert.False(t, photoFound, "Photo should be added only once")
			photoFound = true
			return nil
		})
	})
}

func (s *PlantTestSuite) TestAddPhotoDuplicateError(t provider.T) {
	t.Tags("negative", "photos")
	t.Description("Error when adding duplicate photo to plant")

	// Arrange
	plant, _ := CreatePlant(
		s.validID,
		"Test Plant",
		"Testus Plantus",
		"Test description",
		s.validFileID,
		*s.validPhotos,
		"Test Category",
		s.validSpec,
		s.validTime,
		s.validTime,
	)

	// Act
	var err error
	t.WithNewStep("Attempt to add duplicate photo", func(ctx provider.StepCtx) {
		err = plant.AddPhoto(s.validPhoto)
	})

	// Assert
	t.WithNewStep("Verify duplicate photo error", func(ctx provider.StepCtx) {
		require.Error(t, err, "Should return error for duplicate photo")
	})
}

func (s *PlantTestSuite) TestGetters(t provider.T) {
	t.Tags("positive", "getters")
	t.Description("Test that plant getters return correct values")

	// Arrange
	plant, _ := CreatePlant(
		s.validID,
		"Test Plant",
		"Testus Plantus",
		"Test description",
		s.validFileID,
		*s.validPhotos,
		"Test Category",
		s.validSpec,
		s.validTime,
		s.validTime,
	)

	// Act & Assert
	t.WithNewStep("Verify all getter methods", func(ctx provider.StepCtx) {
		assert.Equal(t, s.validID, plant.ID(), "ID getter should return correct value")
		assert.Equal(t, "Test Plant", plant.GetName(), "Name getter should return correct value")
		assert.Equal(t, "Testus Plantus", plant.GetLatinName(), "Latin name getter should return correct value")
		assert.Equal(t, "Test description", plant.GetDescription(), "Description getter should return correct value")
		assert.Equal(t, s.validFileID, plant.MainPhotoID(), "Main photo ID getter should return correct value")
		assert.Equal(t, *s.validPhotos, plant.GetPhotos(), "Photos getter should return correct value")
		assert.Equal(t, "Test Category", plant.GetCategory(), "Category getter should return correct value")
		assert.Equal(t, s.validSpec, plant.GetSpecification(), "Specification getter should return correct value")
		assert.Equal(t, s.validTime, plant.CreatedAt(), "CreatedAt getter should return correct value")
		assert.Equal(t, s.validTime, plant.UpdatedAt(), "UpdatedAt getter should return correct value")
	})
}

func TestPlant(t *testing.T) {
	suite.RunSuite(t, new(PlantTestSuite))
}
