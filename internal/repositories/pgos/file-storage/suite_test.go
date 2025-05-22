//go:build integration

package filestorage_test

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	filedir "PlantSite/internal/infra/os/file-dir"
	"PlantSite/internal/infra/sqpgx"
	"PlantSite/internal/models"
	filestorage "PlantSite/internal/repositories/pgos/file-storage"
	"PlantSite/internal/repositories/tests"
	"PlantSite/internal/testutils/pgtest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type FileStorageTestSuite struct {
	suite.Suite
	dbContainer    testcontainers.Container
	minioContainer testcontainers.Container
	db             *sqpgx.SquirrelPgx
	fileClient     *filedir.FileClient
	storage        *filestorage.PgOsFileStorage
	prevDir        string
	testBucketName string
}

func TestFileStorageSuite(t *testing.T) {
	suite.Run(t, new(FileStorageTestSuite))
}

func (s *FileStorageTestSuite) SetupSuite() {
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

	s.testBucketName = "bucket"

	// Setup dir in os
	err = os.MkdirAll("test-dir/"+s.testBucketName, 0755)
	require.NoError(s.T(), err)

	// Create file client
	s.fileClient, err = filedir.NewFileClient("test-dir")
	require.NoError(s.T(), err)

	// Create storage instance
	s.storage = filestorage.NewPgOsFileStorage(s.testBucketName, s.fileClient, s.db)
	require.NotNil(s.T(), s.storage)
}

func (s *FileStorageTestSuite) TearDownSuite() {
	ctx := context.Background()
	if s.dbContainer != nil {
		s.dbContainer.Terminate(ctx)
	}
	os.RemoveAll("test-dir")
	// Restore original directory
	os.Chdir(s.prevDir)
}

func createTestFileData() models.FileData {
	id := uuid.New()
	return models.FileData{
		Name:        id.String() + ".txt",
		Reader:      bytes.NewReader([]byte("test content")),
		ContentType: "text/plain; charset=utf-8",
	}
}
