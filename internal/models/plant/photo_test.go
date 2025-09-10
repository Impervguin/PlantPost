//go:build unit

package plant

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PlantPhotoTestSuite struct {
	suite.Suite
	validID          uuid.UUID
	validFileID      uuid.UUID
	validDescription string
}

func (s *PlantPhotoTestSuite) BeforeEach(t provider.T) {
	t.Epic("Plant Photos")
	t.Feature("Plant Photo Management")

	s.validID = uuid.New()
	s.validFileID = uuid.New()
	s.validDescription = "Test description"
}

func (s *PlantPhotoTestSuite) TestCreatePlantPhotoSuccess(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful creation of plant photo with valid parameters")

	var photo *PlantPhoto
	var err error

	t.WithNewStep("Create plant photo", func(ctx provider.StepCtx) {
		photo, err = CreatePlantPhoto(s.validID, s.validFileID, s.validDescription)
	})

	t.WithNewStep("Verify photo properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, s.validID, photo.id)
		assert.Equal(t, s.validFileID, photo.fileID)
		assert.Equal(t, s.validDescription, photo.description)
	})
}

func (s *PlantPhotoTestSuite) TestCreatePlantPhotoValidationErrors(t provider.T) {
	t.Tags("negative", "validation")
	t.Description("Test validation errors during plant photo creation")

	testCases := []struct {
		name        string
		id          uuid.UUID
		fileID      uuid.UUID
		expectedErr string
	}{
		{
			name:        "Empty ID",
			id:          uuid.Nil,
			fileID:      s.validFileID,
			expectedErr: "plant photo ID cannot be empty",
		},
		{
			name:        "Empty file ID",
			id:          s.validID,
			fileID:      uuid.Nil,
			expectedErr: "plant photo file ID cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t provider.T) {
			t.Tags("negative", "validation")

			t.WithNewStep("Prepare invalid parameters", func(ctx provider.StepCtx) {
				ctx.WithParameters(
					allure.NewParameter("ID", tc.id.String()),
					allure.NewParameter("FileID", tc.fileID.String()),
				)
			})

			var err error
			t.WithNewStep("Attempt to create photo", func(ctx provider.StepCtx) {
				_, err = CreatePlantPhoto(tc.id, tc.fileID, s.validDescription)
			})

			t.WithNewStep("Verify validation error", func(ctx provider.StepCtx) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
			})
		})
	}
}

func (s *PlantPhotoTestSuite) TestNewPlantPhotoSuccess(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful creation of plant photo with auto-generated ID")

	var photo *PlantPhoto
	var err error

	t.WithNewStep("Create plant photo with NewPlantPhoto", func(ctx provider.StepCtx) {
		photo, err = NewPlantPhoto(s.validFileID, s.validDescription)
	})

	t.WithNewStep("Verify auto-generated properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, photo.id)
		assert.Equal(t, s.validFileID, photo.fileID)
		assert.Equal(t, s.validDescription, photo.description)
	})
}

func (s *PlantPhotoTestSuite) TestComparePhotos(t provider.T) {
	t.Tags("comparison", "functionality")
	t.Description("Test photo comparison functionality")

	photo1, _ := CreatePlantPhoto(s.validID, s.validFileID, "Photo 1")
	photo2, _ := CreatePlantPhoto(uuid.New(), uuid.New(), "Photo 2")
	photo3, _ := CreatePlantPhoto(s.validID, uuid.New(), "Photo 3")
	photo4, _ := CreatePlantPhoto(uuid.New(), s.validFileID, "Photo 4")

	t.WithNewStep("Compare different photo scenarios", func(ctx provider.StepCtx) {
		assert.True(t, photo1.Compare(photo1), "Same object should match")
		assert.True(t, photo1.Compare(photo3), "Same ID should match")
		assert.True(t, photo1.Compare(photo4), "Same fileID should match")
		assert.False(t, photo1.Compare(photo2), "Different ID and fileID should not match")
		assert.False(t, photo1.Compare(nil), "Nil comparison should return false")
	})
}

func (s *PlantPhotoTestSuite) TestGetters(t provider.T) {
	t.Tags("getters", "functionality")
	t.Description("Test plant photo getter methods")

	photo, _ := CreatePlantPhoto(s.validID, s.validFileID, s.validDescription)

	t.WithNewStep("Verify getter values", func(ctx provider.StepCtx) {
		assert.Equal(t, s.validID, photo.ID())
		assert.Equal(t, s.validFileID, photo.FileID())
		assert.Equal(t, s.validDescription, photo.Description())
	})
}

func (s *PlantPhotoTestSuite) TestValidate(t provider.T) {
	t.Tags("validation", "functionality")
	t.Description("Test plant photo validation")

	validPhoto := &PlantPhoto{id: s.validID, fileID: s.validFileID}
	invalidPhoto := &PlantPhoto{id: uuid.Nil, fileID: s.validFileID}

	t.WithNewStep("Validate different photo states", func(ctx provider.StepCtx) {
		assert.NoError(t, validPhoto.Validate(), "Valid photo should pass validation")
		assert.Error(t, invalidPhoto.Validate(), "Invalid photo should fail validation")
	})
}

type PlantPhotosTestSuite struct {
	suite.Suite
	photo1 *PlantPhoto
	photo2 *PlantPhoto
}

func (s *PlantPhotosTestSuite) BeforeEach(t provider.T) {
	t.Epic("Plant Photos")
	t.Feature("Plant Photos Collection Management")

	s.photo1, _ = CreatePlantPhoto(uuid.New(), uuid.New(), "Photo 1")
	s.photo2, _ = CreatePlantPhoto(uuid.New(), uuid.New(), "Photo 2")
}

func (s *PlantPhotosTestSuite) TestNewPlantPhotos(t provider.T) {
	t.Tags("creation", "initialization")
	t.Description("Test creation of new plant photos collection")

	pp := NewPlantPhotos()

	t.WithNewStep("Verify empty collection", func(ctx provider.StepCtx) {
		assert.NotNil(t, pp)
		assert.Equal(t, 0, pp.Len())
	})
}

func (s *PlantPhotosTestSuite) TestAddSuccess(t provider.T) {
	t.Tags("positive", "add")
	t.Description("Test successful addition of photos to collection")

	pp := NewPlantPhotos()

	t.WithNewStep("Add first photo", func(ctx provider.StepCtx) {
		require.NoError(t, pp.Add(s.photo1))
		assert.Equal(t, 1, pp.Len())
	})

	t.WithNewStep("Add second photo", func(ctx provider.StepCtx) {
		require.NoError(t, pp.Add(s.photo2))
		assert.Equal(t, 2, pp.Len())
	})
}

func (s *PlantPhotosTestSuite) TestAddDuplicate(t provider.T) {
	t.Tags("negative", "add")
	t.Description("Test duplicate photo addition error")

	pp := NewPlantPhotos()
	require.NoError(t, pp.Add(s.photo1))

	var err error
	t.WithNewStep("Attempt to add duplicate photo", func(ctx provider.StepCtx) {
		err = pp.Add(s.photo1)
	})

	t.WithNewStep("Verify duplicate error", func(ctx provider.StepCtx) {
		require.Error(t, err)
		assert.Equal(t, "photo already exists", err.Error())
		assert.Equal(t, 1, pp.Len())
	})
}

func (s *PlantPhotosTestSuite) TestRemoveSuccess(t provider.T) {
	t.Tags("positive", "remove")
	t.Description("Test successful removal of photos from collection")

	pp := NewPlantPhotos()
	require.NoError(t, pp.Add(s.photo1))
	require.NoError(t, pp.Add(s.photo2))

	t.WithNewStep("Remove first photo", func(ctx provider.StepCtx) {
		require.NoError(t, pp.Remove(s.photo1))
		assert.Equal(t, 1, pp.Len())
	})
}

func (s *PlantPhotosTestSuite) TestRemoveNotFound(t provider.T) {
	t.Tags("negative", "remove")
	t.Description("Test removal of non-existent photo error")

	pp := NewPlantPhotos()
	require.NoError(t, pp.Add(s.photo1))

	var err error
	t.WithNewStep("Attempt to remove non-existent photo", func(ctx provider.StepCtx) {
		err = pp.Remove(s.photo2)
	})

	t.WithNewStep("Verify not found error", func(ctx provider.StepCtx) {
		require.Error(t, err)
		assert.Equal(t, "photo not found", err.Error())
		assert.Equal(t, 1, pp.Len())
	})
}

func (s *PlantPhotosTestSuite) TestIterateSuccess(t provider.T) {
	t.Tags("iteration", "functionality")
	t.Description("Test successful iteration over photos collection")

	pp := NewPlantPhotos()
	require.NoError(t, pp.Add(s.photo1))
	require.NoError(t, pp.Add(s.photo2))

	var collected []PlantPhoto
	var err error
	t.WithNewStep("Iterate over photos", func(ctx provider.StepCtx) {
		err = pp.Iterate(func(e PlantPhoto) error {
			collected = append(collected, e)
			return nil
		})
	})

	t.WithNewStep("Verify iteration results", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Len(t, collected, 2)
		assert.Contains(t, collected, *s.photo1)
		assert.Contains(t, collected, *s.photo2)
	})
}

func (s *PlantPhotosTestSuite) TestIterateWithError(t provider.T) {
	t.Tags("negative", "iteration")
	t.Description("Test iteration with custom error")

	pp := NewPlantPhotos()
	require.NoError(t, pp.Add(s.photo1))
	require.NoError(t, pp.Add(s.photo2))

	testErr := errors.New("test error")
	var err error
	t.WithNewStep("Iterate with error condition", func(ctx provider.StepCtx) {
		err = pp.Iterate(func(e PlantPhoto) error {
			if e.Compare(s.photo2) {
				return testErr
			}
			return nil
		})
	})

	t.WithNewStep("Verify iteration error", func(ctx provider.StepCtx) {
		require.Error(t, err)
		assert.Equal(t, testErr, err)
	})
}

func (s *PlantPhotosTestSuite) TestIterateUpdateSuccess(t provider.T) {
	t.Tags("positive", "update")
	t.Description("Test successful update iteration over photos collection")

	pp := NewPlantPhotos()
	require.NoError(t, pp.Add(s.photo1))
	require.NoError(t, pp.Add(s.photo2))

	newDesc := "Updated description"
	var err error
	t.WithNewStep("Update photos during iteration", func(ctx provider.StepCtx) {
		err = pp.IterateUpdate(func(e *PlantPhoto) error {
			if e.Compare(s.photo2) {
				e.description = newDesc
			}
			return nil
		})
	})

	t.WithNewStep("Verify update results", func(ctx provider.StepCtx) {
		require.NoError(t, err)

		var found bool
		_ = pp.Iterate(func(e PlantPhoto) error {
			if e.Compare(s.photo2) {
				assert.Equal(t, newDesc, e.description)
				found = true
			}
			return nil
		})
		assert.True(t, found)
	})
}

func (s *PlantPhotosTestSuite) TestIterateUpdateValidationError(t provider.T) {
	t.Tags("negative", "update")
	t.Description("Test update iteration with validation error")

	pp := NewPlantPhotos()
	require.NoError(t, pp.Add(s.photo1))

	var err error
	t.WithNewStep("Update with invalid data", func(ctx provider.StepCtx) {
		err = pp.IterateUpdate(func(e *PlantPhoto) error {
			e.fileID = uuid.Nil
			return nil
		})
	})

	t.WithNewStep("Verify validation error", func(ctx provider.StepCtx) {
		require.Error(t, err)
		assert.Contains(t, err.Error(), "plant photo file ID cannot be empty")
	})
}

func (s *PlantPhotosTestSuite) TestLen(t provider.T) {
	t.Tags("functionality", "size")
	t.Description("Test collection length functionality")

	pp := NewPlantPhotos()

	t.WithNewStep("Verify initial length", func(ctx provider.StepCtx) {
		assert.Equal(t, 0, pp.Len())
	})

	t.WithNewStep("Verify length after additions", func(ctx provider.StepCtx) {
		require.NoError(t, pp.Add(s.photo1))
		assert.Equal(t, 1, pp.Len())

		require.NoError(t, pp.Add(s.photo2))
		assert.Equal(t, 2, pp.Len())
	})
}

func (s *PlantPhotosTestSuite) TestValidate(t provider.T) {
	t.Tags("validation", "functionality")
	t.Description("Test collection validation")

	pp := NewPlantPhotos()

	t.WithNewStep("Validate empty collection", func(ctx provider.StepCtx) {
		assert.NoError(t, pp.Validate())
	})

	t.WithNewStep("Validate collection with photos", func(ctx provider.StepCtx) {
		require.NoError(t, pp.Add(s.photo1))
		assert.NoError(t, pp.Validate())
	})
}

func TestPlantPhoto(t *testing.T) {
	suite.RunSuite(t, new(PlantPhotoTestSuite))
}

func TestPlantPhotos(t *testing.T) {
	suite.RunSuite(t, new(PlantPhotosTestSuite))
}
