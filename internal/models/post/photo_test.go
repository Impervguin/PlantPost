package post

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostPhoto(t *testing.T) {
	t.Run("CreatePostPhoto", func(t *testing.T) {
		id := uuid.New()
		fileID := uuid.New()
		photo, err := CreatePostPhoto(id, fileID, 1)

		require.NoError(t, err)
		assert.Equal(t, id, photo.ID())
		assert.Equal(t, fileID, photo.FileID())
		assert.Equal(t, 1, photo.PlaceNumber())
	})

	t.Run("NewPostPhoto", func(t *testing.T) {
		fileID := uuid.New()
		photo, err := NewPostPhoto(fileID, 2)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, photo.ID())
		assert.Equal(t, fileID, photo.FileID())
		assert.Equal(t, 2, photo.PlaceNumber())
	})
}

func TestPostPhotos(t *testing.T) {
	photo1, _ := NewPostPhoto(uuid.New(), 1)
	photo2, _ := NewPostPhoto(uuid.New(), 2)
	photo3, _ := NewPostPhoto(uuid.New(), 3)

	t.Run("NewPostPhotos", func(t *testing.T) {
		pp := NewPostPhotos()
		assert.Empty(t, pp.List())
	})

	t.Run("Validate", func(t *testing.T) {
		tests := []struct {
			name    string
			photos  []PostPhoto
			wantErr bool
		}{
			{"Valid", []PostPhoto{*photo1, *photo2}, false},
			{"Duplicate place", []PostPhoto{{placeNumber: 1}, {placeNumber: 1}}, true},
			{"Invalid place number", []PostPhoto{{placeNumber: MaximumPhotoPerPostCount + 1}}, true},
			{"Empty file ID", []PostPhoto{{fileID: uuid.Nil}}, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				pp := &PostPhotos{photos: tt.photos}
				err := pp.Validate()
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Add", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			pp := NewPostPhotos()
			assert.NoError(t, pp.Add(photo1))
			assert.Len(t, pp.List(), 1)
			assert.Equal(t, 1, pp.List()[0].PlaceNumber())
		})

		t.Run("Duplicate", func(t *testing.T) {
			pp := &PostPhotos{photos: []PostPhoto{*photo1}}
			err := pp.Add(photo1)
			assert.Error(t, err)
			assert.Len(t, pp.List(), 1)
		})

		t.Run("Max photos", func(t *testing.T) {
			pp := NewPostPhotos()
			for i := 0; i < MaximumPhotoPerPostCount; i++ {
				p, _ := NewPostPhoto(uuid.New(), i+1)
				assert.NoError(t, pp.Add(p))
			}
			newPhoto, _ := NewPostPhoto(uuid.New(), MaximumPhotoPerPostCount+1)
			err := pp.Add(newPhoto)
			assert.Error(t, err)
		})
	})

	t.Run("Remove", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			pp := &PostPhotos{photos: []PostPhoto{*photo1, *photo2, *photo3}}
			assert.NoError(t, pp.Remove(photo2.ID()))
			assert.Len(t, pp.List(), 2)
			assert.Equal(t, []int{1, 2}, []int{pp.List()[0].PlaceNumber(), pp.List()[1].PlaceNumber()})
		})

		t.Run("Not found", func(t *testing.T) {
			pp := &PostPhotos{photos: []PostPhoto{*photo1}}
			err := pp.Remove(uuid.New())
			assert.Error(t, err)
			assert.Len(t, pp.List(), 1)
		})
	})

	t.Run("RebalancePositions", func(t *testing.T) {
		pp := &PostPhotos{photos: []PostPhoto{
			{placeNumber: 3},
			{placeNumber: 1},
			{placeNumber: 5},
		}}
		pp.RebalancePositions()
		assert.Equal(t, []int{1, 2, 3}, []int{
			pp.photos[0].PlaceNumber(),
			pp.photos[1].PlaceNumber(),
			pp.photos[2].PlaceNumber(),
		})
	})

	t.Run("List", func(t *testing.T) {
		pp := &PostPhotos{photos: []PostPhoto{*photo1, *photo2}}
		list := pp.List()
		assert.Len(t, list, 2)
		assert.NotSame(t, &pp.photos, &list)
	})
}
