package authstorage

import (
	"PlantSite/internal/infra/sqpgx"
	"PlantSite/internal/models/auth"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type PostgresAuthRepository struct {
	db *sqpgx.SquirrelPgx
}

func configString() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s port=%d host=%s pool_max_conns=%d pool_max_conn_lifetime=%s",
		viper.GetString(ConfigPostgresUserKey),
		viper.GetString(ConfigPostgresPasswordKey),
		viper.GetString(ConfigPostgresDbNameKey),
		viper.GetInt(ConfigPostgresPortKey),
		viper.GetString(ConfigPostgresHostKey),
		viper.GetInt(ConfigMaxConnectionsKey),
		viper.GetString(ConfigMaxConnectionLifetimeKey),
	)
}

func NewPostgresAuthRepository(ctx context.Context) (*PostgresAuthRepository, error) {
	connStr := configString()
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	return &PostgresAuthRepository{db: sqpgx.NewSquirrelPgx(pool)}, nil
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
	err := repo.db.QueryRow(ctx,
		squirrel.Select("id", "username", "email", "password_hash", "created_at").
			From("user").
			Where(whereStatement, args...),
	).Scan(&mem.ID, &mem.Name, &mem.Email, &mem.PasswordHash, &mem.CreatedAt)

	if err == pgx.ErrNoRows {
		return nil, auth.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	domainMem, err := auth.CreateMember(
		mem.ID,
		mem.Name,
		mem.Email,
		mem.PasswordHash,
		mem.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return domainMem, err
}

func (repo *PostgresAuthRepository) getAuthor(ctx context.Context, mem *auth.Member, whereStatement interface{}, args ...interface{}) (*auth.Author, error) {
	var aut *Author = &Author{}
	err := repo.db.QueryRow(ctx,
		squirrel.Select("id", "grant_at", "has_rights", "revoke_at").
			From("author").
			Where(whereStatement, args...),
	).Scan(&aut.ID, &aut.GiveTime, &aut.Rights, &aut.RevokeTime)

	if err == pgx.ErrNoRows {
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

	domainAuth, err := repo.getAuthor(ctx, domainMem, whereStatement, args...)
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
		squirrel.Insert("user").
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
		squirrel.Update("user").
			Set("username", updMember.Name()).
			Set("email", updMember.Email()).
			Set("password_hash", updMember.HashedPassword()).
			Where(squirrel.Eq{"id": updMember.ID().String()}),
	)
	return err
}

func (repo *PostgresAuthRepository) updateAuthor(ctx context.Context, updAuth *auth.Author) error {
	_, err := repo.db.Update(ctx,
		squirrel.Update("author").
			Set("grant_at", updAuth.GiveTime).
			Set("has_rights", updAuth.HasRights()).
			Set("revoke_at", updAuth.RevokeTime).
			Where(squirrel.Eq{"id": updAuth.ID().String()}),
	)
	return err
}

func (repo *PostgresAuthRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(auth.User) (auth.User, error)) (auth.User, error) {
	usr, err := repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	updUsr, err := updateFn(usr.(*auth.Member))
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
		return nil, fmt.Errorf("Unsupported user type: %v", updUsr)
	}
}

func (repo *PostgresAuthRepository) GetByName(ctx context.Context, name string) (auth.User, error) {
	return repo.getUser(ctx, squirrel.Eq{"username": name})
}

func (repo *PostgresAuthRepository) GetByEmail(ctx context.Context, email string) (auth.User, error) {
	return repo.getUser(ctx, squirrel.Eq{"email": email})
}
