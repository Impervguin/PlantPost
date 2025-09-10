//go:build unit

package album

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AlbumTestSuite struct {
	suite.Suite
	validID          uuid.UUID
	validName        string
	validDescription string
	validPlantIDs    uuid.UUIDs
	validOwnerID     uuid.UUID
	validCreatedAt   time.Time
	validUpdatedAt   time.Time
}

func (s *AlbumTestSuite) BeforeEach(t provider.T) {
	t.Epic("Albums")
	t.Feature("Album Management")

	s.validID = uuid.New()
	s.validName = "Test Album"
	s.validDescription = "Test Description"
	s.validPlantIDs = uuid.UUIDs{uuid.New(), uuid.New()}
	s.validOwnerID = uuid.New()
	s.validCreatedAt = time.Now().Add(-time.Hour)
	s.validUpdatedAt = time.Now()
}

func (s *AlbumTestSuite) TestCreateAlbumSuccess(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful creation of album with valid parameters")

	var album *Album
	var err error

	t.WithNewStep("Create album with valid data", func(ctx provider.StepCtx) {
		album, err = CreateAlbum(
			s.validID,
			s.validName,
			s.validDescription,
			s.validPlantIDs,
			s.validOwnerID,
			s.validCreatedAt,
			s.validUpdatedAt,
		)
	})

	t.WithNewStep("Verify album properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, s.validID, album.id)
		assert.Equal(t, s.validName, album.name)
		assert.Equal(t, s.validDescription, album.description)
		assert.Equal(t, s.validPlantIDs, album.plantIDs)
		assert.Equal(t, s.validOwnerID, album.ownerID)
		assert.Equal(t, s.validCreatedAt, album.createdAt)
	})
}

func (s *AlbumTestSuite) TestCreateAlbumValidationErrors(t provider.T) {
	t.Tags("negative", "validation")
	t.Description("Test validation errors during album creation")

	testCases := []struct {
		name        string
		id          uuid.UUID
		nameStr     string
		description string
		plantIDs    uuid.UUIDs
		ownerID     uuid.UUID
		createdAt   time.Time
		updatedAt   time.Time
		expectError bool
	}{
		{
			"Empty ID",
			uuid.Nil,
			s.validName,
			s.validDescription,
			s.validPlantIDs,
			s.validOwnerID,
			s.validCreatedAt,
			s.validUpdatedAt,
			true,
		},
		{
			"Empty name",
			s.validID,
			"",
			s.validDescription,
			s.validPlantIDs,
			s.validOwnerID,
			s.validCreatedAt,
			s.validUpdatedAt,
			true,
		},
		{
			"Empty ownerID",
			s.validID,
			s.validName,
			s.validDescription,
			s.validPlantIDs,
			uuid.Nil,
			s.validCreatedAt,
			s.validUpdatedAt,
			true,
		},
		{
			"Nil plantIDs",
			s.validID,
			s.validName,
			s.validDescription,
			nil,
			s.validOwnerID,
			s.validCreatedAt,
			s.validUpdatedAt,
			true,
		},
		{
			"Empty plantID in list",
			s.validID,
			s.validName,
			s.validDescription,
			uuid.UUIDs{uuid.New(), uuid.Nil},
			s.validOwnerID,
			s.validCreatedAt,
			s.validUpdatedAt,
			true,
		},
		{
			"Future update time",
			s.validID,
			s.validName,
			s.validDescription,
			s.validPlantIDs,
			s.validOwnerID,
			s.validCreatedAt,
			time.Now().Add(time.Hour),
			true,
		},
		{
			"Creation after update",
			s.validID,
			s.validName,
			s.validDescription,
			s.validPlantIDs,
			s.validOwnerID,
			s.validUpdatedAt,
			s.validCreatedAt,
			true,
		},
		{
			"Valid data",
			s.validID,
			s.validName,
			s.validDescription,
			s.validPlantIDs,
			s.validOwnerID,
			s.validCreatedAt,
			s.validUpdatedAt,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t provider.T) {
			t.Tags("negative", "validation")

			t.WithNewStep("Prepare test parameters", func(ctx provider.StepCtx) {
				ctx.WithParameters(
					allure.NewParameter("Name", tc.nameStr),
					allure.NewParameter("OwnerID", tc.ownerID.String()),
					allure.NewParameter("ExpectError", tc.expectError),
				)
			})

			var err error
			t.WithNewStep("Attempt to create album", func(ctx provider.StepCtx) {
				_, err = CreateAlbum(
					tc.id,
					tc.nameStr,
					tc.description,
					tc.plantIDs,
					tc.ownerID,
					tc.createdAt,
					tc.updatedAt,
				)
			})

			t.WithNewStep("Verify validation result", func(ctx provider.StepCtx) {
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		})
	}
}

func (s *AlbumTestSuite) TestNewAlbumSuccess(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful creation of album with auto-generated fields")

	var album *Album
	var err error

	t.WithNewStep("Create album with NewAlbum", func(ctx provider.StepCtx) {
		album, err = NewAlbum(
			s.validName,
			s.validDescription,
			s.validPlantIDs,
			s.validOwnerID,
		)
	})

	t.WithNewStep("Verify auto-generated properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, album.id)
		assert.True(t, album.createdAt.Before(time.Now()) || album.createdAt.Equal(time.Now()))
		assert.True(t, album.updatedAt.Before(time.Now()) || album.updatedAt.Equal(time.Now()))
		assert.Equal(t, s.validName, album.name)
		assert.Equal(t, s.validDescription, album.description)
		assert.Equal(t, s.validPlantIDs, album.plantIDs)
		assert.Equal(t, s.validOwnerID, album.ownerID)
	})
}

func (s *AlbumTestSuite) TestGetOwnerID(t provider.T) {
	t.Tags("functionality", "access")
	t.Description("Test owner ID retrieval")

	album := &Album{ownerID: s.validOwnerID}

	t.WithNewStep("Get owner ID", func(ctx provider.StepCtx) {
		assert.Equal(t, s.validOwnerID, album.GetOwnerID())
	})
}

func (s *AlbumTestSuite) TestUpdateName(t provider.T) {
	t.Tags("positive", "update")
	t.Description("Successful update of album name")

	album := &Album{name: "Old Name", updatedAt: s.validUpdatedAt}
	newName := "New Name"

	var err error
	t.WithNewStep("Update album name", func(ctx provider.StepCtx) {
		err = album.UpdateName(newName)
	})

	t.WithNewStep("Verify name update", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, newName, album.name)
		assert.True(t, album.updatedAt.After(s.validUpdatedAt))
	})
}

func (s *AlbumTestSuite) TestUpdateDescription(t provider.T) {
	t.Tags("positive", "update")
	t.Description("Successful update of album description")

	album := &Album{description: "Old Desc", updatedAt: s.validUpdatedAt}
	newDesc := "New Desc"

	var err error
	t.WithNewStep("Update album description", func(ctx provider.StepCtx) {
		err = album.UpdateDescription(newDesc)
	})

	t.WithNewStep("Verify description update", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, newDesc, album.description)
		assert.True(t, album.updatedAt.After(s.validUpdatedAt))
	})
}

func (s *AlbumTestSuite) TestAddPlant(t provider.T) {
	t.Tags("positive", "update")
	t.Description("Successful addition of plant to album")

	album := &Album{plantIDs: s.validPlantIDs, updatedAt: s.validUpdatedAt}
	newPlantID := uuid.New()

	var err error
	t.WithNewStep("Add plant to album", func(ctx provider.StepCtx) {
		err = album.AddPlant(newPlantID)
	})

	t.WithNewStep("Verify plant addition", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Contains(t, album.plantIDs, newPlantID)
		assert.True(t, album.updatedAt.After(s.validUpdatedAt))
	})
}

func (s *AlbumTestSuite) TestRemovePlantSuccess(t provider.T) {
	t.Tags("positive", "update")
	t.Description("Successful removal of plant from album")

	plantToRemove := s.validPlantIDs[0]
	album := &Album{plantIDs: s.validPlantIDs, updatedAt: s.validUpdatedAt}

	var err error
	t.WithNewStep("Remove plant from album", func(ctx provider.StepCtx) {
		err = album.RemovePlant(plantToRemove)
	})

	t.WithNewStep("Verify plant removal", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.NotContains(t, album.plantIDs, plantToRemove)
		assert.True(t, album.updatedAt.After(s.validUpdatedAt))
	})
}

func (s *AlbumTestSuite) TestRemovePlantNotFound(t provider.T) {
	t.Tags("negative", "update")
	t.Description("Failed removal of non-existent plant from album")

	album := &Album{plantIDs: s.validPlantIDs, updatedAt: s.validUpdatedAt}
	nonExistentPlant := uuid.New()

	var err error
	t.WithNewStep("Attempt to remove non-existent plant", func(ctx provider.StepCtx) {
		err = album.RemovePlant(nonExistentPlant)
	})

	t.WithNewStep("Verify removal error", func(ctx provider.StepCtx) {
		assert.Error(t, err)
		assert.Equal(t, s.validPlantIDs, album.plantIDs)
	})
}

func TestAlbum(t *testing.T) {
	suite.RunSuite(t, new(AlbumTestSuite))
}
