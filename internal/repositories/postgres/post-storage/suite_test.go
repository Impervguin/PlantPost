//go:build integration

package poststorage_test

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	minioclient "PlantSite/internal/infra/minio-client"
	"PlantSite/internal/infra/sqpgx"
	"PlantSite/internal/models"
	"PlantSite/internal/models/auth"
	"PlantSite/internal/models/post"
	filestorage "PlantSite/internal/repositories/pgminio/file-storage"
	authstorage "PlantSite/internal/repositories/postgres/auth-storage"
	poststorage "PlantSite/internal/repositories/postgres/post-storage"
	"PlantSite/internal/repositories/tests"
	"PlantSite/internal/testutils/miniotest"
	"PlantSite/internal/testutils/pgtest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type PostRepositoryTestSuite struct {
	suite.Suite
	dbContainer    testcontainers.Container
	minioContainer testcontainers.Container
	db             *sqpgx.SquirrelPgx
	fileRepo       *filestorage.PgMinioStorage
	repo           *poststorage.PostgresPostRepository
	userRepo       *authstorage.PostgresAuthRepository
	prevDir        string
}

func TestPostRepositorySuite(t *testing.T) {
	suite.Run(t, new(PostRepositoryTestSuite))
}

func (s *PostRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	// Save current directory
	prevDir, err := os.Getwd()
	require.NoError(s.T(), err)
	s.prevDir = prevDir

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

	// Run minio migrations
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

	// Create post repository
	s.repo, err = poststorage.NewPostgresPostRepository(ctx, s.db)
	require.NoError(s.T(), err)

	// Create user repository
	s.userRepo, err = authstorage.NewPostgresAuthRepository(ctx, s.db)
	require.NoError(s.T(), err)
}

func (s *PostRepositoryTestSuite) TearDownSuite() {
	ctx := context.Background()
	if s.minioContainer != nil {
		s.minioContainer.Terminate(ctx)
	}
	if s.dbContainer != nil {
		s.dbContainer.Terminate(ctx)
	}
	// Restore original directory
	os.Chdir(s.prevDir)
}

func (s *PostRepositoryTestSuite) uploadTestPhoto(ctx context.Context) uuid.UUID {
	id := uuid.New()
	fileData := models.FileData{
		Name:        id.String() + ".jpg",
		Reader:      bytes.NewReader([]byte("test photo content")),
		ContentType: "image/jpeg",
	}
	file, err := s.fileRepo.UploadFile(ctx, fileData)
	require.NoError(s.T(), err)
	return file.ID
}

func (s *PostRepositoryTestSuite) pushAuthor(ctx context.Context) uuid.UUID {
	nameUU := uuid.New()
	member, err := auth.NewMember(
		nameUU.String()[:8],
		nameUU.String()[:8]+"@test.com",
		[]byte("hassPasword"),
	)
	require.NoError(s.T(), err)
	_, err = s.userRepo.Create(ctx, member)
	require.NoError(s.T(), err)
	_, err = s.userRepo.Update(ctx, member.ID(), func(u auth.User) (auth.User, error) {
		member := u.(*auth.Member)
		return auth.CreateAuthor(*member, time.Now(), true, time.Now().Add(-time.Hour))
	})
	require.NoError(s.T(), err)
	return member.ID()
}

func (s *PostRepositoryTestSuite) createTestPost(ctx context.Context) *post.Post {
	ownerID := s.pushAuthor(ctx)
	return s.createAuthorPost(ctx, ownerID)
}

func (s *PostRepositoryTestSuite) createAuthorPost(ctx context.Context, ownerID uuid.UUID) *post.Post {
	photoIDs := make([]uuid.UUID, 0, 2)
	for i := 0; i < 2; i++ {
		photoID := s.uploadTestPhoto(ctx)
		photoIDs = append(photoIDs, photoID)
	}

	// Create post content
	content, err := post.NewContent("Test post content", post.ContentTypePlainText)
	require.NoError(s.T(), err)

	// Create post photos
	photos := post.NewPostPhotos()
	for i, photoID := range photoIDs {
		photo, err := post.CreatePostPhoto(uuid.New(), photoID, i)
		require.NoError(s.T(), err)
		err = photos.Add(photo)
		require.NoError(s.T(), err)
	}

	// Create post
	pst, err := post.CreatePost(
		uuid.New(),
		"Test Post",
		*content,
		[]string{"test", "post"},
		ownerID, // author ID
		*photos,
		time.Now(),
		time.Now(),
	)
	require.NoError(s.T(), err)

	return pst
}
