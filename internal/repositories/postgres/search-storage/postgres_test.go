//go:build integration

package searchstorage_test

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
	"PlantSite/internal/models/plant"
	"PlantSite/internal/models/post"
	"PlantSite/internal/models/search"
	filestorage "PlantSite/internal/repositories/pgminio/file-storage"
	authstorage "PlantSite/internal/repositories/postgres/auth-storage"
	plantstorage "PlantSite/internal/repositories/postgres/plant-storage"
	poststorage "PlantSite/internal/repositories/postgres/post-storage"
	searchstorage "PlantSite/internal/repositories/postgres/search-storage"
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

type SearchRepositoryTestSuite struct {
	suite.Suite
	dbContainer testcontainers.Container
	dbCreds     *pgtest.PostgresCredentials
	fileCnt     testcontainers.Container
	fileCreds   *miniotest.MinioCredentials
	db          *sqpgx.SquirrelPgx
	fileRepo    *filestorage.PgMinioStorage
	plantRepo   *plantstorage.PostgresPlantRepository
	postRepo    *poststorage.PostgresPostRepository
	searchRepo  *searchstorage.PostgresSearchRepository
	userRepo    *authstorage.PostgresAuthRepository
	prevDir     string
}

func TestSearchRepositorySuite(t *testing.T) {
	suite.RunSuite(t, new(SearchRepositoryTestSuite))
}

func (s *SearchRepositoryTestSuite) BeforeEach(t provider.T) {
	t.Epic("Search Repository")
	t.Feature("Search Storage")
}

func (s *SearchRepositoryTestSuite) BeforeAll(t provider.T) {
	ctx := context.Background()

	// Save current directory
	prevDir, err := os.Getwd()
	require.NoError(t, err)
	s.prevDir = prevDir

	os.Chdir(tests.GetTestWorkingDir())

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
		DBName:                 pgCreds.Database,
		Host:                   pgCreds.Host,
		Port:                   pgCreds.Port,
		MaxConnections:         10,
		MaxConnectionsLifetime: time.Minute,
	}
	s.db, err = sqpgx.NewSquirrelPgx(ctx, dbConfig)
	require.NoError(t, err)

	// Setup MinIO container
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
		require.ErrorIs(t, err, miniotest.ErrBucketDoesNotExist)
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
	// Create repositories
	s.searchRepo, err = searchstorage.NewPostgresSearchRepository(ctx, s.db)
	require.NoError(t, err)

	s.plantRepo, err = plantstorage.NewPostgresPlantRepository(ctx, s.db)
	require.NoError(t, err)

	plntGetter := searchstorage.NewSearchPlantGetter(s.searchRepo)

	s.postRepo, err = poststorage.NewPostgresPostRepository(ctx, s.db, plntGetter)
	require.NoError(t, err)

	s.userRepo, err = authstorage.NewPostgresAuthRepository(ctx, s.db)
	require.NoError(t, err)
}

func (s *SearchRepositoryTestSuite) AfterAll(t provider.T) {
	ctx := context.Background()
	if s.fileCnt != nil {
		s.fileCnt.Terminate(ctx)
	}
	if s.dbContainer != nil {
		s.dbContainer.Terminate(ctx)
	}
	// Restore original directory
	os.Chdir(s.prevDir)
}

func (s *SearchRepositoryTestSuite) AfterEach(t provider.T) {
	// err := pgtest.TruncateTables(context.Background(), s.dbCreds)
	// require.NoError(t, err)
	// err = miniotest.CleanUpBucket(context.Background(), s.fileCreds)
	// require.NoError(t, err)
}

func (s *SearchRepositoryTestSuite) uploadTestPhoto(ctx context.Context, t provider.T) uuid.UUID {
	nameUU := uuid.New()
	fileData := models.FileData{
		Name:        nameUU.String() + ".jpg",
		Reader:      bytes.NewReader([]byte("test photo content")),
		ContentType: "image/jpeg",
	}
	file, err := s.fileRepo.Upload(ctx, &fileData)
	require.NoError(t, err)
	return file.ID
}

func (s *SearchRepositoryTestSuite) pushAuthor(ctx context.Context, t provider.T) uuid.UUID {
	nameUU := uuid.New()
	member, err := auth.NewMember(
		nameUU.String()[:8],
		nameUU.String()[:8]+"@test.com",
		[]byte("hassPasword"),
	)
	require.NoError(t, err)
	_, err = s.userRepo.Create(ctx, member)
	require.NoError(t, err)
	_, err = s.userRepo.Update(ctx, member.ID(), func(u auth.User) (auth.User, error) {
		member := u.(*auth.Member)
		return auth.CreateAuthor(*member, time.Now(), true, time.Now().Add(-time.Hour))
	})
	require.NoError(t, err)
	return member.ID()
}

func (s *SearchRepositoryTestSuite) createTestPost(ctx context.Context, t provider.T) *post.Post {
	ownerID := s.pushAuthor(ctx, t)
	return s.createAuthorPost(ctx, t, ownerID)
}

func (s *SearchRepositoryTestSuite) createAuthorPost(ctx context.Context, t provider.T, ownerID uuid.UUID) *post.Post {
	photoIDs := make([]uuid.UUID, 0, 2)
	for i := 0; i < 2; i++ {
		photoID := s.uploadTestPhoto(ctx, t)
		photoIDs = append(photoIDs, photoID)
	}

	// Create post content
	content, err := post.NewContent("Test post content", post.ContentTypePlainText)
	require.NoError(t, err)

	// Create post photos
	photos := post.NewPostPhotos()
	for i, photoID := range photoIDs {
		photo, err := post.CreatePostPhoto(uuid.New(), photoID, i)
		require.NoError(t, err)
		err = photos.Add(photo)
		require.NoError(t, err)
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
	require.NoError(t, err)

	return pst
}

func (s *SearchRepositoryTestSuite) createConiferousPlant(
	ctx context.Context,
	t provider.T,
	name string,
	height, diameter float64,
	moisture plant.SoilMoisture,
	acid plant.SoilAcidity,
	light plant.LightRelation,
	soilType plant.Soil,
	winterHardiness plant.WinterHardiness,
) *plant.Plant {
	// Upload main photo
	mainPhotoID := s.uploadTestPhoto(ctx, t)

	// Create coniferous specification
	spec, err := plant.NewConiferousSpecification(
		height,          // heightM
		diameter,        // diameterM
		acid,            // soilAcidity
		moisture,        // soilMoisture
		light,           // lightRelation
		soilType,        // soilType
		winterHardiness, // winterHardiness
	)
	require.NoError(t, err)

	// Create plant photos collection
	photos := plant.NewPlantPhotos()
	additionalPhotoID := s.uploadTestPhoto(ctx, t)

	photo, err := plant.CreatePlantPhoto(uuid.New(), additionalPhotoID, "Test photo")
	require.NoError(t, err)
	err = photos.Add(photo)
	require.NoError(t, err)

	// Create plant
	plnt, err := plant.CreatePlant(
		uuid.New(),
		name,
		"Testus Plantus",
		"Test description",
		mainPhotoID,
		*photos,
		spec.Category(),
		spec,
		time.Now(),
		time.Now(),
	)
	require.NoError(t, err)

	return plnt
}

func (s *SearchRepositoryTestSuite) createDeciduousPlant(
	ctx context.Context,
	t provider.T,
	name string,
	height, diameter float64,
	moisture plant.SoilMoisture,
	acid plant.SoilAcidity,
	light plant.LightRelation,
	soilType plant.Soil,
	winterHardiness plant.WinterHardiness,
	flowering plant.FloweringPeriod,
) *plant.Plant {
	// Upload main photo
	mainPhotoID := s.uploadTestPhoto(ctx, t)

	// Create deciduous specification
	spec, err := plant.NewDeciduousSpecification(
		height,          // heightM
		diameter,        // diameterM
		flowering,       // floweringPeriod
		acid,            // soilAcidity
		moisture,        // soilMoisture
		light,           // lightRelation
		soilType,        // soilType
		winterHardiness, // winterHardiness
	)
	require.NoError(t, err)

	// Create plant photos collection
	photos := plant.NewPlantPhotos()
	additionalPhotoID := s.uploadTestPhoto(ctx, t)

	photo, err := plant.CreatePlantPhoto(uuid.New(), additionalPhotoID, "Test photo")
	require.NoError(t, err)
	err = photos.Add(photo)
	require.NoError(t, err)

	// Create plant
	plnt, err := plant.CreatePlant(
		uuid.New(),
		name,
		"Testus Plantus",
		"Test description",
		mainPhotoID,
		*photos,
		spec.Category(),
		spec,
		time.Now(),
		time.Now(),
	)
	require.NoError(t, err)

	return plnt
}

func (s *SearchRepositoryTestSuite) TestSearchPosts(t provider.T) {
	t.Tags("search", "posts")
	t.Description("Test post search functionality")

	t.Run("Search posts by title", func(t provider.T) {
		ctx := context.Background()
		u := uuid.NewString()
		t.WithNewStep("Create test posts", func(pctx provider.StepCtx) {
			post1 := s.createTestPost(ctx, t)
			post1.UpdateTitle(u + "Unique Title 1")

			post2 := s.createAuthorPost(ctx, t, post1.AuthorID())
			post2.UpdateTitle(u + "Unique Title 2")

			_, err := s.postRepo.Create(ctx, post1)
			require.NoError(t, err)
			_, err = s.postRepo.Create(ctx, post2)
			require.NoError(t, err)
		})

		t.WithNewStep("Search by title", func(pctx provider.StepCtx) {
			srch := search.NewPostSearch()
			srch.AddFilter(search.NewPostTitleContainsFilter(u + "Unique Title 1"))

			posts, err := s.searchRepo.SearchPosts(ctx, srch)
			require.NoError(t, err)

			assert.Len(t, posts, 1)
			assert.Equal(t, u+"Unique Title 1", posts[0].Title())
		})
	})

	t.Run("Search posts by tag", func(t provider.T) {
		ctx := context.Background()

		u := uuid.NewString()

		t.WithNewStep("Create test posts", func(pctx provider.StepCtx) {
			post1 := s.createTestPost(ctx, t)
			post1.UpdateTags([]string{u + "gardening", u + "tips"})

			post2 := s.createTestPost(ctx, t)
			post2.UpdateTags([]string{u + "plants", u + "care"})

			_, err := s.postRepo.Create(ctx, post1)
			require.NoError(t, err)
			_, err = s.postRepo.Create(ctx, post2)
			require.NoError(t, err)
		})

		t.WithNewStep("Search by tag", func(pctx provider.StepCtx) {
			srch := search.NewPostSearch()
			srch.AddFilter(search.NewPostTagFilter([]string{u + "gardening"}))

			posts, err := s.searchRepo.SearchPosts(ctx, srch)
			require.NoError(t, err)

			assert.Len(t, posts, 1)
			assert.Contains(t, posts[0].Tags(), u+"gardening")
		})
	})
}

func (s *SearchRepositoryTestSuite) TestSearchPlants(t provider.T) {
	t.Tags("search", "plants")
	t.Description("Test plant search functionality")

	t.Run("Search plants by type", func(t provider.T) {
		ctx := context.Background()

		u := uuid.NewString()

		t.WithNewStep("Create test plants", func(pctx provider.StepCtx) {
			u := uuid.NewString()
			coniferousPlant := s.createDeciduousPlant(
				ctx,
				t,
				u+"Pine Tree",
				1.5,
				0.5,
				plant.MediumMoisture,
				10,
				plant.Light,
				plant.MediumSoil,
				plant.WinterHardiness(10),
				plant.Spring,
			)
			deciduousPlant := s.createDeciduousPlant(
				ctx,
				t,
				u+"Oak Tree",
				5.0,
				1.0,
				plant.HighMoisture,
				10,
				plant.Light,
				plant.MediumSoil,
				plant.WinterHardiness(10),
				plant.Spring,
			)

			_, err := s.plantRepo.Create(ctx, coniferousPlant)
			require.NoError(t, err)
			_, err = s.plantRepo.Create(ctx, deciduousPlant)
			require.NoError(t, err)
		})

		t.WithNewStep("Search coniferous plants", func(pctx provider.StepCtx) {
			coniferousSearch := search.NewPlantSearch()
			coniferousSearch.AddFilter(search.NewPlantCategoryFilter(plant.ConiferousCategory))
			coniferousSearch.AddFilter(search.NewPlantNameFilter(u))
			coniferousPlants, err := s.searchRepo.SearchPlants(ctx, coniferousSearch)
			require.NoError(t, err)

			assert.Len(t, coniferousPlants, 0)
		})
	})

	t.Run("Search plants by height", func(t provider.T) {
		ctx := context.Background()

		u := uuid.NewString()

		t.WithNewStep("Create test plants", func(pctx provider.StepCtx) {
			tallConifer := s.createConiferousPlant(
				ctx,
				t,
				u+"Tall Pine",
				10.0,
				2.0,
				plant.MediumMoisture,
				10,
				plant.Light,
				plant.MediumSoil,
				plant.WinterHardiness(10),
			)
			shortConifer := s.createConiferousPlant(
				ctx,
				t,
				u+"Short Pine",
				1.5,
				0.5,
				plant.MediumMoisture,
				10,
				plant.Light,
				plant.MediumSoil,
				plant.WinterHardiness(10),
			)

			_, err := s.plantRepo.Create(ctx, tallConifer)
			require.NoError(t, err)
			_, err = s.plantRepo.Create(ctx, shortConifer)
			require.NoError(t, err)
		})

		t.WithNewStep("Search tall plants", func(pctx provider.StepCtx) {
			tallSearch := search.NewPlantSearch()
			tallSearch.AddFilter(search.NewPlantHeightFilter(8.0, 15.0))
			tallSearch.AddFilter(search.NewPlantNameFilter(u))
			tallPlants, err := s.searchRepo.SearchPlants(ctx, tallSearch)
			require.NoError(t, err)

			assert.Len(t, tallPlants, 1)
			assert.Equal(t, u+"Tall Pine", tallPlants[0].GetName())
		})
	})
}

func (s *SearchRepositoryTestSuite) TestGetPostByID(t provider.T) {
	t.Tags("get", "post")
	t.Description("Test post retrieval by ID")

	ctx := context.Background()
	testPost := s.createTestPost(ctx, t)

	t.WithNewStep("Create post", func(pctx provider.StepCtx) {
		_, err := s.postRepo.Create(ctx, testPost)
		require.NoError(t, err)
	})

	t.WithNewStep("Retrieve post", func(pctx provider.StepCtx) {
		retrievedPost, err := s.searchRepo.GetPostByID(ctx, testPost.ID())
		require.NoError(t, err)

		assert.Equal(t, testPost.ID(), retrievedPost.ID())
		assert.Equal(t, testPost.Title(), retrievedPost.Title())
		assert.Equal(t, testPost.Content().Text, retrievedPost.Content().Text)
		assert.Equal(t, testPost.AuthorID(), retrievedPost.AuthorID())
	})
}

func (s *SearchRepositoryTestSuite) TestGetPlantByID(t provider.T) {
	t.Tags("get", "plant")
	t.Description("Test plant retrieval by ID")

	ctx := context.Background()
	testPlant := s.createConiferousPlant(
		ctx,
		t,
		"Test Plant",
		1.0,
		0.5,
		plant.MediumMoisture,
		10,
		plant.Light,
		plant.MediumSoil,
		plant.WinterHardiness(10),
	)

	t.WithNewStep("Create plant", func(pctx provider.StepCtx) {
		_, err := s.plantRepo.Create(ctx, testPlant)
		require.NoError(t, err)
	})

	t.WithNewStep("Retrieve plant", func(pctx provider.StepCtx) {
		retrievedPlant, err := s.searchRepo.GetPlantByID(ctx, testPlant.ID())
		require.NoError(t, err)

		assert.Equal(t, testPlant.ID(), retrievedPlant.ID())
		assert.Equal(t, testPlant.GetName(), retrievedPlant.GetName())
		assert.Equal(t, testPlant.GetLatinName(), retrievedPlant.GetLatinName())
	})
}
