package authstorage

import (
	"PlantSite/internal/infra/sqdb"
	"PlantSite/internal/models/auth"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type PostgresAuthRepository struct {
	db sqdb.SquirrelDatabase
}

func NewPostgresAuthRepository(_ context.Context, db sqdb.SquirrelDatabase) (*PostgresAuthRepository, error) {
	return &PostgresAuthRepository{db: db}, nil
}

type Member struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash []byte
	CreatedAt    time.Time
}

type Author struct {
	ID         uuid.UUID
	Rights     bool
	GiveTime   time.Time
	RevokeTime time.Time
}

func (repo *PostgresAuthRepository) getMember(ctx context.Context, whereStatement interface{}, args ...interface{}) (*auth.Member, error) {
	var mem *Member = &Member{}
	row, err := repo.db.QueryRow(ctx,
		squirrel.Select("id", "username", "email", "password_hash", "created_at").
			From("app_user").
			Where(whereStatement, args...),
	)

	if err != nil {
		return nil, fmt.Errorf("PostgresAuthRepository.getMember failed %w", err)
	}

	err = row.Scan(&mem.ID, &mem.Name, &mem.Email, &mem.PasswordHash, &mem.CreatedAt)

	if err == sqdb.ErrNoRows {
		return nil, auth.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("PostgresAuthRepository.getMember failed %w", err)
	}

	domainMem, err := auth.CreateMember(
		mem.ID,
		mem.Name,
		mem.Email,
		mem.PasswordHash,
		mem.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("PostgresAuthRepository.getMember failed %w", err)
	}

	return domainMem, err
}

func (repo *PostgresAuthRepository) getAuthor(ctx context.Context, mem *auth.Member) (*auth.Author, error) {
	var aut *Author = &Author{}
	row, err := repo.db.QueryRow(ctx,
		squirrel.Select("id", "grant_at", "has_rights", "revoke_at").
			From("author").
			Where(squirrel.Eq{`"id"`: mem.ID()}),
	)

	if err != nil {
		return nil, fmt.Errorf("PostgresAuthRepository.getAuthor failed %w", err)
	}

	err = row.Scan(&aut.ID, &aut.GiveTime, &aut.Rights, &aut.RevokeTime)

	if err == sqdb.ErrNoRows {
		return nil, auth.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	domainAuth, err := auth.CreateAuthor(
		*mem,
		aut.GiveTime,
		aut.Rights,
		aut.RevokeTime,
	)
	if err != nil {
		return nil, err
	}

	return domainAuth, err
}

func (repo *PostgresAuthRepository) getUser(ctx context.Context, whereStatement interface{}, args ...interface{}) (auth.User, error) {
	domainMem, err := repo.getMember(ctx, whereStatement, args...)
	if errors.Is(err, auth.ErrUserNotFound) {
		return nil, auth.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	domainAuth, err := repo.getAuthor(ctx, domainMem)
	if errors.Is(err, auth.ErrUserNotFound) {
		return domainMem, nil
	} else if err != nil {
		return nil, err
	}

	return domainAuth, err

}

func (repo *PostgresAuthRepository) Get(ctx context.Context, id uuid.UUID) (auth.User, error) {
	return repo.getUser(ctx, squirrel.Eq{"id": id.String()})
}

func (repo *PostgresAuthRepository) Create(ctx context.Context, mem *auth.Member) (auth.User, error) {
	_, err := repo.db.Insert(ctx,
		squirrel.Insert("app_user").
			Columns("id", "username", "email", "password_hash", "created_at").
			Values(mem.ID(), mem.Name(), mem.Email(), mem.HashedPassword(), mem.CreatedAt()),
	)
	if err != nil {
		return nil, err
	}

	return mem, nil
}

func (repo *PostgresAuthRepository) updateMember(ctx context.Context, updMember *auth.Member) error {
	_, err := repo.db.Update(ctx,
		squirrel.Update("app_user").
			Set("username", updMember.Name()).
			Set("email", updMember.Email()).
			Set("password_hash", updMember.HashedPassword()).
			Where(squirrel.Eq{"id": updMember.ID()}),
	)
	return err
}

func (repo *PostgresAuthRepository) updateAuthor(ctx context.Context, updAuth *auth.Author) error {
	_, err := repo.db.Insert(ctx, squirrel.Insert("author").
		Columns("id", "has_rights", "grant_at", "revoke_at").
		Values(updAuth.ID(), updAuth.HasRights(), updAuth.GiveTime(), updAuth.RevokeTime()).
		Suffix("ON CONFLICT (id) DO UPDATE SET has_rights = ?, grant_at = ?, revoke_at = ?", updAuth.HasRights(), updAuth.GiveTime(), updAuth.RevokeTime()),
	)
	if err != nil {
		return err
	}
	return err
}

func (repo *PostgresAuthRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(auth.User) (auth.User, error)) (auth.User, error) {
	usr, err := repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	updUsr, err := updateFn(usr)
	if err != nil {
		return nil, err
	}
	switch fact := updUsr.(type) {
	case *auth.Member:
		err = repo.updateMember(ctx, fact)
		if err != nil {
			return nil, err
		}
		return updUsr, nil
	case *auth.Author:
		err = repo.updateAuthor(ctx, fact)
		if err != nil {
			return nil, err
		}
		err = repo.updateMember(ctx, &fact.Member)
		if err != nil {
			return nil, err
		}
		return updUsr, nil
	default:
		return nil, fmt.Errorf("unsupported user type: %v", updUsr)
	}
}

func (repo *PostgresAuthRepository) GetByName(ctx context.Context, name string) (auth.User, error) {
	return repo.getUser(ctx, squirrel.Eq{`"username"`: name})
}

func (repo *PostgresAuthRepository) GetByEmail(ctx context.Context, email string) (auth.User, error) {
	return repo.getUser(ctx, squirrel.Eq{`"email"`: email})
}
