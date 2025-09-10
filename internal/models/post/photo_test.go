//go:build unit

package post

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PostPhotoTestSuite struct {
	suite.Suite
	photo1 *PostPhoto
	photo2 *PostPhoto
	photo3 *PostPhoto
}

func (s *PostPhotoTestSuite) BeforeEach(t provider.T) {
	t.Epic("Posts")
	t.Feature("Post Photos Management")

	s.photo1, _ = NewPostPhoto(uuid.New(), 1)
	s.photo2, _ = NewPostPhoto(uuid.New(), 2)
	s.photo3, _ = NewPostPhoto(uuid.New(), 3)
}

func (s *PostPhotoTestSuite) TestCreatePostPhoto(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful creation of post photo with specific ID")

	id := uuid.New()
	fileID := uuid.New()

	var photo *PostPhoto
	var err error

	t.WithNewStep("Create post photo with CreatePostPhoto", func(ctx provider.StepCtx) {
		photo, err = CreatePostPhoto(id, fileID, 1)
	})

	t.WithNewStep("Verify photo properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, id, photo.ID())
		assert.Equal(t, fileID, photo.FileID())
		assert.Equal(t, 1, photo.PlaceNumber())
	})
}

func (s *PostPhotoTestSuite) TestNewPostPhoto(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful creation of post photo with auto-generated ID")

	fileID := uuid.New()

	var photo *PostPhoto
	var err error

	t.WithNewStep("Create post photo with NewPostPhoto", func(ctx provider.StepCtx) {
		photo, err = NewPostPhoto(fileID, 2)
	})

	t.WithNewStep("Verify auto-generated properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, photo.ID())
		assert.Equal(t, fileID, photo.FileID())
		assert.Equal(t, 2, photo.PlaceNumber())
	})
}

func (s *PostPhotoTestSuite) TestNewPostPhotos(t provider.T) {
	t.Tags("creation", "initialization")
	t.Description("Test creation of new post photos collection")

	pp := NewPostPhotos()

	t.WithNewStep("Verify empty collection", func(ctx provider.StepCtx) {
		assert.Empty(t, pp.List())
	})
}

func (s *PostPhotoTestSuite) TestValidate(t provider.T) {
	t.Tags("validation", "functionality")
	t.Description("Test validation of post photos collection")

	tests := []struct {
		name    string
		photos  []PostPhoto
		wantErr bool
	}{
		{"Valid photos", []PostPhoto{*s.photo1, *s.photo2}, false},
		{"Duplicate place numbers", []PostPhoto{{placeNumber: 1}, {placeNumber: 1}}, true},
		{"Invalid place number", []PostPhoto{{placeNumber: MaximumPhotoPerPostCount + 1}}, true},
		{"Empty file ID", []PostPhoto{{fileID: uuid.Nil}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t provider.T) {
			t.Tags("validation")

			pp := &PostPhotos{photos: tt.photos}
			var err error

			t.WithNewStep("Validate photos collection", func(ctx provider.StepCtx) {
				err = pp.Validate()
			})

			t.WithNewStep("Verify validation result", func(ctx provider.StepCtx) {
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		})
	}
}

func (s *PostPhotoTestSuite) TestAddSuccess(t provider.T) {
	t.Tags("positive", "add")
	t.Description("Successful addition of photo to collection")

	pp := NewPostPhotos()

	t.WithNewStep("Add first photo", func(ctx provider.StepCtx) {
		assert.NoError(t, pp.Add(s.photo1))
	})

	t.WithNewStep("Verify photo addition", func(ctx provider.StepCtx) {
		assert.Len(t, pp.List(), 1)
		assert.Equal(t, 1, pp.List()[0].PlaceNumber())
	})
}

func (s *PostPhotoTestSuite) TestAddDuplicate(t provider.T) {
	t.Tags("negative", "add")
	t.Description("Test duplicate photo addition error")

	pp := &PostPhotos{photos: []PostPhoto{*s.photo1}}

	var err error
	t.WithNewStep("Attempt to add duplicate photo", func(ctx provider.StepCtx) {
		err = pp.Add(s.photo1)
	})

	t.WithNewStep("Verify duplicate error", func(ctx provider.StepCtx) {
		assert.Error(t, err)
		assert.Len(t, pp.List(), 1)
	})
}

func (s *PostPhotoTestSuite) TestAddMaxPhotos(t provider.T) {
	t.Tags("negative", "add")
	t.Description("Test maximum photos limit enforcement")

	pp := NewPostPhotos()

	t.WithNewStep("Add maximum allowed photos", func(ctx provider.StepCtx) {
		for i := 0; i < MaximumPhotoPerPostCount; i++ {
			p, _ := NewPostPhoto(uuid.New(), i+1)
			assert.NoError(t, pp.Add(p))
		}
	})

	newPhoto, _ := NewPostPhoto(uuid.New(), MaximumPhotoPerPostCount+1)
	var err error
	t.WithNewStep("Attempt to exceed maximum photos", func(ctx provider.StepCtx) {
		err = pp.Add(newPhoto)
	})

	t.WithNewStep("Verify max photos error", func(ctx provider.StepCtx) {
		assert.Error(t, err)
	})
}

func (s *PostPhotoTestSuite) TestRemoveSuccess(t provider.T) {
	t.Tags("positive", "remove")
	t.Description("Successful removal of photo from collection")

	pp := &PostPhotos{photos: []PostPhoto{*s.photo1, *s.photo2, *s.photo3}}

	t.WithNewStep("Remove middle photo", func(ctx provider.StepCtx) {
		assert.NoError(t, pp.Remove(s.photo2.ID()))
	})

	t.WithNewStep("Verify photo removal and rebalancing", func(ctx provider.StepCtx) {
		assert.Len(t, pp.List(), 2)
		assert.Equal(t, []int{1, 2}, []int{pp.List()[0].PlaceNumber(), pp.List()[1].PlaceNumber()})
	})
}

func (s *PostPhotoTestSuite) TestRemoveNotFound(t provider.T) {
	t.Tags("negative", "remove")
	t.Description("Test removal of non-existent photo error")

	pp := &PostPhotos{photos: []PostPhoto{*s.photo1}}

	var err error
	t.WithNewStep("Attempt to remove non-existent photo", func(ctx provider.StepCtx) {
		err = pp.Remove(uuid.New())
	})

	t.WithNewStep("Verify not found error", func(ctx provider.StepCtx) {
		assert.Error(t, err)
		assert.Len(t, pp.List(), 1)
	})
}

func (s *PostPhotoTestSuite) TestRebalancePositions(t provider.T) {
	t.Tags("functionality", "rebalancing")
	t.Description("Test position rebalancing in photos collection")

	pp := &PostPhotos{photos: []PostPhoto{
		{placeNumber: 3},
		{placeNumber: 1},
		{placeNumber: 5},
	}}

	t.WithNewStep("Rebalance photo positions", func(ctx provider.StepCtx) {
		pp.RebalancePositions()
	})

	t.WithNewStep("Verify positions are rebalanced", func(ctx provider.StepCtx) {
		assert.Equal(t, []int{1, 2, 3}, []int{
			pp.photos[0].PlaceNumber(),
			pp.photos[1].PlaceNumber(),
			pp.photos[2].PlaceNumber(),
		})
	})
}

func (s *PostPhotoTestSuite) TestList(t provider.T) {
	t.Tags("functionality", "access")
	t.Description("Test retrieval of photos list")

	pp := &PostPhotos{photos: []PostPhoto{*s.photo1, *s.photo2}}

	var list []PostPhoto
	t.WithNewStep("Get photos list", func(ctx provider.StepCtx) {
		list = pp.List()
	})

	t.WithNewStep("Verify list properties", func(ctx provider.StepCtx) {
		assert.Len(t, list, 2)
		assert.NotSame(t, &pp.photos, &list)
	})
}

func TestPostPhoto(t *testing.T) {
	suite.RunSuite(t, new(PostPhotoTestSuite))
}
