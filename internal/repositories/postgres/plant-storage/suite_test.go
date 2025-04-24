//go:build integration

package plantstorage_test

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	minioclient "PlantSite/internal/infra/minio-client"
	"PlantSite/internal/infra/sqpgx"
	"PlantSite/internal/models"
	"PlantSite/internal/models/plant"
	filestorage "PlantSite/internal/repositories/pgminio/file-storage"
	plantstorage "PlantSite/internal/repositories/postgres/plant-storage"
	"PlantSite/internal/repositories/tests"
	"PlantSite/internal/testutils/miniotest"
	"PlantSite/internal/testutils/pgtest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type PlantRepositoryTestSuite struct {
	suite.Suite
	dbContainer    testcontainers.Container
	minioContainer testcontainers.Container
	db             *sqpgx.SquirrelPgx
	fileRepo       *filestorage.PgMinioStorage
	repo           *plantstorage.PostgresPlantRepository
	prevDir        string
}

func TestPlantRepositorySuite(t *testing.T) {
	suite.Run(t, new(PlantRepositoryTestSuite))
}

func (s *PlantRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	// Save current directory
	prevDir, err := os.Getwd()
	require.NoError(s.T(), err)
	s.prevDir = prevDir

	// Change directory to test working directory
	err = os.Chdir(tests.GetTestWorkingDir())
	require.NoError(s.T(), err)

	// Setup PostgreSQL container
	dbContainer, dbCreds, err := pgtest.NewTestPostgres(ctx)
	require.NoError(s.T(), err)
	s.dbContainer = dbContainer

	// Run migrations
	err = pgtest.Migrate(ctx, &dbCreds)
	require.NoError(s.T(), err)

	// Create database connection
	dbConfig := &sqpgx.SqpgxConfig{
		User:                   dbCreds.User,
		Password:               dbCreds.Password,
		DbName:                 dbCreds.Database,
		Host:                   dbCreds.Host,
		Port:                   dbCreds.Port,
		MaxConnections:         10,
		MaxConnectionsLifetime: time.Minute,
	}
	s.db, err = sqpgx.NewSquirrelPgx(ctx, dbConfig)
	require.NoError(s.T(), err)

	// Setup MinIO container
	minioContainer, minioCreds, err := miniotest.NewTestMinio(ctx)
	require.NoError(s.T(), err)
	s.minioContainer = minioContainer

	// Run migrations
	err = miniotest.Migrate(ctx, minioCreds)
	require.NoError(s.T(), err)

	// Create MinIO client
	minioConfig, err := minioclient.NewMinioConfig(
		minioCreds.GetEndpoint(),
		minioCreds.User,
		minioCreds.Password,
		minioCreds.Bucket,
	)
	require.NoError(s.T(), err)

	minioClient, err := minioclient.NewMinioClient(minioConfig)
	require.NoError(s.T(), err)

	// Create file repository
	s.fileRepo, err = filestorage.NewPgMinioStorage(ctx, s.db, minioClient)
	require.NoError(s.T(), err)

	// Create plant repository
	s.repo, err = plantstorage.NewPostgresPlantRepository(ctx, s.db)
	require.NoError(s.T(), err)
}

func (s *PlantRepositoryTestSuite) TearDownSuite() {
	ctx := context.Background()
	if s.minioContainer != nil {
		s.minioContainer.Terminate(ctx)
	}
	if s.dbContainer != nil {
		s.dbContainer.Terminate(ctx)
	}
	os.Chdir(s.prevDir)
}

func (s *PlantRepositoryTestSuite) pushTestPhoto(ctx context.Context) uuid.UUID {
	fileData := models.FileData{
		Name:        "test_photo.jpg",
		Reader:      bytes.NewReader([]byte("test photo content")),
		ContentType: "image/jpeg",
	}
	file, err := s.fileRepo.Upload(ctx, &fileData)
	require.NoError(s.T(), err)
	return file.ID
}

func (s *PlantRepositoryTestSuite) createTestPlant(ctx context.Context) *plant.Plant {
	// Upload main photo
	mainPhotoID := s.pushTestPhoto(ctx)
	additionalPhotoID := s.pushTestPhoto(ctx)

	// Create plant specification
	spec, err := plant.NewConiferousSpecification(
		1.5,                  // heightM
		0.5,                  // diameterM
		10,                   // soilAcidity
		plant.MediumMoisture, // soilMoisture
		plant.Light,          // lightRelation
		plant.MediumSoil,     // soilType
		10,                   // winterHardiness
	)
	require.NoError(s.T(), err)

	// Create plant photos collection
	photos := plant.NewPlantPhotos()
	photo, err := plant.NewPlantPhoto(additionalPhotoID, "Test photo")
	require.NoError(s.T(), err)
	photos.Add(photo)

	// Create plant
	plnt, err := plant.CreatePlant(
		uuid.New(),
		"Test Plant",
		"Testus Plantus",
		"Test description",
		mainPhotoID,
		*photos,
		plant.ConiferousCategory,
		spec,
		time.Now(),
		time.Now(),
	)
	require.NoError(s.T(), err)

	return plnt
}
