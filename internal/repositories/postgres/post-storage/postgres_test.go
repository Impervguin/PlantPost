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

type PostRepositoryTestSuite struct {
	suite.Suite
	dbContainer testcontainers.Container
	dbCreds     *pgtest.PostgresCredentials
	fileCnt     testcontainers.Container
	fileCreds   *miniotest.MinioCredentials
	db          *sqpgx.SquirrelPgx
	fileRepo    *filestorage.PgMinioStorage
	repo        *poststorage.PostgresPostRepository
	userRepo    *authstorage.PostgresAuthRepository
	plantRepo   *plantstorage.PostgresPlantRepository
	prevDir     string
}

func TestPostRepositorySuite(t *testing.T) {
	suite.RunSuite(t, new(PostRepositoryTestSuite))
}

func (s *PostRepositoryTestSuite) BeforeEach(t provider.T) {
	t.Epic("Post Repository")
	t.Feature("Post Storage")
}

func (s *PostRepositoryTestSuite) BeforeAll(t provider.T) {
	ctx := context.Background()

	// Save current directory
	prevDir, err := os.Getwd()
	require.NoError(t, err)
	s.prevDir = prevDir

	err = os.Chdir(tests.GetTestWorkingDir())
	require.NoError(t, err)

	// Setup PostgreSQL container
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
	s.plantRepo, err = plantstorage.NewPostgresPlantRepository(ctx, s.db)
	require.NoError(t, err)

	// Create search repository
	searchRepo, err := searchstorage.NewPostgresSearchRepository(ctx, s.db)
	require.NoError(t, err)

	// Create search plant getter
	plntGetter := searchstorage.NewSearchPlantGetter(searchRepo)

	// Create post repository
	s.repo, err = poststorage.NewPostgresPostRepository(ctx, s.db, plntGetter)
	require.NoError(t, err)

	// Create user repository
	s.userRepo, err = authstorage.NewPostgresAuthRepository(ctx, s.db)
	require.NoError(t, err)
}

func (s *PostRepositoryTestSuite) AfterAll(t provider.T) {
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

func (s *PostRepositoryTestSuite) AfterEach(t provider.T) {
	err := pgtest.TruncateTables(context.Background(), s.dbCreds)
	require.NoError(t, err)
	err = miniotest.CleanUpBucket(context.Background(), s.fileCreds)
	require.NoError(t, err)
}

func (s *PostRepositoryTestSuite) uploadTestPhoto(ctx context.Context, t provider.T) uuid.UUID {
	id := uuid.New()
	fileData := models.FileData{
		Name:        id.String() + ".jpg",
		Reader:      bytes.NewReader([]byte("test photo content")),
		ContentType: "image/jpeg",
	}
	file, err := s.fileRepo.Upload(ctx, &fileData)
	require.NoError(t, err)
	return file.ID
}

func (s *PostRepositoryTestSuite) pushAuthor(ctx context.Context, t provider.T) uuid.UUID {
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

func (s *PostRepositoryTestSuite) createTestPost(ctx context.Context, t provider.T) *post.Post {
	ownerID := s.pushAuthor(ctx, t)
	return s.createAuthorPost(ctx, t, ownerID)
}

func (s *PostRepositoryTestSuite) createAuthorPost(ctx context.Context, t provider.T, ownerID uuid.UUID) *post.Post {
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

func (s *PostRepositoryTestSuite) TestCreatePost(t provider.T) {
	t.Tags("create", "post")
	t.Description("Test post creation functionality")

	t.Run("Successful post creation", func(t provider.T) {
		ctx := context.Background()
		testPost := s.createTestPost(ctx, t)

		t.WithNewStep("Create post", func(pctx provider.StepCtx) {
			createdPost, err := s.repo.Create(ctx, testPost)
			require.NoError(t, err)

			assert.Equal(t, testPost.ID(), createdPost.ID())
			assert.Equal(t, testPost.Title(), createdPost.Title())
			assert.Equal(t, testPost.Content().Text, createdPost.Content().Text)
			assert.Equal(t, testPost.AuthorID(), createdPost.AuthorID())
			assert.Equal(t, testPost.Photos().Len(), createdPost.Photos().Len())
			assert.ElementsMatch(t, testPost.Tags(), createdPost.Tags())
		})

		t.WithNewStep("Verify post creation", func(pctx provider.StepCtx) {
			fetchedPost, err := s.repo.Get(ctx, testPost.ID())
			require.NoError(t, err)
			assert.Equal(t, testPost.ID(), fetchedPost.ID())
		})
	})

	t.Run("Post with no photos", func(t provider.T) {
		ctx := context.Background()
		testPost := s.createTestPost(ctx, t)

		t.WithNewStep("Remove all photos", func(pctx provider.StepCtx) {
			testPost.ClearPhotos()
		})

		t.WithNewStep("Create post without photos", func(pctx provider.StepCtx) {
			createdPost, err := s.repo.Create(ctx, testPost)
			require.NoError(t, err)
			assert.Equal(t, 0, createdPost.Photos().Len())
		})

		t.WithNewStep("Verify post creation", func(pctx provider.StepCtx) {
			fetchedPost, err := s.repo.Get(ctx, testPost.ID())
			require.NoError(t, err)
			assert.Equal(t, 0, fetchedPost.Photos().Len())
		})
	})
}

func (s *PostRepositoryTestSuite) TestGetPost(t provider.T) {
	t.Tags("get", "post")
	t.Description("Test post retrieval functionality")

	t.Run("Successful post retrieval", func(t provider.T) {
		ctx := context.Background()
		testPost := s.createTestPost(ctx, t)

		t.WithNewStep("Create post", func(pctx provider.StepCtx) {
			_, err := s.repo.Create(ctx, testPost)
			require.NoError(t, err)
		})

		t.WithNewStep("Retrieve post", func(pctx provider.StepCtx) {
			retrievedPost, err := s.repo.Get(ctx, testPost.ID())
			require.NoError(t, err)

			assert.Equal(t, testPost.ID(), retrievedPost.ID())
			assert.Equal(t, testPost.Title(), retrievedPost.Title())
			assert.Equal(t, testPost.Content().Text, retrievedPost.Content().Text)
			assert.Equal(t, testPost.AuthorID(), retrievedPost.AuthorID())
			assert.Equal(t, testPost.Photos().Len(), retrievedPost.Photos().Len())
			assert.ElementsMatch(t, testPost.Tags(), retrievedPost.Tags())
		})
	})

	t.Run("Non-existent post", func(t provider.T) {
		ctx := context.Background()

		t.WithNewStep("Attempt to retrieve non-existent post", func(pctx provider.StepCtx) {
			_, err := s.repo.Get(ctx, uuid.New())
			require.Error(t, err)
		})
	})
}

func (s *PostRepositoryTestSuite) TestUpdatePost(t provider.T) {
	t.Tags("update", "post")
	t.Description("Test post update functionality")

	t.Run("Successful post update", func(t provider.T) {
		ctx := context.Background()
		testPost := s.createTestPost(ctx, t)

		t.WithNewStep("Create initial post", func(pctx provider.StepCtx) {
			_, err := s.repo.Create(ctx, testPost)
			require.NoError(t, err)
		})

		t.WithNewStep("Update post", func(pctx provider.StepCtx) {
			newPhotoID := s.uploadTestPhoto(ctx, t)

			updatedPost, err := s.repo.Update(ctx, testPost.ID(), func(p *post.Post) (*post.Post, error) {
				p.UpdateTitle("Updated Post")
				content, err := post.NewContent("Updated content", post.ContentTypePlainText)
				if err != nil {
					return nil, err
				}
				p.UpdateContent(*content)
				p.UpdateTags([]string{"updated", "tags"})

				photo, err := post.CreatePostPhoto(uuid.New(), newPhotoID, 0)
				if err != nil {
					return nil, err
				}
				err = p.AddPhoto(photo)
				if err != nil {
					return nil, err
				}

				return p, nil
			})
			require.NoError(t, err)

			assert.Equal(t, "Updated Post", updatedPost.Title())
			assert.Equal(t, "Updated content", updatedPost.Content().Text)
			assert.ElementsMatch(t, []string{"updated", "tags"}, updatedPost.Tags())
			assert.Equal(t, 3, updatedPost.Photos().Len())
			assert.True(t, updatedPost.UpdatedAt().After(testPost.UpdatedAt()))
		})
	})
}

func (s *PostRepositoryTestSuite) TestDeletePost(t provider.T) {
	t.Tags("delete", "post")
	t.Description("Test post deletion functionality")

	t.Run("Successful post deletion", func(t provider.T) {
		ctx := context.Background()
		testPost := s.createTestPost(ctx, t)

		t.WithNewStep("Create post", func(pctx provider.StepCtx) {
			_, err := s.repo.Create(ctx, testPost)
			require.NoError(t, err)
		})

		t.WithNewStep("Delete post", func(pctx provider.StepCtx) {
			err := s.repo.Delete(ctx, testPost.ID())
			require.NoError(t, err)
		})

		t.WithNewStep("Verify deletion", func(pctx provider.StepCtx) {
			_, err := s.repo.Get(ctx, testPost.ID())
			require.Error(t, err)
		})
	})
}

func (s *PostRepositoryTestSuite) TestListAuthorPosts(t provider.T) {
	t.Tags("list", "author_posts")
	t.Description("Test listing author posts functionality")

	t.Run("Successful author posts listing", func(t provider.T) {
		ctx := context.Background()
		firstPost := s.createTestPost(ctx, t)

		t.WithNewStep("Create initial post", func(pctx provider.StepCtx) {
			_, err := s.repo.Create(ctx, firstPost)
			require.NoError(t, err)
		})

		t.WithNewStep("Create additional posts", func(pctx provider.StepCtx) {
			for i := 0; i < 2; i++ {
				testPost := s.createAuthorPost(ctx, t, firstPost.AuthorID())
				_, err := s.repo.Create(ctx, testPost)
				require.NoError(t, err)
			}
		})

		t.WithNewStep("List author posts", func(pctx provider.StepCtx) {
			posts, err := s.repo.ListAuthorPosts(ctx, firstPost.AuthorID())
			require.NoError(t, err)

			assert.Len(t, posts, 3)
			for _, p := range posts {
				assert.Equal(t, firstPost.AuthorID(), p.AuthorID())
			}
		})
	})

	t.Run("Empty author posts list", func(t provider.T) {
		ctx := context.Background()

		t.WithNewStep("List posts for non-existent author", func(pctx provider.StepCtx) {
			posts, err := s.repo.ListAuthorPosts(ctx, uuid.New())
			require.NoError(t, err)
			assert.Empty(t, posts)
		})
	})
}
