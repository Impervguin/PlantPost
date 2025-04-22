package authrepo

import (
	"PlantSite/internal/models/auth"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAdminLoginExists   = errors.New("admin login already exists")
	ErrAdminIdExists      = errors.New("admin id already exists")
	ErrAdminNotUpdateable = errors.New("admin not updateable")
)

var _ auth.AuthRepository = (*WithAdminRepository)(nil)

type AdminRepository interface {
	Add(ctx context.Context, admin auth.Admin) error
	Get(ctx context.Context, login string) (auth.Admin, bool)
	GetByID(ctx context.Context, id uuid.UUID) (auth.Admin, bool)
	Iterate(ctx context.Context, fn func(admin auth.Admin) error) error
}

type WithAdminRepository struct {
	admin AdminRepository
	auth  auth.AuthRepository
}

func NewWithAdminRepository(admin AdminRepository, authRepo auth.AuthRepository) *WithAdminRepository {
	if admin == nil {
		panic("nil admin")
	}
	if authRepo == nil {
		panic("nil auth")
	}

	err := admin.Iterate(context.Background(), func(admin auth.Admin) error {
		usr, err := authRepo.GetByName(context.Background(), admin.Login())
		if errors.Is(err, auth.ErrUserNotFound) {
			memb, err := auth.CreateMember(admin.ID(), admin.Login(), admin.Login()+"@admin.com", admin.HashedPassword(), time.Now())
			if err != nil {
				return err
			}
			usr, err = authRepo.Create(context.Background(), memb)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		_, ok := usr.(*auth.Author)
		if !ok {
			usr, err = authRepo.Update(context.Background(), usr.ID(), func(usr auth.User) (auth.User, error) {
				mem := usr.(*auth.Member)
				return auth.CreateAuthor(*mem, time.Now(), true, time.Time{})
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
	return &WithAdminRepository{
		admin: admin,
		auth:  authRepo,
	}
}

func (r *WithAdminRepository) Get(ctx context.Context, id uuid.UUID) (auth.User, error) {
	adm, exists := r.admin.GetByID(ctx, id)
	if exists {
		return &adm, nil
	}
	return r.auth.Get(ctx, id)
}

func (r *WithAdminRepository) GetByName(ctx context.Context, name string) (auth.User, error) {
	adm, exists := r.admin.Get(ctx, name)
	if exists {
		return &adm, nil
	}
	return r.auth.GetByName(ctx, name)
}

func (r *WithAdminRepository) GetByEmail(ctx context.Context, email string) (auth.User, error) {
	return r.auth.GetByEmail(ctx, email)
}

func (r *WithAdminRepository) Create(ctx context.Context, user *auth.Member) (auth.User, error) {
	_, ok := r.admin.Get(ctx, user.Name())
	if ok {
		return nil, ErrAdminLoginExists
	}
	_, ok = r.admin.GetByID(ctx, user.ID())
	if ok {
		return nil, ErrAdminIdExists
	}
	return r.auth.Create(ctx, user)
}

func (r *WithAdminRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(auth.User) (auth.User, error)) (auth.User, error) {
	_, exists := r.admin.GetByID(ctx, id)
	if exists {
		return nil, ErrAdminNotUpdateable
	}
	return r.auth.Update(ctx, id, updateFn)
}
