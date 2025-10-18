package searchservice

import (
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/post"
	"PlantSite/internal/models/search"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TracedSearchService struct {
	inner  SearchServiceContract
	tracer trace.Tracer
}

func NewTracedSearchService(inner SearchServiceContract) *TracedSearchService {
	return &TracedSearchService{
		inner:  inner,
		tracer: otel.Tracer("search-service"),
	}
}

func (s *TracedSearchService) PostAuthors(ctx context.Context) ([]*auth.Author, error) {
	ctx, span := s.tracer.Start(ctx, "PostAuthors")
	defer span.End()

	authors, err := s.inner.PostAuthors(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	span.SetAttributes(
		attribute.Int("authors_count", len(authors)),
	)
	return authors, nil
}

func (s *TracedSearchService) PostTags(ctx context.Context) ([]string, error) {
	ctx, span := s.tracer.Start(ctx, "PostTags")
	defer span.End()

	tags, err := s.inner.PostTags(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	span.SetAttributes(
		attribute.Int("tags_count", len(tags)),
	)
	return tags, nil
}

func (s *TracedSearchService) SearchPosts(ctx context.Context, plSearch *search.PostSearch) ([]*SearchPost, error) {
	ctx, span := s.tracer.Start(ctx, "SearchPosts")
	defer span.End()

	postFilters := make([]string, 0)
	plSearch.Iterate(func(f search.PostFilter) error {
		postFilters = append(postFilters, f.Identifier())
		return nil
	})

	posts, err := s.inner.SearchPosts(ctx, plSearch)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	span.SetAttributes(
		attribute.Int("posts_count", len(posts)),
	)
	return posts, nil
}

func (s *TracedSearchService) SearchPlants(ctx context.Context, plSearch *search.PlantSearch) ([]*SearchPlant, error) {
	ctx, span := s.tracer.Start(ctx, "SearchPlants")
	defer span.End()

	plantFilters := make([]string, 0)
	plSearch.Iterate(func(f search.PlantFilter) error {
		plantFilters = append(plantFilters, f.Identifier())
		return nil
	})

	plants, err := s.inner.SearchPlants(ctx, plSearch)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	span.SetAttributes(
		attribute.Int("plants_count", len(plants)),
	)
	return plants, nil
}

func (s *TracedSearchService) GetPostByID(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	ctx, span := s.tracer.Start(ctx, "GetPostByID")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
	)

	post, err := s.inner.GetPostByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("title", post.Title()),
		attribute.StringSlice("tags", post.Tags()),
		attribute.String("author_id", post.AuthorID().String()),
		attribute.String("created_at", post.CreatedAt().String()),
	)
	return post, nil
}

func (s *TracedSearchService) GetPost(ctx context.Context, id uuid.UUID) (*GetPost, error) {
	ctx, span := s.tracer.Start(ctx, "GetPost")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
	)

	post, err := s.inner.GetPost(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("title", post.Title),
		attribute.StringSlice("tags", post.Tags),
		attribute.String("author_id", post.AuthorID.String()),
		attribute.String("created_at", post.CreatedAt.String()),
	)
	return post, nil
}

func (s *TracedSearchService) GetPlantByID(ctx context.Context, id uuid.UUID) (*GetPlant, error) {
	ctx, span := s.tracer.Start(ctx, "GetPlantByID")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
	)

	plant, err := s.inner.GetPlantByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("name", plant.Name),
		attribute.String("latin_name", plant.LatinName),
		attribute.String("category", plant.Category),
		attribute.String("created_at", plant.CreatedAt.String()),
	)
	return plant, nil
}
