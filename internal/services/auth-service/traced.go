package authservice

import (
	"PlantSite/internal/models/auth"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TracedAuthService struct {
	inner  AuthServiceContract
	tracer trace.Tracer
}

func NewTracedAuthService(inner AuthServiceContract) *TracedAuthService {
	return &TracedAuthService{
		inner:  inner,
		tracer: otel.Tracer("auth-service"),
	}
}

func (s *TracedAuthService) Login(ctx context.Context, identifier, password string) (uuid.UUID, error) {
	ctx, span := s.tracer.Start(ctx, "Login")
	defer span.End()

	span.SetAttributes(
		attribute.String("identifier", identifier),
	)
	sid, err := s.inner.Login(ctx, identifier, password)
	if err != nil {
		span.RecordError(err)
		return sid, err
	}
	span.SetAttributes(
		attribute.String("session_id", sid.String()),
	)
	return sid, nil
}

func (s *TracedAuthService) Register(ctx context.Context, name, email, password string) error {
	ctx, span := s.tracer.Start(ctx, "Register")
	defer span.End()

	span.SetAttributes(
		attribute.String("name", name),
		attribute.String("email", email),
	)
	err := s.inner.Register(ctx, name, email, password)
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (s *TracedAuthService) Logout(ctx context.Context) error {
	ctx, span := s.tracer.Start(ctx, "Logout")
	defer span.End()

	err := s.inner.Logout(ctx)
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (s *TracedAuthService) Authenticate(ctx context.Context, sid uuid.UUID) context.Context {
	ctx, span := s.tracer.Start(ctx, "Authenticate")
	defer span.End()

	ctx = s.inner.Authenticate(ctx, sid)

	span.SetAttributes(
		attribute.String("session_id", sid.String()),
	)
	return ctx
}

func (s *TracedAuthService) UserFromContext(ctx context.Context) auth.User {
	ctx, span := s.tracer.Start(ctx, "UserFromContext")
	defer span.End()

	user := s.inner.UserFromContext(ctx)

	span.SetAttributes(
		attribute.String("user_id", user.ID().String()),
	)
	return user
}
