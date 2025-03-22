package post_test

import (
	"PlantSite/internal/models/post"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPostContent(t *testing.T) {
	testID := 0

	t.Logf("Test %d: plain text content", testID)
	content, err := post.NewContent("Hello, world!", post.ContentTypePlainText)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", content.Text)
	assert.Equal(t, post.ContentTypePlainText, content.ContentType)
	testID++

	t.Logf("Test %d: empty content text", testID)
	content, err = post.NewContent("", post.ContentTypePlainText)
	assert.Error(t, err)
	assert.Nil(t, content)
	testID++

	t.Logf("Test %d: unsupported content format", testID)
	content, err = post.NewContent("Hello, world!", "unsupported_format")
	assert.Error(t, err)
	assert.Nil(t, content)
	testID++
}

func TestPostCreation(t *testing.T) {
	testID := 0

	t.Logf("Test %d: valid post", testID)
	{
		postID := uuid.New()
		authorID := uuid.New()
		postTitle := "Post Title"
		content, err := post.NewContent("Hello, world!", post.ContentTypePlainText)
		assert.NoError(t, err)
		postPhotos := []post.PostPhoto{
			{ID: uuid.New(), PlaceNumber: 1, FileID: uuid.New()},
			{ID: uuid.New(), PlaceNumber: 2, FileID: uuid.New()},
		}
		res, err := post.CreatePost(
			postID,
			postTitle,
			*content,
			[]string{"tag1", "tag2"},
			authorID,
			postPhotos,
			time.Now(),
		)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, authorID, res.AuthorID)
		assert.Equal(t, postID, res.ID)
		assert.Equal(t, postTitle, res.Title)
		assert.Equal(t, 2, len(res.Photos))
		assert.Equal(t, *content, res.Content)

		testID++
	}
	t.Logf("Test %d: no tags", testID)
	{
		postID := uuid.New()
		authorID := uuid.New()
		postTitle := "Post Title"
		content, err := post.NewContent("Hello, world!", post.ContentTypePlainText)
		assert.NoError(t, err)
		postPhotos := []post.PostPhoto{
			{ID: uuid.New(), PlaceNumber: 1, FileID: uuid.New()},
			{ID: uuid.New(), PlaceNumber: 2, FileID: uuid.New()},
		}
		res, err := post.CreatePost(
			postID,
			postTitle,
			*content,
			[]string{},
			authorID,
			postPhotos,
			time.Now(),
		)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, authorID, res.AuthorID)
		assert.Equal(t, postID, res.ID)
		assert.Equal(t, postTitle, res.Title)
		assert.Equal(t, len(res.Tags), 0)
		assert.Equal(t, 2, len(res.Photos))
		assert.Equal(t, *content, res.Content)
		testID++
	}

	t.Logf("Test %d: invalid post (empty ID)", testID)
	{
		postID := uuid.Nil
		authorID := uuid.New()
		postTitle := "Post Title"
		content, err := post.NewContent("Hello, world!", post.ContentTypePlainText)
		assert.NoError(t, err)
		postPhotos := []post.PostPhoto{
			{ID: uuid.New(), PlaceNumber: 1, FileID: uuid.New()},
			{ID: uuid.New(), PlaceNumber: 2, FileID: uuid.New()},
		}
		res, err := post.CreatePost(
			postID,
			postTitle,
			*content,
			[]string{"tag1", "tag2"},
			authorID,
			postPhotos,
			time.Now(),
		)
		assert.Error(t, err)
		assert.Nil(t, res)
		testID++
	}
	t.Logf("Test %d: invalid post (empty AuthorID)", testID)
	{
		postID := uuid.New()
		authorID := uuid.Nil
		postTitle := "Post Title"
		content, err := post.NewContent("Hello, world!", post.ContentTypePlainText)
		assert.NoError(t, err)
		postPhotos := []post.PostPhoto{
			{ID: uuid.New(), PlaceNumber: 1, FileID: uuid.New()},
			{ID: uuid.New(), PlaceNumber: 2, FileID: uuid.New()},
		}
		res, err := post.CreatePost(
			postID,
			postTitle,
			*content,
			[]string{"tag1", "tag2"},
			authorID,
			postPhotos,
			time.Now(),
		)
		assert.Error(t, err)
		assert.Nil(t, res)
		testID++
	}
	t.Logf("Test %d: invalid post (empty Title)", testID)
	{
		postID := uuid.New()
		authorID := uuid.New()
		postTitle := ""
		content, err := post.NewContent("Hello, world!", post.ContentTypePlainText)
		assert.NoError(t, err)
		postPhotos := []post.PostPhoto{
			{ID: uuid.New(), PlaceNumber: 1, FileID: uuid.New()},
			{ID: uuid.New(), PlaceNumber: 2, FileID: uuid.New()},
		}
		res, err := post.CreatePost(
			postID,
			postTitle,
			*content,
			[]string{"tag1", "tag2"},
			authorID,
			postPhotos,
			time.Now(),
		)
		assert.Error(t, err)
		assert.Nil(t, res)
		testID++
	}
	t.Logf("Test %d: invalid post (exceede photo count)", testID)
	{
		postID := uuid.New()
		authorID := uuid.New()
		postTitle := "Post Title"
		content, err := post.NewContent("Hello, world!", post.ContentTypePlainText)
		assert.NoError(t, err)
		postPhotos := []post.PostPhoto{}
		for i := 0; i < post.MaximumPhotoPerPostCount+1; i++ {
			postPhotos = append(postPhotos, post.PostPhoto{ID: uuid.New(), PlaceNumber: i, FileID: uuid.New()})
		}
		res, err := post.CreatePost(
			postID,
			postTitle,
			*content,
			[]string{"tag1", "tag2"},
			authorID,
			postPhotos,
			time.Now(),
		)
		assert.Error(t, err)
		assert.Nil(t, res)
		testID++
	}
	t.Logf("Test %d: invalid post (photos nil)", testID)
	{
		postID := uuid.New()
		authorID := uuid.New()
		postTitle := "Post Title"
		content, err := post.NewContent("Hello, world!", post.ContentTypePlainText)
		assert.NoError(t, err)
		var postPhotos []post.PostPhoto = nil
		res, err := post.CreatePost(
			postID,
			postTitle,
			*content,
			[]string{"tag1", "tag2"},
			authorID,
			postPhotos,
			time.Now(),
		)
		assert.Error(t, err)
		assert.Nil(t, res)
		testID++
	}
	t.Logf("Test %d: invalid post (tags nil)", testID)
	{
		postID := uuid.New()
		authorID := uuid.New()
		postTitle := "Post Title"
		content, err := post.NewContent("Hello, world!", post.ContentTypePlainText)
		assert.NoError(t, err)
		postPhotos := []post.PostPhoto{
			{ID: uuid.New(), PlaceNumber: 1, FileID: uuid.New()},
			{ID: uuid.New(), PlaceNumber: 2, FileID: uuid.New()},
		}
		res, err := post.CreatePost(
			postID,
			postTitle,
			*content,
			nil,
			authorID,
			postPhotos,
			time.Now(),
		)
		assert.Error(t, err)
		assert.Nil(t, res)
		testID++
	}
	t.Logf("Test %d: invalid post (duplicate photo placenumber)", testID)
	{
		postID := uuid.New()
		authorID := uuid.New()
		postTitle := "Post Title"
		content, err := post.NewContent("Hello, world!", post.ContentTypePlainText)
		assert.NoError(t, err)
		postPhotos := []post.PostPhoto{
			{ID: uuid.New(), PlaceNumber: 1, FileID: uuid.New()},
			{ID: uuid.New(), PlaceNumber: 1, FileID: uuid.New()},
		}
		res, err := post.CreatePost(
			postID,
			postTitle,
			*content,
			[]string{"tag1", "tag2"},
			authorID,
			postPhotos,
			time.Now(),
		)
		assert.Error(t, err)
		assert.Nil(t, res)
	}
	t.Logf("Test %d: invalid post (created in future)", testID)
	{
		postID := uuid.New()
		authorID := uuid.New()
		postTitle := "Post Title"
		content, err := post.NewContent("Hello, world!", post.ContentTypePlainText)
		assert.NoError(t, err)
		postPhotos := []post.PostPhoto{
			{ID: uuid.New(), PlaceNumber: 1, FileID: uuid.New()},
			{ID: uuid.New(), PlaceNumber: 2, FileID: uuid.New()},
		}
		res, err := post.CreatePost(
			postID,
			postTitle,
			*content,
			[]string{"tag1", "tag2"},
			authorID,
			postPhotos,
			time.Now().Add(time.Hour),
		)
		assert.Error(t, err)
		assert.Nil(t, res)
	}
}
