//go:build integration

package albumstorage_test

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	minioclient "PlantSite/internal/infra/minio-client"
	"PlantSite/internal/infra/sqpgx"
	"PlantSite/internal/models"
	"PlantSite/internal/models/album"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/plant"
	filestorage "PlantSite/internal/repositories/pgminio/file-storage"
	albumstorage "PlantSite/internal/repositories/postgres/album-storage"
	authstorage "PlantSite/internal/repositories/postgres/auth-storage"
	plantstorage "PlantSite/internal/repositories/postgres/plant-storage"
	"PlantSite/internal/repositories/tests"
	"PlantSite/internal/testutils/miniotest"
	"PlantSite/internal/testutils/pgtest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type AlbumRepositoryTestSuite struct {
	suite.Suite
	container testcontainers.Container
	fileCnt   testcontainers.Container
	db        *sqpgx.SquirrelPgx
	albumRepo *albumstorage.PostgresAlbumRepository
	plantRepo *plantstorage.PostgresPlantRepository
	userRepo  *authstorage.PostgresAuthRepository
	fileRepo  *filestorage.PgMinioStorage
	prevDir   string
}

func TestAlbumRepositorySuite(t *testing.T) {
	suite.Run(t, new(AlbumRepositoryTestSuite))
}

func (s *AlbumRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	// Save current directory
	prevDir, err := os.Getwd()
	require.NoError(s.T(), err)
	s.prevDir = prevDir

	// Change directory to test working directory
	err = os.Chdir(tests.GetTestWorkingDir())
	require.NoError(s.T(), err)

	// Create new container for each test
	container, creds, err := pgtest.NewTestPostgres(ctx)
	require.NoError(s.T(), err)
	s.container = container

	// Run migrations
	err = pgtest.Migrate(ctx, &creds)
	require.NoError(s.T(), err)

	// Create database connection config
	config := &sqpgx.SqpgxConfig{
		User:                   creds.User,
		Password:               creds.Password,
		DbName:                 creds.Database,
		Host:                   creds.Host,
		Port:                   creds.Port,
		MaxConnections:         10,
		MaxConnectionsLifetime: time.Minute,
	}

	// Create database connection
	db, err := sqpgx.NewSquirrelPgx(ctx, config)
	require.NoError(s.T(), err)
	s.db = db

	// Create repositories
	s.albumRepo, err = albumstorage.NewPostgresAlbumRepository(ctx, db)
	require.NoError(s.T(), err)

	s.plantRepo, err = plantstorage.NewPostgresPlantRepository(ctx, db)
	require.NoError(s.T(), err)

	s.userRepo, err = authstorage.NewPostgresAuthRepository(ctx, db)
	require.NoError(s.T(), err)

	// Create file repository
	fileCnt, fileCreds, err := miniotest.NewTestMinio(ctx)
	require.NoError(s.T(), err)
	s.fileCnt = fileCnt
	err = miniotest.Migrate(context.Background(), fileCreds)
	require.NoError(s.T(), err)

	minioConfig, err := minioclient.NewMinioConfig(
		fileCreds.GetEndpoint(),
		fileCreds.User,
		fileCreds.Password,
		fileCreds.Bucket,
	)
	require.NoError(s.T(), err)

	minioCl, err := minioclient.NewMinioClient(
		minioConfig,
	)
	require.NoError(s.T(), err)
	s.fileRepo, err = filestorage.NewPgMinioStorage(ctx, db, minioCl)
	require.NoError(s.T(), err)
}

func (s *AlbumRepositoryTestSuite) TearDownSuite() {
	ctx := context.Background()
	if s.fileCnt != nil {
		s.fileCnt.Terminate(ctx)
	}
	if s.container != nil {
		s.container.Terminate(ctx)
	}
	err := os.Chdir(s.prevDir)
	require.NoError(s.T(), err)
}

func (s *AlbumRepositoryTestSuite) SetupTest()    {}
func (s *AlbumRepositoryTestSuite) TearDownTest() {}

func (s *AlbumRepositoryTestSuite) pushTestPlant() *plant.Plant {
	ctx := context.Background()

	plUUID := uuid.New()
	// Create file
	photo, err := s.fileRepo.Upload(ctx, &models.FileData{
		Name:        plUUID.String(),
		Reader:      bytes.NewReader([]byte("test")),
		ContentType: "image/jpeg",
	})
	require.NoError(s.T(), err)

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

	// Create plant
	plnt, err := plant.CreatePlant(
		plUUID,
		plUUID.String(),
		"Testus plantus",
		"Test description",
		photo.ID, // main photo ID
		*plant.NewPlantPhotos(),
		plant.ConiferousCategory,
		spec,
		time.Now(),
		time.Now(),
	)
	require.NoError(s.T(), err)

	// Store in database
	createdPlant, err := s.plantRepo.Create(ctx, plnt)
	require.NoError(s.T(), err)

	return createdPlant
}

func (s *AlbumRepositoryTestSuite) createTestMember() *auth.Member {
	memId := uuid.New()

	user, err := auth.CreateMember(
		memId,
		memId.String()[:8],
		memId.String()+"@test.com",
		[]byte("test"),
		time.Now(),
	)
	require.NoError(s.T(), err)

	return user
}

func (s *AlbumRepositoryTestSuite) pushTestUser() auth.User {
	ctx := context.Background()

	user := s.createTestMember()

	_, err := s.userRepo.Create(ctx, user)
	require.NoError(s.T(), err)

	return user
}

func (s *AlbumRepositoryTestSuite) createTestAlbum(plantIDs uuid.UUIDs, ownerID uuid.UUID) *album.Album {
	albID := uuid.New()
	alb, err := album.CreateAlbum(
		albID,
		albID.String(),
		"Test Description",
		plantIDs,
		ownerID, // owner ID
		time.Now(),
		time.Now(),
	)
	require.NoError(s.T(), err)

	return alb
}
