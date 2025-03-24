package post

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPost(t *testing.T) {
	// Подготовка тестовых данных
	validID := uuid.New()
	validTitle := "Test Post"
	validContent, _ := NewContent("Test content", ContentTypePlainText)
	validTags := []string{"tag1", "tag2"}
	validAuthorID := uuid.New()
	validPhotos := PostPhotos{}
	validTime := time.Now().Add(-time.Hour)

	t.Run("CreatePost - успешное создание", func(t *testing.T) {
		post, err := CreatePost(
			validID,
			validTitle,
			*validContent,
			validTags,
			validAuthorID,
			validPhotos,
			validTime,
			validTime,
		)

		require.NoError(t, err)
		assert.Equal(t, validID, post.ID())
		assert.Equal(t, validTitle, post.Title())
		assert.Equal(t, *validContent, post.Content())
		assert.Equal(t, validTags, post.Tags())
		assert.Equal(t, validAuthorID, post.AuthorID())
		assert.Equal(t, validPhotos, post.Photos())
		assert.Equal(t, validTime, post.CreatedAt())
		assert.Equal(t, validTime, post.UpdatedAt())
	})

	t.Run("CreatePost - ошибки валидации", func(t *testing.T) {
		testCases := []struct {
			name        string
			id          uuid.UUID
			title       string
			content     Content
			tags        []string
			authorID    uuid.UUID
			photos      PostPhotos
			createdAt   time.Time
			updatedAt   time.Time
			expectError bool
		}{
			{
				"Пустой ID",
				uuid.Nil, validTitle, *validContent, validTags, validAuthorID, validPhotos, validTime, validTime, true,
			},
			{
				"Пустой заголовок",
				validID, "", *validContent, validTags, validAuthorID, validPhotos, validTime, validTime, true,
			},
			{
				"Пустой authorID",
				validID, validTitle, *validContent, validTags, uuid.Nil, validPhotos, validTime, validTime, true,
			},
			{
				"Nil теги",
				validID, validTitle, *validContent, nil, validAuthorID, validPhotos, validTime, validTime, true,
			},
			{
				"Пустой тег",
				validID, validTitle, *validContent, []string{""}, validAuthorID, validPhotos, validTime, validTime, true,
			},
			{
				"Обновление до создания",
				validID, validTitle, *validContent, validTags, validAuthorID, validPhotos, validTime, validTime.Add(-time.Hour), true,
			},
			{
				"Обновление в будущем",
				validID, validTitle, *validContent, validTags, validAuthorID, validPhotos, validTime, time.Now().Add(time.Hour), true,
			},
			{
				"Корректные данные",
				validID, validTitle, *validContent, validTags, validAuthorID, validPhotos, validTime, validTime, false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := CreatePost(
					tc.id,
					tc.title,
					tc.content,
					tc.tags,
					tc.authorID,
					tc.photos,
					tc.createdAt,
					tc.updatedAt,
				)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("NewPost - успешное создание", func(t *testing.T) {
		post, err := NewPost(
			validTitle,
			*validContent,
			validTags,
			validAuthorID,
			&validPhotos,
		)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, post.ID())
		assert.True(t, post.CreatedAt().Before(time.Now()) || post.CreatedAt().Equal(time.Now()))
		assert.True(t, post.UpdatedAt().Before(time.Now()) || post.UpdatedAt().Equal(time.Now()))
	})

	t.Run("UpdateTitle", func(t *testing.T) {
		post := &Post{title: "Old Title", updatedAt: validTime}
		newTitle := "New Title"

		err := post.UpdateTitle(newTitle)
		require.NoError(t, err)
		assert.Equal(t, newTitle, post.title)
		assert.True(t, post.updatedAt.After(validTime))
	})

	t.Run("UpdateContent", func(t *testing.T) {
		post := &Post{content: *validContent, updatedAt: validTime}
		newContent, _ := NewContent("New content", ContentTypePlainText)

		err := post.UpdateContent(*newContent)
		require.NoError(t, err)
		assert.Equal(t, *newContent, post.content)
		assert.True(t, post.updatedAt.After(validTime))
	})

	t.Run("UpdateTags", func(t *testing.T) {
		post := &Post{tags: validTags, updatedAt: validTime}
		newTags := []string{"new1", "new2"}

		err := post.UpdateTags(newTags)
		require.NoError(t, err)
		assert.Equal(t, newTags, post.tags)
		assert.True(t, post.updatedAt.After(validTime))
	})

	t.Run("AddTag", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			post := &Post{tags: validTags}
			newTag := "newTag"

			err := post.AddTag(newTag)
			require.NoError(t, err)
			assert.Contains(t, post.tags, newTag)
		})

		t.Run("Duplicate", func(t *testing.T) {
			post := &Post{tags: validTags}

			err := post.AddTag(validTags[0])
			assert.Error(t, err)
		})

		t.Run("Max tags", func(t *testing.T) {
			tags := make([]string, MaximumTagCount)
			for i := range tags {
				tags[i] = fmt.Sprintf("tag%d", i)
			}
			post := &Post{tags: tags}

			err := post.AddTag("new")
			assert.Error(t, err)
		})
	})

	t.Run("RemoveTag", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			post := &Post{tags: validTags}

			err := post.RemoveTag(validTags[0])
			require.NoError(t, err)
			assert.NotContains(t, post.tags, validTags[0])
		})

		t.Run("Not found", func(t *testing.T) {
			post := &Post{tags: validTags}

			err := post.RemoveTag("nonexistent")
			assert.Error(t, err)
		})
	})

	t.Run("AddPhoto", func(t *testing.T) {
		post := &Post{photos: PostPhotos{}}
		photo := &PostPhoto{}

		err := post.AddPhoto(photo)
		assert.Error(t, err)
	})

	t.Run("Getters", func(t *testing.T) {
		post := &Post{
			id:        validID,
			title:     validTitle,
			content:   *validContent,
			tags:      validTags,
			authorID:  validAuthorID,
			photos:    validPhotos,
			createdAt: validTime,
			updatedAt: validTime,
		}

		assert.Equal(t, validID, post.ID())
		assert.Equal(t, validTitle, post.Title())
		assert.Equal(t, *validContent, post.Content())
		assert.Equal(t, validTags, post.Tags())
		assert.Equal(t, validAuthorID, post.AuthorID())
		assert.Equal(t, validPhotos, post.Photos())
		assert.Equal(t, validTime, post.CreatedAt())
		assert.Equal(t, validTime, post.UpdatedAt())
	})
}
