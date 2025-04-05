package search

import (
	"PlantSite/internal/models/post"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPost создает тестовый пост с заданными параметрами
func mockPost(title string, tags []string, authorID uuid.UUID) (*post.Post, error) {
	content, _ := post.NewContent("Test content", post.ContentTypePlainText)
	photos := post.NewPostPhotos()

	return post.NewPost(
		title,
		*content,
		tags,
		authorID,
		photos,
	)
}

func TestPostFilters(t *testing.T) {
	// Создаем тестовые данные
	authorID1 := uuid.New()
	authorID2 := uuid.New()
	testPost1, err := mockPost("First Post", []string{"tech", "golang"}, authorID1)
	require.NoError(t, err)
	testPost2, err := mockPost("Second Post About Programming", []string{"programming", "java"}, authorID2)
	require.NoError(t, err)
	testPost3, err := mockPost("Third Post", []string{"tech", "python"}, authorID1)
	require.NoError(t, err)

	t.Run("PostTitleFilter", func(t *testing.T) {
		filter := NewPostTitleFilter("First Post")
		assert.True(t, filter.Filter(testPost1))
		assert.False(t, filter.Filter(testPost2))
		assert.False(t, filter.Filter(testPost3))
	})

	t.Run("PostTitleContainsFilter", func(t *testing.T) {
		tests := []struct {
			name  string
			part  string
			post1 bool
			post2 bool
			post3 bool
		}{
			{"Exact match", "First Post", true, false, false},
			{"Partial match", "Post", true, true, true},
			{"Case insensitive", "first", true, false, false},
			{"Part of word", "Prog", false, true, false},
			{"No match", "Nonexistent", false, false, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				filter := NewPostTitleContainsFilter(tt.part)
				assert.Equal(t, tt.post1, filter.Filter(testPost1))
				assert.Equal(t, tt.post2, filter.Filter(testPost2))
				assert.Equal(t, tt.post3, filter.Filter(testPost3))
			})
		}
	})

	t.Run("PostTagFilter", func(t *testing.T) {
		tests := []struct {
			name  string
			tag   string
			post1 bool
			post2 bool
			post3 bool
		}{
			{"Tech tag", "tech", true, false, true},
			{"Programming tag", "programming", false, true, false},
			{"Java tag", "java", false, true, false},
			{"Nonexistent tag", "ruby", false, false, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				filter := NewPostTagFilter([]string{tt.tag})
				assert.Equal(t, tt.post1, filter.Filter(testPost1))
				assert.Equal(t, tt.post2, filter.Filter(testPost2))
				assert.Equal(t, tt.post3, filter.Filter(testPost3))
			})
		}
	})

	t.Run("PostAuthorFilter", func(t *testing.T) {
		filter := NewPostAuthorFilter(authorID1)
		assert.True(t, filter.Filter(testPost1))
		assert.False(t, filter.Filter(testPost2))
		assert.True(t, filter.Filter(testPost3))

		filter = NewPostAuthorFilter(authorID2)
		assert.False(t, filter.Filter(testPost1))
		assert.True(t, filter.Filter(testPost2))
		assert.False(t, filter.Filter(testPost3))
	})
}
