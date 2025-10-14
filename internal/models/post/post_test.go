//go:build unit

package post

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

type PostTestSuite struct {
	suite.Suite
	validID       uuid.UUID
	validTitle    string
	validContent  *Content
	validTags     []string
	validAuthorID uuid.UUID
	validPhotos   PostPhotos
	validTime     time.Time
}

func (s *PostTestSuite) BeforeEach(t provider.T) {
	t.Epic("Posts")
	t.Feature("Post Management")

	s.validID = uuid.New()
	s.validTitle = "Test Post"
	s.validContent, _ = NewContent("Test content", ContentTypePlainText)
	s.validTags = []string{"tag1", "tag2"}
	s.validAuthorID = uuid.New()
	s.validPhotos = PostPhotos{}
	s.validTime = time.Now().Add(-time.Hour)
}

func (s *PostTestSuite) TestCreatePostSuccess(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful creation of post with valid parameters")

	var post *Post
	var err error

	t.WithNewStep("Create post with valid data", func(ctx provider.StepCtx) {
		post, err = CreatePost(
			s.validID,
			s.validTitle,
			*s.validContent,
			s.validTags,
			s.validAuthorID,
			s.validPhotos,
			s.validTime,
			s.validTime,
		)
	})

	t.WithNewStep("Verify post properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, s.validID, post.ID())
		assert.Equal(t, s.validTitle, post.Title())
		assert.Equal(t, *s.validContent, post.Content())
		assert.Equal(t, s.validTags, post.Tags())
		assert.Equal(t, s.validAuthorID, post.AuthorID())
		assert.Equal(t, s.validPhotos, post.Photos())
		assert.Equal(t, s.validTime, post.CreatedAt())
		assert.Equal(t, s.validTime, post.UpdatedAt())
	})
}

func (s *PostTestSuite) TestCreatePostValidationErrors(t provider.T) {
	t.Tags("negative", "validation")
	t.Description("Test validation errors during post creation")

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
			"Empty ID",
			uuid.Nil, s.validTitle, *s.validContent, s.validTags, s.validAuthorID, s.validPhotos, s.validTime, s.validTime, true,
		},
		{
			"Empty title",
			s.validID, "", *s.validContent, s.validTags, s.validAuthorID, s.validPhotos, s.validTime, s.validTime, true,
		},
		{
			"Empty authorID",
			s.validID, s.validTitle, *s.validContent, s.validTags, uuid.Nil, s.validPhotos, s.validTime, s.validTime, true,
		},
		{
			"Nil tags",
			s.validID, s.validTitle, *s.validContent, nil, s.validAuthorID, s.validPhotos, s.validTime, s.validTime, true,
		},
		{
			"Empty tag",
			s.validID, s.validTitle, *s.validContent, []string{""}, s.validAuthorID, s.validPhotos, s.validTime, s.validTime, true,
		},
		{
			"Update before creation",
			s.validID, s.validTitle, *s.validContent, s.validTags, s.validAuthorID, s.validPhotos, s.validTime, s.validTime.Add(-time.Hour), true,
		},
		{
			"Future update time",
			s.validID, s.validTitle, *s.validContent, s.validTags, s.validAuthorID, s.validPhotos, s.validTime, time.Now().Add(time.Hour), true,
		},
		{
			"Valid data",
			s.validID, s.validTitle, *s.validContent, s.validTags, s.validAuthorID, s.validPhotos, s.validTime, s.validTime, false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t provider.T) {
			t.Tags("negative", "validation")

			t.WithNewStep("Prepare test parameters", func(ctx provider.StepCtx) {
				ctx.WithParameters(
					allure.NewParameter("Title", tc.title),
					allure.NewParameter("AuthorID", tc.authorID.String()),
					allure.NewParameter("ExpectError", tc.expectError),
				)
			})

			var err error
			t.WithNewStep("Attempt to create post", func(ctx provider.StepCtx) {
				_, err = CreatePost(
					tc.id,
					tc.title,
					tc.content,
					tc.tags,
					tc.authorID,
					tc.photos,
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

func (s *PostTestSuite) TestNewPostSuccess(t provider.T) {
	t.Tags("positive", "creation")
	t.Description("Successful creation of post with auto-generated fields")

	var post *Post
	var err error

	t.WithNewStep("Create post with NewPost", func(ctx provider.StepCtx) {
		post, err = NewPost(
			s.validTitle,
			*s.validContent,
			s.validTags,
			s.validAuthorID,
			&s.validPhotos,
		)
	})

	t.WithNewStep("Verify auto-generated properties", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, post.ID())
		assert.True(t, post.CreatedAt().Before(time.Now()) || post.CreatedAt().Equal(time.Now()))
		assert.True(t, post.UpdatedAt().Before(time.Now()) || post.UpdatedAt().Equal(time.Now()))
	})
}

func (s *PostTestSuite) TestUpdateTitle(t provider.T) {
	t.Tags("positive", "update")
	t.Description("Successful update of post title")

	post := &Post{title: "Old Title", updatedAt: s.validTime}
	newTitle := "New Title"

	var err error
	t.WithNewStep("Update post title", func(ctx provider.StepCtx) {
		err = post.UpdateTitle(newTitle)
	})

	t.WithNewStep("Verify title update", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, newTitle, post.title)
		assert.True(t, post.updatedAt.After(s.validTime))
	})
}

func (s *PostTestSuite) TestUpdateContent(t provider.T) {
	t.Tags("positive", "update")
	t.Description("Successful update of post content")

	post := &Post{content: *s.validContent, updatedAt: s.validTime}
	newContent, _ := NewContent("New content", ContentTypePlainText)

	var err error
	t.WithNewStep("Update post content", func(ctx provider.StepCtx) {
		err = post.UpdateContent(*newContent)
	})

	t.WithNewStep("Verify content update", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, *newContent, post.content)
		assert.True(t, post.updatedAt.After(s.validTime))
	})
}

func (s *PostTestSuite) TestUpdateTags(t provider.T) {
	t.Tags("positive", "update")
	t.Description("Successful update of post tags")

	post := &Post{tags: s.validTags, updatedAt: s.validTime}
	newTags := []string{"new1", "new2"}

	var err error
	t.WithNewStep("Update post tags", func(ctx provider.StepCtx) {
		err = post.UpdateTags(newTags)
	})

	t.WithNewStep("Verify tags update", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, newTags, post.tags)
		assert.True(t, post.updatedAt.After(s.validTime))
	})
}

func (s *PostTestSuite) TestAddTag(t provider.T) {
	t.Tags("positive", "tags")
	t.Description("Test tag addition functionality")

	t.Run("Success", func(t provider.T) {
		post := &Post{tags: s.validTags}
		newTag := "newTag"

		var err error
		t.WithNewStep("Add new tag", func(ctx provider.StepCtx) {
			err = post.AddTag(newTag)
		})

		t.WithNewStep("Verify tag addition", func(ctx provider.StepCtx) {
			require.NoError(t, err)
			assert.Contains(t, post.tags, newTag)
		})
	})

	t.Run("Duplicate", func(t provider.T) {
		post := &Post{tags: s.validTags}

		var err error
		t.WithNewStep("Attempt to add duplicate tag", func(ctx provider.StepCtx) {
			err = post.AddTag(s.validTags[0])
		})

		t.WithNewStep("Verify duplicate error", func(ctx provider.StepCtx) {
			assert.Error(t, err)
		})
	})

	t.Run("Max tags", func(t provider.T) {
		tags := make([]string, MaximumTagCount)
		for i := range tags {
			tags[i] = fmt.Sprintf("tag%d", i)
		}
		post := &Post{tags: tags}

		var err error
		t.WithNewStep("Attempt to exceed maximum tags", func(ctx provider.StepCtx) {
			err = post.AddTag("new")
		})

		t.WithNewStep("Verify max tags error", func(ctx provider.StepCtx) {
			assert.Error(t, err)
		})
	})
}

func (s *PostTestSuite) TestRemoveTag(t provider.T) {
	t.Tags("positive", "tags")
	t.Description("Test tag removal functionality")

	t.Run("Success", func(t provider.T) {
		post := &Post{tags: s.validTags}

		var err error
		t.WithNewStep("Remove existing tag", func(ctx provider.StepCtx) {
			err = post.RemoveTag(s.validTags[0])
		})

		t.WithNewStep("Verify tag removal", func(ctx provider.StepCtx) {
			require.NoError(t, err)
			assert.NotContains(t, post.tags, s.validTags[0])
		})
	})

	t.Run("Not found", func(t provider.T) {
		post := &Post{tags: s.validTags}

		var err error
		t.WithNewStep("Attempt to remove non-existent tag", func(ctx provider.StepCtx) {
			err = post.RemoveTag("nonexistent")
		})

		t.WithNewStep("Verify not found error", func(ctx provider.StepCtx) {
			assert.Error(t, err)
		})
	})
}

func (s *PostTestSuite) TestAddPhoto(t provider.T) {
	t.Tags("negative", "photos")
	t.Description("Test photo addition error handling")

	post := &Post{photos: PostPhotos{}}
	photo := &PostPhoto{}

	var err error
	t.WithNewStep("Attempt to add photo", func(ctx provider.StepCtx) {
		err = post.AddPhoto(photo)
	})

	t.WithNewStep("Verify photo addition error", func(ctx provider.StepCtx) {
		assert.Error(t, err)
	})
}

func (s *PostTestSuite) TestGetters(t provider.T) {
	t.Tags("getters", "functionality")
	t.Description("Test post getter methods")

	post := &Post{
		id:        s.validID,
		title:     s.validTitle,
		content:   *s.validContent,
		tags:      s.validTags,
		authorID:  s.validAuthorID,
		photos:    s.validPhotos,
		createdAt: s.validTime,
		updatedAt: s.validTime,
	}

	t.WithNewStep("Verify all getter values", func(ctx provider.StepCtx) {
		assert.Equal(t, s.validID, post.ID())
		assert.Equal(t, s.validTitle, post.Title())
		assert.Equal(t, *s.validContent, post.Content())
		assert.Equal(t, s.validTags, post.Tags())
		assert.Equal(t, s.validAuthorID, post.AuthorID())
		assert.Equal(t, s.validPhotos, post.Photos())
		assert.Equal(t, s.validTime, post.CreatedAt())
		assert.Equal(t, s.validTime, post.UpdatedAt())
	})
}

func TestPost(t *testing.T) {
	suite.RunSuite(t, new(PostTestSuite))
}
