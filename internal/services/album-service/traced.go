package albumservice

import (
	"PlantSite/internal/models/album"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TracedAlbumService struct {
	inner  AlbumServiceContract
	tracer trace.Tracer
}

func NewTracedAlbumService(inner AlbumServiceContract) *TracedAlbumService {
	return &TracedAlbumService{
		inner:  inner,
		tracer: otel.Tracer("search-service"),
	}
}

var _ AlbumServiceContract = (*TracedAlbumService)(nil)

func (s *TracedAlbumService) CreateAlbum(ctx context.Context, alb *album.Album) (*album.Album, error) {
	ctx, span := s.tracer.Start(ctx, "CreateAlbum")
	defer span.End()

	span.SetAttributes(
		attribute.String("name", alb.Name()),
		attribute.String("description", alb.Description()),
		attribute.StringSlice("plant_ids", alb.PlantIDs().Strings()),
		attribute.String("owner_id", alb.GetOwnerID().String()),
		attribute.String("created_at", alb.CreatedAt().String()),
		attribute.String("updated_at", alb.UpdatedAt().String()),
	)

	alb, err := s.inner.CreateAlbum(ctx, alb)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return alb, nil
}

func (s *TracedAlbumService) GetAlbum(ctx context.Context, id uuid.UUID) (*album.Album, error) {
	ctx, span := s.tracer.Start(ctx, "GetAlbum")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
	)

	alb, err := s.inner.GetAlbum(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("name", alb.Name()),
		attribute.String("description", alb.Description()),
		attribute.StringSlice("plant_ids", alb.PlantIDs().Strings()),
		attribute.String("owner_id", alb.GetOwnerID().String()),
		attribute.String("created_at", alb.CreatedAt().String()),
	)
	return alb, nil
}

func (s *TracedAlbumService) UpdateAlbumName(ctx context.Context, id uuid.UUID, name string) error {
	ctx, span := s.tracer.Start(ctx, "UpdateAlbumName")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
		attribute.String("name", name),
	)

	err := s.inner.UpdateAlbumName(ctx, id, name)
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (s *TracedAlbumService) UpdateAlbumDescription(ctx context.Context, id uuid.UUID, description string) error {
	ctx, span := s.tracer.Start(ctx, "UpdateAlbumDescription")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
		attribute.String("description", description),
	)

	err := s.inner.UpdateAlbumDescription(ctx, id, description)
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (s *TracedAlbumService) AddPlantToAlbum(ctx context.Context, id uuid.UUID, plantID uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "AddPlantToAlbum")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
		attribute.String("plant_id", plantID.String()),
	)

	err := s.inner.AddPlantToAlbum(ctx, id, plantID)
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (s *TracedAlbumService) RemovePlantFromAlbum(ctx context.Context, id uuid.UUID, plantID uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "RemovePlantFromAlbum")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
		attribute.String("plant_id", plantID.String()),
	)

	err := s.inner.RemovePlantFromAlbum(ctx, id, plantID)
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (s *TracedAlbumService) DeleteAlbum(ctx context.Context, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "DeleteAlbum")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
	)

	err := s.inner.DeleteAlbum(ctx, id)
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (s *TracedAlbumService) ListAlbums(ctx context.Context) ([]*album.Album, error) {
	ctx, span := s.tracer.Start(ctx, "ListAlbums")
	defer span.End()

	albs, err := s.inner.ListAlbums(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.Int("albums_count", len(albs)),
	)
	return albs, nil
}
