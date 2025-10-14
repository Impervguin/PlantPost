//go:build unit

package search

import (
	"PlantSite/internal/models/post"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PostBuilder реализует паттерн Data Builder для создания тестовых постов
type PostBuilder struct {
	title    string
	tags     []string
	authorID uuid.UUID
}

func NewPostBuilder() *PostBuilder {
	return &PostBuilder{
		title:    "Default Post",
		tags:     []string{"default"},
		authorID: uuid.New(),
	}
}

func (b *PostBuilder) WithTitle(title string) *PostBuilder {
	b.title = title
	return b
}

func (b *PostBuilder) WithTags(tags []string) *PostBuilder {
	b.tags = tags
	return b
}

func (b *PostBuilder) WithAuthorID(authorID uuid.UUID) *PostBuilder {
	b.authorID = authorID
	return b
}

func (b *PostBuilder) Build() (*post.Post, error) {
	content, _ := post.NewContent("Test content", post.ContentTypePlainText)
	photos := post.NewPostPhotos()

	return post.NewPost(
		b.title,
		*content,
		b.tags,
		b.authorID,
		photos,
	)
}

type PostFiltersTestSuite struct {
	suite.Suite
	testPost1 *post.Post
	testPost2 *post.Post
	testPost3 *post.Post
	authorID1 uuid.UUID
	authorID2 uuid.UUID
}

func (s *PostFiltersTestSuite) BeforeEach(t provider.T) {
	t.Epic("Search")
	t.Feature("Post Filters")

	s.authorID1 = uuid.New()
	s.authorID2 = uuid.New()

	builder := NewPostBuilder()

	var err error
	s.testPost1, err = builder.
		WithTitle("First Post").
		WithTags([]string{"tech", "golang"}).
		WithAuthorID(s.authorID1).
		Build()
	require.NoError(t, err)

	s.testPost2, err = builder.
		WithTitle("Second Post About Programming").
		WithTags([]string{"programming", "java"}).
		WithAuthorID(s.authorID2).
		Build()
	require.NoError(t, err)

	s.testPost3, err = builder.
		WithTitle("Third Post").
		WithTags([]string{"tech", "python"}).
		WithAuthorID(s.authorID1).
		Build()
	require.NoError(t, err)
}

func (s *PostFiltersTestSuite) TestPostTitleFilter(t provider.T) {
	t.Tags("filter", "title", "exact")
	t.Description("Test exact post title filtering")

	filter := NewPostTitleFilter("First Post")
	t.WithNewStep("Filter by exact title 'First Post'", func(ctx provider.StepCtx) {
		assert.True(t, filter.Filter(s.testPost1))
		assert.False(t, filter.Filter(s.testPost2))
		assert.False(t, filter.Filter(s.testPost3))
	})
}

func (s *PostFiltersTestSuite) TestPostTitleContainsFilter(t provider.T) {
	t.Tags("filter", "title", "partial")
	t.Description("Test partial post title matching")

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
		t.Run(tt.name, func(t provider.T) {
			t.WithNewStep("Prepare title contains filter", func(ctx provider.StepCtx) {
				ctx.WithParameters(
					allure.NewParameter("SearchPart", tt.part),
				)
			})

			filter := NewPostTitleContainsFilter(tt.part)
			t.WithNewStep("Apply title contains filter", func(ctx provider.StepCtx) {
				assert.Equal(t, tt.post1, filter.Filter(s.testPost1))
				assert.Equal(t, tt.post2, filter.Filter(s.testPost2))
				assert.Equal(t, tt.post3, filter.Filter(s.testPost3))
			})
		})
	}
}

func (s *PostFiltersTestSuite) TestPostTagFilter(t provider.T) {
	t.Tags("filter", "tags")
	t.Description("Test post tag filtering")

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
		t.Run(tt.name, func(t provider.T) {
			t.WithNewStep("Prepare tag filter", func(ctx provider.StepCtx) {
				ctx.WithParameters(
					allure.NewParameter("Tag", tt.tag),
				)
			})

			filter := NewPostTagFilter([]string{tt.tag})
			t.WithNewStep("Apply tag filter", func(ctx provider.StepCtx) {
				assert.Equal(t, tt.post1, filter.Filter(s.testPost1))
				assert.Equal(t, tt.post2, filter.Filter(s.testPost2))
				assert.Equal(t, tt.post3, filter.Filter(s.testPost3))
			})
		})
	}
}

func (s *PostFiltersTestSuite) TestPostAuthorFilter(t provider.T) {
	t.Tags("filter", "author")
	t.Description("Test post author filtering")

	filter := NewPostAuthorFilter(s.authorID1)
	t.WithNewStep("Filter by first author", func(ctx provider.StepCtx) {
		assert.True(t, filter.Filter(s.testPost1))
		assert.False(t, filter.Filter(s.testPost2))
		assert.True(t, filter.Filter(s.testPost3))
	})

	filter = NewPostAuthorFilter(s.authorID2)
	t.WithNewStep("Filter by second author", func(ctx provider.StepCtx) {
		assert.False(t, filter.Filter(s.testPost1))
		assert.True(t, filter.Filter(s.testPost2))
		assert.False(t, filter.Filter(s.testPost3))
	})
}

func TestPostFilters(t *testing.T) {
	suite.RunSuite(t, new(PostFiltersTestSuite))
}
