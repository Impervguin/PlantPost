package albumstorage

import (
	"PlantSite/internal/infra/sqpgx"
	"PlantSite/internal/models/album"
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type PostgresAlbumRepository struct {
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

func NewPostgresAlbumRepository(ctx context.Context) (*PostgresAlbumRepository, error) {
	connStr := configString()
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	return &PostgresAlbumRepository{db: sqpgx.NewSquirrelPgx(pool)}, err
}

type Album struct {
	ID          uuid.UUID
	Name        string
	Description string
	OwnerID     uuid.UUID
	PlantIDs    uuid.UUIDs
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (repo *PostgresAlbumRepository) Create(ctx context.Context, alb *album.Album) (*album.Album, error) {
	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Insert(ctx, squirrel.Insert("album").
		Columns("id", "name", "description", "owner_id", "created_at", "updated_at").
		Values(alb.ID(), alb.Name(), alb.Description(), alb.GetOwnerID(), alb.CreatedAt(), alb.UpdatedAt()),
	)

	if err != nil {
		return nil, err
	}

	query := squirrel.Insert("plant_album").
		Columns("id", "album_id", "plant_id")

	for _, id := range alb.PlantIDs() {
		query.Values(uuid.New(), alb.ID(), id)
	}

	_, err = tx.Insert(ctx, query)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)

	return alb, err
}

func (repo *PostgresAlbumRepository) Get(ctx context.Context, id uuid.UUID) (*album.Album, error) {
	var tmpAlbum Album
	err := repo.db.QueryRow(ctx, squirrel.Select("id", "name", "description", "owner_id", "created_at", "updated_at").
		From("album").
		Where(squirrel.Eq{"id": id}),
	).Scan(&tmpAlbum.ID, &tmpAlbum.Name, &tmpAlbum.Description, &tmpAlbum.OwnerID, &tmpAlbum.CreatedAt, &tmpAlbum.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, album.ErrAlbumNotFound
	} else if err != nil {
		return nil, err
	}
	var plantIDs uuid.UUIDs
	rows, err := repo.db.Query(ctx, squirrel.Select("plant_id").
		From("plant_album").
		Where(squirrel.Eq{"album_id": id}),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var plantID uuid.UUID
		err := rows.Scan(&plantID)
		if err != nil {
			return nil, err
		}
		plantIDs = append(plantIDs, plantID)
	}
	tmpAlbum.PlantIDs = plantIDs

	alb, err := album.CreateAlbum(
		tmpAlbum.ID,
		tmpAlbum.Name,
		tmpAlbum.Description,
		tmpAlbum.PlantIDs,
		tmpAlbum.OwnerID,
		tmpAlbum.CreatedAt,
		tmpAlbum.UpdatedAt,
	)
	return alb, err
}

func (repo *PostgresAlbumRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(*album.Album) (*album.Album, error)) (*album.Album, error) {
	alb, err := repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	alb, err = updateFn(alb)
	if err != nil {
		return nil, err
	}

	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Update(ctx, squirrel.Update("album").
		Set("name", alb.Name()).
		Set("description", alb.Description()).
		Set("owner_id", alb.GetOwnerID()).
		Set("updated_at", alb.UpdatedAt()).
		Where(squirrel.Eq{"id": alb.ID()}),
	)
	if err != nil {
		return nil, err
	}
	_, err = tx.Delete(ctx, squirrel.Delete("plant_album").
		Where(squirrel.Eq{"album_id": id}),
	)

	if err != nil {
		return nil, err
	}

	query := squirrel.Insert("plant_album").
		Columns("id", "album_id", "plant_id")
	for _, id := range alb.PlantIDs() {
		query.Values(uuid.New(), alb.ID(), id)
	}

	_, err = tx.Insert(ctx, query)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)

	return alb, err
}

func (repo *PostgresAlbumRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Delete(ctx, squirrel.Delete("plant_album").
		Where(squirrel.Eq{"album_id": id}),
	)

	if err != nil {
		return err
	}

	_, err = tx.Delete(ctx, squirrel.Delete("album").
		Where(squirrel.Eq{"id": id}),
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (repo *PostgresAlbumRepository) List(ctx context.Context, ownerID uuid.UUID) ([]*album.Album, error) {
	var albums []*album.Album
	rows, err := repo.db.Query(ctx, squirrel.Select("id", "name", "description", "owner_id", "created_at", "updated_at").
		From("album").
		Where(squirrel.Eq{"owner_id": ownerID}),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tmpAlbum Album
		err := rows.Scan(&tmpAlbum.ID, &tmpAlbum.Name, &tmpAlbum.Description, &tmpAlbum.OwnerID, &tmpAlbum.CreatedAt, &tmpAlbum.UpdatedAt)
		if err != nil {
			return nil, err
		}
		var plantIDs uuid.UUIDs
		subRows, err := repo.db.Query(ctx, squirrel.Select("plant_id").
			From("plant_album").
			Where(squirrel.Eq{"album_id": tmpAlbum.ID}),
		)
		if err != nil {
			return nil, err
		}
		defer subRows.Close()
		for subRows.Next() {
			var plantID uuid.UUID
			err := subRows.Scan(&plantID)
			if err != nil {
				return nil, err
			}
			plantIDs = append(plantIDs, plantID)
		}
		tmpAlbum.PlantIDs = plantIDs
		alb, err := album.CreateAlbum(
			tmpAlbum.ID,
			tmpAlbum.Name,
			tmpAlbum.Description,
			tmpAlbum.PlantIDs,
			tmpAlbum.OwnerID,
			tmpAlbum.CreatedAt,
			tmpAlbum.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		albums = append(albums, alb)
	}
	return albums, nil
}
