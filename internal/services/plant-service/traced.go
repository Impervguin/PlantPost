package plantservice

import (
	"PlantSite/internal/models"
	"PlantSite/internal/models/plant"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TracedPlantService struct {
	inner  PlantServiceContract
	tracer trace.Tracer
}

func NewTracedPlantService(inner PlantServiceContract) *TracedPlantService {
	return &TracedPlantService{
		inner:  inner,
		tracer: otel.Tracer("plant-service"),
	}
}

var _ PlantServiceContract = (*TracedPlantService)(nil)

func (s *TracedPlantService) UpdatePlantSpec(ctx context.Context, id uuid.UUID, spec plant.PlantSpecification) error {
	ctx, span := s.tracer.Start(ctx, "UpdatePlantSpec")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
		attribute.String("category", spec.Category()),
	)

	err := s.inner.UpdatePlantSpec(ctx, id, spec)
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (s *TracedPlantService) DeletePlant(ctx context.Context, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "DeletePlant")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
	)

	err := s.inner.DeletePlant(ctx, id)
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (s *TracedPlantService) UploadPlantPhoto(
	ctx context.Context,
	id uuid.UUID,
	fdata models.FileData,
	description string,
) error {
	ctx, span := s.tracer.Start(ctx, "UploadPlantPhoto")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
		attribute.String("description", description),
		attribute.String("file_name", fdata.Name),
		attribute.String("content_type", fdata.ContentType),
	)

	err := s.inner.UploadPlantPhoto(ctx, id, fdata, description)
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (s *TracedPlantService) GetPlantCategory(ctx context.Context, name string) (*plant.PlantCategory, error) {
	ctx, span := s.tracer.Start(ctx, "GetPlantCategory")
	defer span.End()

	span.SetAttributes(
		attribute.String("name", name),
	)

	cat, err := s.inner.GetPlantCategory(ctx, name)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return cat, nil
}

func (s *TracedPlantService) ListCategories(ctx context.Context) ([]plant.PlantCategory, error) {
	ctx, span := s.tracer.Start(ctx, "ListCategories")
	defer span.End()

	cats, err := s.inner.ListCategories(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	catstrs := make([]string, 0, len(cats))
	for _, cat := range cats {
		catstrs = append(catstrs, cat.Name)
	}

	span.SetAttributes(
		attribute.StringSlice("categories", catstrs),
	)
	return cats, nil
}

func (s *TracedPlantService) GetPlant(ctx context.Context, id uuid.UUID) (*GetPlant, error) {
	ctx, span := s.tracer.Start(ctx, "GetPlant")
	defer span.End()

	span.SetAttributes(
		attribute.String("id", id.String()),
	)

	pl, err := s.inner.GetPlant(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("name", pl.Name),
		attribute.String("latin_name", pl.LatinName),
		attribute.String("description", pl.Description),
		attribute.String("category", pl.Category),
		attribute.String("created_at", pl.CreatedAt.String()),
	)
	return pl, nil
}

func (s *TracedPlantService) CreatePlant(
	ctx context.Context,
	data CreatePlantData,
	mainPhotoFile models.FileData,
) error {
	ctx, span := s.tracer.Start(ctx, "CreatePlant")
	defer span.End()

	span.SetAttributes(
		attribute.String("name", data.Name),
		attribute.String("latin_name", data.LatinName),
		attribute.String("description", data.Description),
		attribute.String("category", data.Category),
	)

	span.SetAttributes(
		attribute.String("photo_name", mainPhotoFile.Name),
		attribute.String("photo_content_type", mainPhotoFile.ContentType),
	)

	err := s.inner.CreatePlant(ctx, data, mainPhotoFile)
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}
