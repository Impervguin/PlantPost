//go:build integration

package authstorage_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"PlantSite/internal/infra/sqpgx"
	"PlantSite/internal/models/auth"
	authstorage "PlantSite/internal/repositories/postgres/auth-storage"
	"PlantSite/internal/repositories/tests"
	"PlantSite/internal/testutils/pgtest"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

type AuthRepositoryTestSuite struct {
	suite.Suite
	container testcontainers.Container
	db        *sqpgx.SquirrelPgx
	repo      *authstorage.PostgresAuthRepository
	prevDir   string
	cntCreds  pgtest.PostgresCredentials
}

func TestAuthRepositorySuite(t *testing.T) {
	suite.RunSuite(t, new(AuthRepositoryTestSuite))
}

func (s *AuthRepositoryTestSuite) BeforeEach(t provider.T) {
	t.Epic("Auth Repository")
	t.Feature("Authentication Storage")

	err := pgtest.Migrate(context.Background(), &s.cntCreds)
	require.NoError(t, err)
}

func (s *AuthRepositoryTestSuite) BeforeAll(t provider.T) {
	ctx := context.Background()

	// Save current directory
	prevDir, err := os.Getwd()
	require.NoError(t, err)
	s.prevDir = prevDir

	// Change directory to test working directory
	err = os.Chdir(tests.GetTestWorkingDir())
	require.NoError(t, err)

	// Create new container
	container, creds, err := pgtest.NewTestPostgres(ctx)
	require.NoError(t, err)
	s.container = container
	s.cntCreds = creds

	// Create database connection
	config := &sqpgx.SqpgxConfig{
		User:                   creds.User,
		Password:               creds.Password,
		DbName:                 creds.Database,
		Host:                   creds.Host,
		Port:                   creds.Port,
		MaxConnections:         10,
		MaxConnectionsLifetime: time.Minute,
	}

	db, err := sqpgx.NewSquirrelPgx(ctx, config)
	require.NoError(t, err)
	s.db = db

	// Create repository
	repo, err := authstorage.NewPostgresAuthRepository(ctx, db)
	require.NoError(t, err)
	s.repo = repo
}

func (s *AuthRepositoryTestSuite) AfterAll(t provider.T) {
	if s.container != nil {
		s.container.Terminate(context.Background())
	}
	err := os.Chdir(s.prevDir)
	require.NoError(t, err)
}

func (s *AuthRepositoryTestSuite) AfterEach(t provider.T) {
	err := pgtest.MigrateDown(context.Background(), &s.cntCreds)
	require.NoError(t, err)
}

func (s *AuthRepositoryTestSuite) createTestMember(t provider.T) *auth.Member {
	memID := uuid.New()
	user, err := auth.CreateMember(
		memID,
		memID.String()[:8],
		memID.String()[:8]+"@example.com",
		[]byte("hashedpassword"),
		time.Now(),
	)
	require.NoError(t, err)
	return user
}

func (s *AuthRepositoryTestSuite) createTestAuthor(t provider.T) *auth.Author {
	member := s.createTestMember(t)
	author, err := auth.CreateAuthor(
		*member,
		time.Now(),
		true,
		time.Now().Add(-512*time.Hour),
	)
	require.NoError(t, err)
	return author
}

func (s *AuthRepositoryTestSuite) TestCreateMember(t provider.T) {
	t.Tags("create", "member")
	t.Description("Test member creation functionality")

	ctx := context.Background()
	testMember := s.createTestMember(t)

	t.WithNewStep("Create member", func(pctx provider.StepCtx) {
		createdUser, err := s.repo.Create(ctx, testMember)
		require.NoError(t, err)

		createdMember, ok := createdUser.(*auth.Member)
		require.True(t, ok)
		assert.Equal(t, testMember.ID(), createdMember.ID())
		assert.Equal(t, testMember.Name(), createdMember.Name())
		assert.Equal(t, testMember.Email(), createdMember.Email())
		assert.Equal(t, testMember.HashedPassword(), createdMember.HashedPassword())
	})

	t.WithNewStep("Verify member retrieval", func(pctx provider.StepCtx) {
		fetchedUser, err := s.repo.Get(ctx, testMember.ID())
		require.NoError(t, err)
		fetchedMember, ok := fetchedUser.(*auth.Member)
		require.True(t, ok)
		assert.Equal(t, testMember.ID(), fetchedMember.ID())
	})
}

func (s *AuthRepositoryTestSuite) TestCreateDuplicateMember(t provider.T) {
	t.Tags("create", "duplicate")
	t.Description("Test duplicate member creation handling")

	ctx := context.Background()
	testMember := s.createTestMember(t)

	t.WithNewStep("Create first member", func(pctx provider.StepCtx) {
		_, err := s.repo.Create(ctx, testMember)
		require.NoError(t, err)
	})

	t.WithNewStep("Attempt duplicate creation", func(pctx provider.StepCtx) {
		_, err := s.repo.Create(ctx, testMember)
		require.Error(t, err)
	})
}

func (s *AuthRepositoryTestSuite) TestGetMember(t provider.T) {
	t.Tags("get", "member")
	t.Description("Test member retrieval functionality")

	ctx := context.Background()
	testMember := s.createTestMember(t)

	t.WithNewStep("Create member", func(pctx provider.StepCtx) {
		_, err := s.repo.Create(ctx, testMember)
		require.NoError(t, err)
	})

	t.WithNewStep("Retrieve member", func(pctx provider.StepCtx) {
		fetchedUser, err := s.repo.Get(ctx, testMember.ID())
		require.NoError(t, err)

		fetchedMember, ok := fetchedUser.(*auth.Member)
		require.True(t, ok)
		assert.Equal(t, testMember.ID(), fetchedMember.ID())
		assert.Equal(t, testMember.Name(), fetchedMember.Name())
		assert.Equal(t, testMember.Email(), fetchedMember.Email())
		assert.Equal(t, testMember.HashedPassword(), fetchedMember.HashedPassword())
	})
}

func (s *AuthRepositoryTestSuite) TestGetNonExistentMember(t provider.T) {
	t.Tags("get", "not_found")
	t.Description("Test non-existent member retrieval handling")

	ctx := context.Background()
	nonExistentID := uuid.New()

	t.WithNewStep("Attempt get non-existent member", func(pctx provider.StepCtx) {
		_, err := s.repo.Get(ctx, nonExistentID)
		require.Error(t, err)
		assert.Equal(t, auth.ErrUserNotFound, err)
	})
}

func (s *AuthRepositoryTestSuite) TestGetByName(t provider.T) {
	t.Tags("get", "by_name")
	t.Description("Test member retrieval by name functionality")

	ctx := context.Background()
	testMember := s.createTestMember(t)

	t.WithNewStep("Create member", func(pctx provider.StepCtx) {
		_, err := s.repo.Create(ctx, testMember)
		require.NoError(t, err)
	})

	t.WithNewStep("Retrieve by name", func(pctx provider.StepCtx) {
		fetchedUser, err := s.repo.GetByName(ctx, testMember.Name())
		require.NoError(t, err)

		fetchedMember, ok := fetchedUser.(*auth.Member)
		require.True(t, ok)
		assert.Equal(t, testMember.ID(), fetchedMember.ID())
		assert.Equal(t, testMember.Name(), fetchedMember.Name())
	})
}

func (s *AuthRepositoryTestSuite) TestGetByEmail(t provider.T) {
	t.Tags("get", "by_email")
	t.Description("Test member retrieval by email functionality")

	ctx := context.Background()
	testMember := s.createTestMember(t)

	t.WithNewStep("Create member", func(pctx provider.StepCtx) {
		_, err := s.repo.Create(ctx, testMember)
		require.NoError(t, err)
	})

	t.WithNewStep("Retrieve by email", func(pctx provider.StepCtx) {
		fetchedUser, err := s.repo.GetByEmail(ctx, testMember.Email())
		require.NoError(t, err)

		fetchedMember, ok := fetchedUser.(*auth.Member)
		require.True(t, ok)
		assert.Equal(t, testMember.ID(), fetchedMember.ID())
		assert.Equal(t, testMember.Email(), fetchedMember.Email())
	})
}

func (s *AuthRepositoryTestSuite) TestUpdateMember(t provider.T) {
	t.Tags("update", "member")
	t.Description("Test member update functionality")

	ctx := context.Background()
	testMember := s.createTestMember(t)

	t.WithNewStep("Create initial member", func(pctx provider.StepCtx) {
		_, err := s.repo.Create(ctx, testMember)
		require.NoError(t, err)
	})

	t.WithNewStep("Update member", func(pctx provider.StepCtx) {
		newName := "updateduser"
		newEmail := "updated@example.com"
		newPassword := []byte("updatedhash")

		updatedUser, err := s.repo.Update(ctx, testMember.ID(), func(u auth.User) (auth.User, error) {
			switch fact := u.(type) {
			case *auth.Member:
				fact.UpdateName(newName)
				fact.UpdateEmail(newEmail)
				fact.UpdateHashedPassword(newPassword)
				return fact, nil
			case *auth.Author:
				fact.UpdateName(newName)
				fact.UpdateEmail(newEmail)
				fact.UpdateHashedPassword(newPassword)
				return fact, nil
			default:
				return nil, fmt.Errorf("unsupported user type: %v", u)
			}
		})
		require.NoError(t, err)

		updatedMember, ok := updatedUser.(*auth.Member)
		require.True(t, ok)
		assert.Equal(t, newName, updatedMember.Name())
		assert.Equal(t, newEmail, updatedMember.Email())
		assert.Equal(t, newPassword, updatedMember.HashedPassword())
	})

	t.WithNewStep("Verify persistence", func(pctx provider.StepCtx) {
		fetchedUser, err := s.repo.Get(ctx, testMember.ID())
		require.NoError(t, err)
		fetchedMember, ok := fetchedUser.(*auth.Member)
		require.True(t, ok)
		assert.Equal(t, "updateduser", fetchedMember.Name())
		assert.Equal(t, "updated@example.com", fetchedMember.Email())
		assert.Equal(t, []byte("updatedhash"), fetchedMember.HashedPassword())
	})
}

func (s *AuthRepositoryTestSuite) TestUpdateNonExistentMember(t provider.T) {
	t.Tags("update", "not_found")
	t.Description("Test non-existent member update handling")

	ctx := context.Background()
	nonExistentID := uuid.New()

	t.WithNewStep("Attempt update non-existent member", func(pctx provider.StepCtx) {
		_, err := s.repo.Update(ctx, nonExistentID, func(u auth.User) (auth.User, error) {
			switch fact := u.(type) {
			case *auth.Member:
				fact.UpdateName("shouldfail")
				return fact, nil
			case *auth.Author:
				fact.UpdateName("shouldfail")
				return fact, nil
			default:
				return nil, fmt.Errorf("unsupported user type: %v", u)
			}
		})
		require.Error(t, err)
		assert.Equal(t, auth.ErrUserNotFound, err)
	})
}

func (s *AuthRepositoryTestSuite) TestPromoteToAuthor(t provider.T) {
	t.Tags("promote", "author")
	t.Description("Test member promotion to author functionality")

	ctx := context.Background()
	testMember := s.createTestMember(t)

	t.WithNewStep("Create initial member", func(pctx provider.StepCtx) {
		_, err := s.repo.Create(ctx, testMember)
		require.NoError(t, err)
	})

	t.WithNewStep("Promote to author", func(pctx provider.StepCtx) {
		grantTime := time.Now()
		updatedUser, err := s.repo.Update(ctx, testMember.ID(), func(u auth.User) (auth.User, error) {
			member := u.(*auth.Member)
			return auth.CreateAuthor(*member, grantTime, true, grantTime.Add(-512*time.Hour))
		})
		require.NoError(t, err)

		updatedAuthor, ok := updatedUser.(*auth.Author)
		require.True(t, ok)
		assert.Equal(t, testMember.ID(), updatedAuthor.ID())
		assert.Equal(t, testMember.Name(), updatedAuthor.Name())
		assert.True(t, updatedAuthor.HasRights())
	})

	t.WithNewStep("Verify author persistence", func(pctx provider.StepCtx) {
		fetchedUser, err := s.repo.Get(ctx, testMember.ID())
		require.NoError(t, err)
		fetchedAuthor, ok := fetchedUser.(*auth.Author)
		require.True(t, ok)
		assert.Equal(t, testMember.ID(), fetchedAuthor.ID())
		assert.True(t, fetchedAuthor.HasRights())
	})
}

func (s *AuthRepositoryTestSuite) TestRevokeAuthorRights(t provider.T) {
	t.Tags("revoke", "author")
	t.Description("Test author rights revocation functionality")

	ctx := context.Background()
	testAuthor := s.createTestAuthor(t)

	t.WithNewStep("Create initial author", func(pctx provider.StepCtx) {
		_, err := s.repo.Create(ctx, &testAuthor.Member)
		require.NoError(t, err)
		_, err = s.repo.Update(ctx, testAuthor.ID(), func(u auth.User) (auth.User, error) {
			return testAuthor, nil
		})
		require.NoError(t, err)
	})

	t.WithNewStep("Revoke author rights", func(pctx provider.StepCtx) {
		updatedUser, err := s.repo.Update(ctx, testAuthor.ID(), func(u auth.User) (auth.User, error) {
			switch fact := u.(type) {
			case *auth.Author:
				fact.RevokeAuthorRights()
				return fact, nil
			default:
				return nil, fmt.Errorf("unsupported user type: %v", u)
			}
		})
		require.NoError(t, err)

		updatedAuthor, ok := updatedUser.(*auth.Author)
		require.True(t, ok)
		assert.False(t, updatedAuthor.HasRights())
	})

	t.WithNewStep("Verify rights revocation persistence", func(pctx provider.StepCtx) {
		fetchedUser, err := s.repo.Get(ctx, testAuthor.ID())
		require.NoError(t, err)
		fetchedAuthor, ok := fetchedUser.(*auth.Author)
		require.True(t, ok)
		assert.False(t, fetchedAuthor.HasRights())
	})
}

func (s *AuthRepositoryTestSuite) TestGetAuthor(t provider.T) {
	t.Tags("get", "author")
	t.Description("Test author retrieval functionality")

	ctx := context.Background()
	testAuthor := s.createTestAuthor(t)

	t.WithNewStep("Create author", func(pctx provider.StepCtx) {
		_, err := s.repo.Create(ctx, &testAuthor.Member)
		require.NoError(t, err)
		_, err = s.repo.Update(ctx, testAuthor.ID(), func(u auth.User) (auth.User, error) {
			return testAuthor, nil
		})
		require.NoError(t, err)
	})

	t.WithNewStep("Retrieve author", func(pctx provider.StepCtx) {
		fetchedUser, err := s.repo.Get(ctx, testAuthor.ID())
		require.NoError(t, err)

		fetchedAuthor, ok := fetchedUser.(*auth.Author)
		require.True(t, ok)
		assert.Equal(t, testAuthor.ID(), fetchedAuthor.ID())
		assert.Equal(t, testAuthor.Name(), fetchedAuthor.Name())
		assert.True(t, fetchedAuthor.HasRights())
	})
}
