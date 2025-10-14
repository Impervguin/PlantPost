//go:build integration

package plantstorage_test

import (
	"bytes"
	"context"
	"fmt"
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
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

type PlantRepositoryTestSuite struct {
	suite.Suite
	dbContainer testcontainers.Container
	dbCreds     *pgtest.PostgresCredentials
	fileCnt     testcontainers.Container
	fileCreds   *miniotest.MinioCredentials
	db          *sqpgx.SquirrelPgx
	fileRepo    *filestorage.PgMinioStorage
	repo        *plantstorage.PostgresPlantRepository
	prevDir     string
}

func TestPlantRepositorySuite(t *testing.T) {
	suite.RunSuite(t, new(PlantRepositoryTestSuite))
}

func (s *PlantRepositoryTestSuite) BeforeEach(t provider.T) {
	t.Epic("Plant Repository")
	t.Feature("Plant Storage")
}

func (s *PlantRepositoryTestSuite) BeforeAll(t provider.T) {
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
	s.dbContainer = container

	if pgConfig.External {
		err = pgtest.CheckMigrationVersion(ctx, pgCreds, pgCreds.Database)
		require.NoError(t, err)
	} else {
		err = pgtest.Migrate(ctx, pgCreds, pgCreds.Database)
		require.NoError(t, err)
	}

	// Create database connection
	dbConfig := &sqpgx.SqpgxConfig{
		User:                   pgCreds.User,
		Password:               pgCreds.Password,
		DbName:                 pgCreds.Database,
		Host:                   pgCreds.Host,
		Port:                   pgCreds.Port,
		MaxConnections:         10,
		MaxConnectionsLifetime: time.Minute,
	}
	s.db, err = sqpgx.NewSquirrelPgx(ctx, dbConfig)
	require.NoError(t, err)

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
	s.fileRepo, err = filestorage.NewPgMinioStorage(ctx, s.db, minioCl)
	require.NoError(t, err)
	// Create plant repository
	s.repo, err = plantstorage.NewPostgresPlantRepository(ctx, s.db)
	require.NoError(t, err)
}

func (s *PlantRepositoryTestSuite) AfterAll(t provider.T) {
	ctx := context.Background()
	if s.fileCnt != nil {
		s.fileCnt.Terminate(ctx)
	}
	if s.dbContainer != nil {
		s.dbContainer.Terminate(ctx)
	}
	os.Chdir(s.prevDir)
}

func (s *PlantRepositoryTestSuite) AfterEach(t provider.T) {
	// err := pgtest.TruncateTables(context.Background(), s.dbCreds)
	// require.NoError(t, err)
	// err = miniotest.CleanUpBucket(context.Background(), s.fileCreds)
	// require.NoError(t, err)
}

func (s *PlantRepositoryTestSuite) pushTestPhoto(ctx context.Context, t provider.T) uuid.UUID {
	fileData := models.FileData{
		Name:        fmt.Sprintf("test-photo-%s.jpg", uuid.NewString()),
		Reader:      bytes.NewReader([]byte("test photo content")),
		ContentType: "image/jpeg",
	}
	file, err := s.fileRepo.Upload(ctx, &fileData)
	require.NoError(t, err)
	return file.ID
}

func (s *PlantRepositoryTestSuite) createTestPlant(ctx context.Context, t provider.T) *plant.Plant {
	// Upload main photo
	mainPhotoID := s.pushTestPhoto(ctx, t)
	additionalPhotoID := s.pushTestPhoto(ctx, t)

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

	// Create plant photos collection
	photos := plant.NewPlantPhotos()
	photo, err := plant.NewPlantPhoto(additionalPhotoID, "Test photo")
	require.NoError(t, err)
	photos.Add(photo)

	// Create plant
	plnt, err := plant.CreatePlant(
		uuid.New(),
		"Test Plant"+uuid.NewString(),
		"Testus Plantus"+uuid.NewString(),
		"Test description"+uuid.NewString(),
		mainPhotoID,
		*photos,
		plant.ConiferousCategory,
		spec,
		time.Now(),
		time.Now(),
	)
	require.NoError(t, err)

	return plnt
}

func (s *PlantRepositoryTestSuite) TestCreatePlant(t provider.T) {
	t.Tags("create", "plant")
	t.Description("Test plant creation functionality")

	t.Run("Successful plant creation", func(t provider.T) {
		ctx := context.Background()
		testPlant := s.createTestPlant(ctx, t)

		t.WithNewStep("Create plant", func(pctx provider.StepCtx) {
			createdPlant, err := s.repo.Create(ctx, testPlant)
			require.NoError(t, err)

			assert.Equal(t, testPlant.ID(), createdPlant.ID())
			assert.Equal(t, testPlant.GetName(), createdPlant.GetName())
			assert.Equal(t, testPlant.GetLatinName(), createdPlant.GetLatinName())
			assert.Equal(t, testPlant.GetDescription(), createdPlant.GetDescription())
			assert.Equal(t, testPlant.MainPhotoID(), createdPlant.MainPhotoID())
			assert.Equal(t, testPlant.GetCategory(), createdPlant.GetCategory())
			assert.Equal(t, testPlant.GetPhotos().Len(), createdPlant.GetPhotos().Len())
		})

		t.WithNewStep("Verify plant creation", func(pctx provider.StepCtx) {
			retrievedPlant, err := s.repo.Get(ctx, testPlant.ID())
			require.NoError(t, err)

			assert.Equal(t, testPlant.ID(), retrievedPlant.ID())
			assert.Equal(t, testPlant.GetName(), retrievedPlant.GetName())
			assert.Equal(t, testPlant.GetLatinName(), retrievedPlant.GetLatinName())
			assert.Equal(t, testPlant.GetDescription(), retrievedPlant.GetDescription())
			assert.Equal(t, testPlant.MainPhotoID(), retrievedPlant.MainPhotoID())
			assert.Equal(t, testPlant.GetCategory(), retrievedPlant.GetCategory())
			assert.Equal(t, testPlant.GetPhotos().Len(), retrievedPlant.GetPhotos().Len())
		})
	})

	t.Run("Plant with no photos", func(t provider.T) {
		ctx := context.Background()
		testPlant := s.createTestPlant(ctx, t)

		t.WithNewStep("Remove all photos", func(pctx provider.StepCtx) {
			photoIDs := make([]uuid.UUID, 0, testPlant.GetPhotos().Len())
			testPlant.GetPhotos().Iterate(func(e plant.PlantPhoto) error {
				photoIDs = append(photoIDs, e.ID())
				return nil
			})
			for _, photoID := range photoIDs {
				err := testPlant.DeletePhoto(photoID)
				require.NoError(t, err)
			}
			require.NoError(t, testPlant.Validate())
		})

		t.WithNewStep("Create plant without photos", func(pctx provider.StepCtx) {
			createdPlant, err := s.repo.Create(ctx, testPlant)
			require.NoError(t, err)
			assert.Equal(t, 0, createdPlant.GetPhotos().Len())
		})

		t.WithNewStep("Verify plant creation", func(pctx provider.StepCtx) {
			retrievedPlant, err := s.repo.Get(ctx, testPlant.ID())
			require.NoError(t, err)
			assert.Equal(t, 0, retrievedPlant.GetPhotos().Len())
		})
	})

	t.Run("Deciduous plant", func(t provider.T) {
		ctx := context.Background()
		testPlant := s.createTestPlant(ctx, t)

		t.WithNewStep("Create deciduous plant", func(pctx provider.StepCtx) {
			decSpec, err := plant.NewDeciduousSpecification(
				1.5,                  // heightM
				0.5,                  // diameterM
				plant.April,          // floweringPeriod
				10,                   // soilAcidity
				plant.MediumMoisture, // soilMoisture
				plant.Light,          // lightRelation
				plant.MediumSoil,     // soilType
				10,                   // winterHardiness
			)
			require.NoError(t, err)

			testPlant.UpdateSpec(decSpec)
			require.NoError(t, testPlant.Validate())
		})

		t.WithNewStep("Create plant in storage", func(pctx provider.StepCtx) {
			createdPlant, err := s.repo.Create(ctx, testPlant)
			require.NoError(t, err)
			assert.Equal(t, testPlant.ID(), createdPlant.ID())
		})

		t.WithNewStep("Verify plant creation", func(pctx provider.StepCtx) {
			retrievedPlant, err := s.repo.Get(ctx, testPlant.ID())
			require.NoError(t, err)
			assert.Equal(t, testPlant.ID(), retrievedPlant.ID())
			assert.Equal(t, testPlant.GetName(), retrievedPlant.GetName())
			assert.Equal(t, testPlant.GetLatinName(), retrievedPlant.GetLatinName())
			assert.Equal(t, testPlant.GetDescription(), retrievedPlant.GetDescription())
			assert.Equal(t, testPlant.MainPhotoID(), retrievedPlant.MainPhotoID())
			assert.Equal(t, testPlant.GetCategory(), retrievedPlant.GetCategory())
			assert.Equal(t, testPlant.GetPhotos().Len(), retrievedPlant.GetPhotos().Len())
		})
	})

	t.Run("Coniferous plant", func(t provider.T) {
		ctx := context.Background()
		testPlant := s.createTestPlant(ctx, t)

		t.WithNewStep("Create coniferous plant", func(pctx provider.StepCtx) {
			conSpec, err := plant.NewConiferousSpecification(
				1.5,                  // heightM
				0.5,                  // diameterM
				10,                   // soilAcidity
				plant.MediumMoisture, // soilMoisture
				plant.Light,          // lightRelation
				plant.MediumSoil,     // soilType
				10,                   // winterHardiness
			)
			require.NoError(t, err)

			testPlant.UpdateSpec(conSpec)
			require.NoError(t, testPlant.Validate())
		})

		t.WithNewStep("Create plant in storage", func(pctx provider.StepCtx) {
			createdPlant, err := s.repo.Create(ctx, testPlant)
			require.NoError(t, err)
			assert.Equal(t, testPlant.ID(), createdPlant.ID())
		})

		t.WithNewStep("Verify plant creation", func(pctx provider.StepCtx) {
			retrievedPlant, err := s.repo.Get(ctx, testPlant.ID())
			require.NoError(t, err)
			assert.Equal(t, testPlant.ID(), retrievedPlant.ID())
			assert.Equal(t, testPlant.GetName(), retrievedPlant.GetName())
			assert.Equal(t, testPlant.GetLatinName(), retrievedPlant.GetLatinName())
			assert.Equal(t, testPlant.GetDescription(), retrievedPlant.GetDescription())
			assert.Equal(t, testPlant.MainPhotoID(), retrievedPlant.MainPhotoID())
			assert.Equal(t, testPlant.GetCategory(), retrievedPlant.GetCategory())
			assert.Equal(t, testPlant.GetPhotos().Len(), retrievedPlant.GetPhotos().Len())
		})
	})

	t.Run("Duplicate plant", func(t provider.T) {
		ctx := context.Background()
		testPlant := s.createTestPlant(ctx, t)

		t.WithNewStep("Create plant in storage", func(pctx provider.StepCtx) {
			_, err := s.repo.Create(ctx, testPlant)
			require.NoError(t, err)
		})

		t.WithNewStep("Attempt to create duplicate plant", func(pctx provider.StepCtx) {
			_, err := s.repo.Create(ctx, testPlant)
			require.Error(t, err)
		})
	})
}

func (s *PlantRepositoryTestSuite) TestGetPlant(t provider.T) {
	t.Tags("get", "plant")
	t.Description("Test plant retrieval functionality")

	t.Run("Successful plant retrieval", func(t provider.T) {
		ctx := context.Background()
		testPlant := s.createTestPlant(ctx, t)

		t.WithNewStep("Create plant", func(pctx provider.StepCtx) {
			_, err := s.repo.Create(ctx, testPlant)
			require.NoError(t, err)
		})

		t.WithNewStep("Retrieve plant", func(pctx provider.StepCtx) {
			retrievedPlant, err := s.repo.Get(ctx, testPlant.ID())
			require.NoError(t, err)

			assert.Equal(t, testPlant.ID(), retrievedPlant.ID())
			assert.Equal(t, testPlant.GetName(), retrievedPlant.GetName())
			assert.Equal(t, testPlant.GetLatinName(), retrievedPlant.GetLatinName())
			assert.Equal(t, testPlant.GetDescription(), retrievedPlant.GetDescription())
			assert.Equal(t, testPlant.MainPhotoID(), retrievedPlant.MainPhotoID())
			assert.Equal(t, testPlant.GetCategory(), retrievedPlant.GetCategory())
			assert.Equal(t, testPlant.GetPhotos().Len(), retrievedPlant.GetPhotos().Len())
		})
	})

	t.Run("Non-existent plant", func(t provider.T) {
		ctx := context.Background()

		t.WithNewStep("Attempt to retrieve non-existent plant", func(pctx provider.StepCtx) {
			_, err := s.repo.Get(ctx, uuid.New())
			require.Error(t, err)
		})
	})
}

func (s *PlantRepositoryTestSuite) TestUpdatePlant(t provider.T) {
	t.Tags("update", "plant")
	t.Description("Test plant update functionality")

	t.Run("Successful plant update", func(t provider.T) {
		ctx := context.Background()
		testPlant := s.createTestPlant(ctx, t)

		t.WithNewStep("Create initial plant", func(pctx provider.StepCtx) {
			_, err := s.repo.Create(ctx, testPlant)
			require.NoError(t, err)
		})

		t.WithNewStep("Update plant", func(pctx provider.StepCtx) {
			newPhotoID := s.pushTestPhoto(ctx, t)

			updatedPlant, err := s.repo.Update(ctx, testPlant.ID(), func(p *plant.Plant) (*plant.Plant, error) {
				p.UpdateName("Updated Plant")
				p.UpdateLatinName("Updated Latin")
				p.UpdateDescription("Updated description")
				p.UpdateMainPhotoID(newPhotoID)

				newAdditionalPhotoID := s.pushTestPhoto(ctx, t)
				photo, err := plant.NewPlantPhoto(newAdditionalPhotoID, "New photo")
				require.NoError(t, err)
				p.AddPhoto(photo)

				return p, nil
			})
			require.NoError(t, err)

			assert.Equal(t, "Updated Plant", updatedPlant.GetName())
			assert.Equal(t, "Updated Latin", updatedPlant.GetLatinName())
			assert.Equal(t, "Updated description", updatedPlant.GetDescription())
			assert.Equal(t, newPhotoID, updatedPlant.MainPhotoID())
			assert.Equal(t, 2, updatedPlant.GetPhotos().Len())
			assert.True(t, updatedPlant.UpdatedAt().After(testPlant.UpdatedAt()))
		})
	})

	t.Run("Update plant specification", func(t provider.T) {
		ctx := context.Background()
		testPlant := s.createTestPlant(ctx, t)

		t.WithNewStep("Create initial plant", func(pctx provider.StepCtx) {
			_, err := s.repo.Create(ctx, testPlant)
			require.NoError(t, err)
		})

		var newSpec plant.PlantSpecification
		var err error
		t.WithNewStep("Update plant specification", func(pctx provider.StepCtx) {
			newSpec, err = plant.NewConiferousSpecification(
				1.5,                  // heightM
				0.5,                  // diameterM
				10,                   // soilAcidity
				plant.MediumMoisture, // soilMoisture
				plant.Light,          // lightRelation
				plant.MediumSoil,     // soilType
				10,                   // winterHardiness
			)
			require.NoError(t, err)

			updatedPlant, err := s.repo.Update(ctx, testPlant.ID(), func(p *plant.Plant) (*plant.Plant, error) {
				p.UpdateSpec(newSpec)
				return p, nil
			})
			require.NoError(t, err)

			assert.Equal(t, newSpec, updatedPlant.GetSpecification())
		})

		t.WithNewStep("Verify plant update", func(pctx provider.StepCtx) {
			retrievedPlant, err := s.repo.Get(ctx, testPlant.ID())
			require.NoError(t, err)
			assert.Equal(t, newSpec, retrievedPlant.GetSpecification())
		})
	})

	t.Run("Update non-existent plant", func(t provider.T) {
		ctx := context.Background()

		t.WithNewStep("Attempt to update non-existent plant", func(pctx provider.StepCtx) {
			_, err := s.repo.Update(ctx, uuid.New(), func(p *plant.Plant) (*plant.Plant, error) {
				p.UpdateName("Should fail")
				return p, nil
			})
			require.Error(t, err)
		})
	})
}

func (s *PlantRepositoryTestSuite) TestDeletePlant(t provider.T) {
	t.Tags("delete", "plant")
	t.Description("Test plant deletion functionality")

	t.Run("Successful plant deletion", func(t provider.T) {
		ctx := context.Background()
		testPlant := s.createTestPlant(ctx, t)

		t.WithNewStep("Create plant", func(pctx provider.StepCtx) {
			_, err := s.repo.Create(ctx, testPlant)
			require.NoError(t, err)
		})

		t.WithNewStep("Delete plant", func(pctx provider.StepCtx) {
			err := s.repo.Delete(ctx, testPlant.ID())
			require.NoError(t, err)
		})

		t.WithNewStep("Verify deletion", func(pctx provider.StepCtx) {
			_, err := s.repo.Get(ctx, testPlant.ID())
			require.Error(t, err)
		})
	})

	t.Run("Delete non-existent plant", func(t provider.T) {
		ctx := context.Background()

		t.WithNewStep("Delete non-existent plant", func(pctx provider.StepCtx) {
			err := s.repo.Delete(ctx, uuid.New())
			require.Error(t, err)
		})
	})
}
