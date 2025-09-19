//go:build integration

package albumstorage_test

import (
	"bytes"
	"context"
	"os"
	"slices"
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
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

type AlbumRepositoryTestSuite struct {
	suite.Suite
	container testcontainers.Container
	fileCnt   testcontainers.Container
	fileCreds *miniotest.MinioCredentials
	db        *sqpgx.SquirrelPgx
	albumRepo *albumstorage.PostgresAlbumRepository
	plantRepo *plantstorage.PostgresPlantRepository
	userRepo  *authstorage.PostgresAuthRepository
	fileRepo  *filestorage.PgMinioStorage
	prevDir   string
	dbCreds   *pgtest.PostgresCredentials
}

func TestAlbumRepositorySuite(t *testing.T) {
	suite.RunSuite(t, new(AlbumRepositoryTestSuite))
}

func (s *AlbumRepositoryTestSuite) BeforeEach(t provider.T) {
	t.Epic("Album Repository")
	t.Feature("Album Storage")
}

func (s *AlbumRepositoryTestSuite) BeforeAll(t provider.T) {
	ctx := context.Background()

	// Save current directory
	prevDir, err := os.Getwd()
	require.NoError(t, err)
	s.prevDir = prevDir

	// Change directory to test working directory
	err = os.Chdir(tests.GetTestWorkingDir())
	require.NoError(t, err)

	pgConfig := pgtest.GetConfig()
	var pgCreds *pgtest.PostgresCredentials
	var container testcontainers.Container
	if pgConfig.External {
		pgCreds, err = pgtest.NewTestExternalPostgres(ctx, pgConfig)
		require.NoError(t, err)
	} else {
		container, pgCreds, err = pgtest.NewTestPostgres(ctx)
		require.NoError(t, err)
	}
	s.dbCreds = pgCreds
	s.container = container

	if pgConfig.External {
		err = pgtest.CheckMigrationVersion(ctx, pgCreds, pgCreds.Database)
		require.NoError(t, err)
	} else {
		err = pgtest.Migrate(ctx, pgCreds, pgCreds.Database)
		require.NoError(t, err)
	}

	// Create database connection config
	config := &sqpgx.SqpgxConfig{
		User:                   pgCreds.User,
		Password:               pgCreds.Password,
		DbName:                 pgCreds.Database,
		Host:                   pgCreds.Host,
		Port:                   pgCreds.Port,
		MaxConnections:         10,
		MaxConnectionsLifetime: time.Minute,
	}

	// Create database connection
	db, err := sqpgx.NewSquirrelPgx(ctx, config)
	require.NoError(t, err)
	s.db = db

	// Create repositories
	s.albumRepo, err = albumstorage.NewPostgresAlbumRepository(ctx, db)
	require.NoError(t, err)

	s.plantRepo, err = plantstorage.NewPostgresPlantRepository(ctx, db)
	require.NoError(t, err)

	s.userRepo, err = authstorage.NewPostgresAuthRepository(ctx, db)
	require.NoError(t, err)

	// Create file repository
	var fileCnt testcontainers.Container
	var fileCreds *miniotest.MinioCredentials
	minConfig := miniotest.GetConfig()
	if minConfig.External {
		fileCreds, err = miniotest.NewTestExternalMinio(ctx, minConfig)
		require.NoError(t, err)
	} else {
		fileCnt, fileCreds, err = miniotest.NewTestMinio(ctx)
		require.NoError(t, err)
	}
	s.fileCnt = fileCnt
	s.fileCreds = fileCreds

	err = miniotest.CheckBucketExists(context.Background(), fileCreds)
	if err != nil {
		require.ErrorIs(t, err, miniotest.BucketDoesNotExistError)
		err = miniotest.Migrate(context.Background(), fileCreds)
		require.NoError(t, err)
	} else {
		require.NoError(t, err)
	}

	minioConfig, err := minioclient.NewMinioConfig(
		fileCreds.GetEndpoint(),
		fileCreds.User,
		fileCreds.Password,
		fileCreds.Bucket,
	)
	require.NoError(t, err)

	minioCl, err := minioclient.NewMinioClient(
		minioConfig,
	)
	require.NoError(t, err)
	s.fileRepo, err = filestorage.NewPgMinioStorage(ctx, db, minioCl)
	require.NoError(t, err)
}

func (s *AlbumRepositoryTestSuite) AfterAll(t provider.T) {
	ctx := context.Background()
	if s.fileCnt != nil {
		s.fileCnt.Terminate(ctx)
	}
	if s.container != nil {
		s.container.Terminate(ctx)
	}
	err := os.Chdir(s.prevDir)
	require.NoError(t, err)
}

func (s *AlbumRepositoryTestSuite) AfterEach(t provider.T) {
	// err := miniotest.CleanUpBucket(context.Background(), s.fileCreds)
	// require.NoError(t, err)
	// err = pgtest.TruncateTables(context.Background(), s.dbCreds)
	// require.NoError(t, err)
}

func (s *AlbumRepositoryTestSuite) pushTestPlant(t provider.T) *plant.Plant {
	ctx := context.Background()

	plUUID := uuid.New()
	// Create file
	photo, err := s.fileRepo.Upload(ctx, &models.FileData{
		Name:        plUUID.String(),
		Reader:      bytes.NewReader([]byte("test")),
		ContentType: "image/jpeg",
	})
	require.NoError(t, err)

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
	require.NoError(t, err)

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
	require.NoError(t, err)

	// Store in database
	createdPlant, err := s.plantRepo.Create(ctx, plnt)
	require.NoError(t, err)

	return createdPlant
}

func (s *AlbumRepositoryTestSuite) createTestMember(t provider.T) *auth.Member {
	memId := uuid.New()

	user, err := auth.CreateMember(
		memId,
		memId.String()[:8],
		memId.String()+"@test.com",
		[]byte("test"),
		time.Now(),
	)
	require.NoError(t, err)

	return user
}

func (s *AlbumRepositoryTestSuite) pushTestUser(t provider.T) auth.User {
	ctx := context.Background()

	user := s.createTestMember(t)

	_, err := s.userRepo.Create(ctx, user)
	require.NoError(t, err)

	return user
}

func (s *AlbumRepositoryTestSuite) createTestAlbum(t provider.T, plantIDs uuid.UUIDs, ownerID uuid.UUID) *album.Album {
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
	require.NoError(t, err)

	return alb
}

func (s *AlbumRepositoryTestSuite) TestCreateAlbum(t provider.T) {
	t.Tags("create", "album")
	t.Description("Test album creation functionality")

	ctx := context.Background()
	t.Run("Successful album creation", func(t provider.T) {
		t.WithNewStep("Prepare test data", func(pctx provider.StepCtx) {
			plant := s.pushTestPlant(t)
			user := s.pushTestUser(t)

			album := s.createTestAlbum(t, uuid.UUIDs{plant.ID()}, user.ID())

			t.WithNewStep("Create album", func(pctx provider.StepCtx) {
				createdAlbum, err := s.albumRepo.Create(ctx, album)
				require.NoError(t, err)

				assert.Equal(t, album.ID(), createdAlbum.ID())
				assert.Equal(t, album.Name(), createdAlbum.Name())
				assert.Equal(t, album.Description(), createdAlbum.Description())
				assert.Equal(t, album.PlantIDs(), createdAlbum.PlantIDs())
				assert.Equal(t, album.GetOwnerID(), createdAlbum.GetOwnerID())
			})

			t.WithNewStep("Verify album persistence", func(pctx provider.StepCtx) {
				fetchedAlbum, err := s.albumRepo.Get(ctx, album.ID())
				require.NoError(t, err)

				assert.Equal(t, album.ID(), fetchedAlbum.ID())
				assert.Equal(t, album.Name(), fetchedAlbum.Name())
				assert.Equal(t, album.Description(), fetchedAlbum.Description())
				assert.Equal(t, album.PlantIDs(), fetchedAlbum.PlantIDs())
				assert.Equal(t, album.GetOwnerID(), fetchedAlbum.GetOwnerID())
			})
		})
	})

	t.Run("No owner trying to create album", func(t provider.T) {
		t.WithNewStep("Prepare test data", func(pctx provider.StepCtx) {
			plant := s.pushTestPlant(t)
			album := s.createTestAlbum(t, uuid.UUIDs{plant.ID()}, uuid.New())

			t.WithNewStep("Attempt to create album", func(pctx provider.StepCtx) {
				_, err := s.albumRepo.Create(ctx, album)
				require.Error(t, err)
			})

			t.WithNewStep("Verify no album creation", func(pctx provider.StepCtx) {
				_, err := s.albumRepo.Get(ctx, album.ID())
				require.Error(t, err)
			})
		})
	})
}

func (s *AlbumRepositoryTestSuite) TestGetAlbum(t provider.T) {
	t.Tags("get", "album")
	t.Description("Test album retrieval functionality")

	ctx := context.Background()
	t.Run("Successful album retrieval", func(t provider.T) {
		t.WithNewStep("Prepare test data", func(pctx provider.StepCtx) {
			plant := s.pushTestPlant(t)
			user := s.pushTestUser(t)
			album := s.createTestAlbum(t, uuid.UUIDs{plant.ID()}, user.ID())

			t.WithNewStep("Create album", func(pctx provider.StepCtx) {
				_, err := s.albumRepo.Create(ctx, album)
				require.NoError(t, err)
			})

			t.WithNewStep("Retrieve album", func(pctx provider.StepCtx) {
				retrievedAlbum, err := s.albumRepo.Get(ctx, album.ID())
				require.NoError(t, err)

				assert.Equal(t, album.ID(), retrievedAlbum.ID())
				assert.Equal(t, album.Name(), retrievedAlbum.Name())
				assert.Equal(t, album.Description(), retrievedAlbum.Description())
				assert.Equal(t, album.PlantIDs(), retrievedAlbum.PlantIDs())
				assert.Equal(t, album.GetOwnerID(), retrievedAlbum.GetOwnerID())
			})
		})
	})

	t.Run("Non-existent album retrieval", func(t provider.T) {
		t.WithNewStep("Prepare test data", func(pctx provider.StepCtx) {
			plant := s.pushTestPlant(t)
			user := s.pushTestUser(t)
			album := s.createTestAlbum(t, uuid.UUIDs{plant.ID()}, user.ID())

			t.WithNewStep("Attempt to retrieve non-existent album", func(pctx provider.StepCtx) {
				_, err := s.albumRepo.Get(ctx, uuid.New())
				require.Error(t, err)
			})

			t.WithNewStep("Verify no album retrieval", func(pctx provider.StepCtx) {
				_, err := s.albumRepo.Get(ctx, album.ID())
				require.Error(t, err)
			})
		})
	})
}

func (s *AlbumRepositoryTestSuite) TestUpdateAlbum(t provider.T) {
	t.Tags("update", "album")
	t.Description("Test album update functionality")

	ctx := context.Background()
	t.Run("Successful album update", func(t provider.T) {
		t.WithNewStep("Prepare test data", func(pctx provider.StepCtx) {
			plant1 := s.pushTestPlant(t)
			plant2 := s.pushTestPlant(t)
			user := s.pushTestUser(t)
			talbum := s.createTestAlbum(t, uuid.UUIDs{plant1.ID()}, user.ID())

			t.WithNewStep("Create album", func(pctx provider.StepCtx) {
				_, err := s.albumRepo.Create(ctx, talbum)
				require.NoError(t, err)
			})

			t.WithNewStep("Update album", func(pctx provider.StepCtx) {
				newName := "Updated Album Name"
				newDescription := "Updated description"
				newPlantIDs := uuid.UUIDs{plant2.ID()}

				updatedAlbum, err := s.albumRepo.Update(ctx, talbum.ID(), func(a *album.Album) (*album.Album, error) {
					a.UpdateName(newName)
					a.UpdateDescription(newDescription)
					for _, plantID := range talbum.PlantIDs() {
						if !slices.Contains(newPlantIDs, plantID) {
							a.RemovePlant(plantID)
						}
					}
					for _, plantID := range newPlantIDs {
						if !slices.Contains(talbum.PlantIDs(), plantID) {
							a.AddPlant(plantID)
						}
					}
					return a, nil
				})
				require.NoError(t, err)

				assert.Equal(t, newName, updatedAlbum.Name())
				assert.Equal(t, newDescription, updatedAlbum.Description())
				assert.Equal(t, newPlantIDs, updatedAlbum.PlantIDs())
			})
		})
	})

	t.Run("Non-existent album update", func(t provider.T) {
		t.WithNewStep("Prepare test data", func(pctx provider.StepCtx) {
			plant := s.pushTestPlant(t)
			user := s.pushTestUser(t)
			talbum := s.createTestAlbum(t, uuid.UUIDs{plant.ID()}, user.ID())

			t.WithNewStep("Create album", func(pctx provider.StepCtx) {
				_, err := s.albumRepo.Create(ctx, talbum)
				require.NoError(t, err)
			})

			t.WithNewStep("Attempt to update non-existent album", func(pctx provider.StepCtx) {
				_, err := s.albumRepo.Update(ctx, uuid.New(), func(a *album.Album) (*album.Album, error) {
					a.UpdateName("Should fail")
					return a, nil
				})
				require.Error(t, err)
			})

			t.WithNewStep("Verify no album update", func(pctx provider.StepCtx) {
				fetchedAlbum, err := s.albumRepo.Get(ctx, talbum.ID())
				require.NoError(t, err)

				assert.Equal(t, talbum.ID(), fetchedAlbum.ID())
				assert.Equal(t, talbum.Name(), fetchedAlbum.Name())
				assert.Equal(t, talbum.Description(), fetchedAlbum.Description())
				assert.Equal(t, talbum.PlantIDs(), fetchedAlbum.PlantIDs())
				assert.Equal(t, talbum.GetOwnerID(), fetchedAlbum.GetOwnerID())
			})
		})
	})
}

func (s *AlbumRepositoryTestSuite) TestDeleteAlbum(t provider.T) {
	t.Tags("delete", "album")
	t.Description("Test album deletion functionality")

	ctx := context.Background()
	t.Run("Successful album deletion", func(t provider.T) {
		t.WithNewStep("Prepare test data", func(pctx provider.StepCtx) {
			plant := s.pushTestPlant(t)
			user := s.pushTestUser(t)
			album := s.createTestAlbum(t, uuid.UUIDs{plant.ID()}, user.ID())

			t.WithNewStep("Create album", func(pctx provider.StepCtx) {
				_, err := s.albumRepo.Create(ctx, album)
				require.NoError(t, err)
			})

			t.WithNewStep("Delete album", func(pctx provider.StepCtx) {
				err := s.albumRepo.Delete(ctx, album.ID())
				require.NoError(t, err)
			})

			t.WithNewStep("Verify deletion", func(pctx provider.StepCtx) {
				_, err := s.albumRepo.Get(ctx, album.ID())
				require.Error(t, err)
			})
		})
	})

	t.Run("Non-existent album deletion", func(t provider.T) {
		t.WithNewStep("Prepare test data", func(pctx provider.StepCtx) {
			plant := s.pushTestPlant(t)
			user := s.pushTestUser(t)
			album := s.createTestAlbum(t, uuid.UUIDs{plant.ID()}, user.ID())

			t.WithNewStep("Create album", func(pctx provider.StepCtx) {
				_, err := s.albumRepo.Create(ctx, album)
				require.NoError(t, err)
			})

			t.WithNewStep("Attempt to delete non-existent album", func(pctx provider.StepCtx) {
				err := s.albumRepo.Delete(ctx, uuid.New())
				require.Error(t, err)
			})

			t.WithNewStep("Verify no album deletion", func(pctx provider.StepCtx) {
				_, err := s.albumRepo.Get(ctx, album.ID())
				require.NoError(t, err)
			})
		})
	})
}

func (s *AlbumRepositoryTestSuite) TestGetAlbumsByOwner(t provider.T) {
	t.Tags("get", "albums_by_owner")
	t.Description("Test retrieving albums by owner functionality")

	ctx := context.Background()
	t.Run("Successful album retrieval by owner", func(t provider.T) {
		t.WithNewStep("Prepare test data", func(pctx provider.StepCtx) {
			plant := s.pushTestPlant(t)
			user1 := s.pushTestUser(t)
			user2 := s.pushTestUser(t)

			album1 := s.createTestAlbum(t, uuid.UUIDs{plant.ID()}, user1.ID())
			album2 := s.createTestAlbum(t, uuid.UUIDs{plant.ID()}, user1.ID())
			album3 := s.createTestAlbum(t, uuid.UUIDs{plant.ID()}, user2.ID())

			t.WithNewStep("Create albums", func(pctx provider.StepCtx) {
				_, err := s.albumRepo.Create(ctx, album1)
				require.NoError(t, err)
				_, err = s.albumRepo.Create(ctx, album2)
				require.NoError(t, err)
				_, err = s.albumRepo.Create(ctx, album3)
				require.NoError(t, err)
			})

			t.WithNewStep("Get albums by owner", func(pctx provider.StepCtx) {
				albums, err := s.albumRepo.List(ctx, user1.ID())
				require.NoError(t, err)

				assert.Len(t, albums, 2)
				albumIDs := make([]uuid.UUID, len(albums))
				for i, alb := range albums {
					albumIDs[i] = alb.ID()
				}
				assert.Contains(t, albumIDs, album1.ID())
				assert.Contains(t, albumIDs, album2.ID())
				assert.NotContains(t, albumIDs, album3.ID())
			})
		})
	})
}
