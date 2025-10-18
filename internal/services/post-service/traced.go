package postservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/post"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TracedPostService struct {
	inner  PostServiceContract
	tracer trace.Tracer
}

func NewTracedPostService(inner PostServiceContract) *TracedPostService {
	return &TracedPostService{
		inner:  inner,
		tracer: otel.Tracer("search-service"),
	}
}

var _ PostServiceContract = (*TracedPostService)(nil)

func (s *TracedPostService) CreatePost(ctx context.Context, data CreatePostTextData, files []models.FileData) (*post.Post, error) {
	ctx, span := s.tracer.Start(ctx, "CreatePost")
	defer span.End()

	span.SetAttributes(
		attribute.String("title", data.Title),
		attribute.StringSlice("tags", data.Tags),
	)

	post, err := s.inner.CreatePost(ctx, data, files)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return post, nil
}

func (s *TracedPostService) GetPost(ctx context.Context, id uuid.UUID) (*GetPost, error) {
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

func (s *TracedPostService) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "Delete")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
	)

	err := s.inner.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (s *TracedPostService) UpdatePost(ctx context.Context, id uuid.UUID, data UpdatePostTextData) (*post.Post, error) {
	ctx, span := s.tracer.Start(ctx, "UpdatePost")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
		attribute.String("title", data.Title),
		attribute.StringSlice("tags", data.Tags),
	)

	post, err := s.inner.UpdatePost(ctx, id, data)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return post, nil
}
