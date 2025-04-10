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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type AuthRepositoryTestSuite struct {
	suite.Suite
	container testcontainers.Container
	db        *sqpgx.SquirrelPgx
	repo      *authstorage.PostgresAuthRepository
	prevDir   string
}

func TestAuthRepositorySuite(t *testing.T) {
	suite.Run(t, new(AuthRepositoryTestSuite))
}

func (s *AuthRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	// Save current directory
	prevDir, err := os.Getwd()
	require.NoError(s.T(), err)
	s.prevDir = prevDir

	// Change directory to test working directory
	err = os.Chdir(tests.GetTestWorkingDir())
	require.NoError(s.T(), err)

	// Create new container
	container, creds, err := pgtest.NewTestPostgres(ctx)
	require.NoError(s.T(), err)
	s.container = container

	// Run migrations
	err = pgtest.Migrate(ctx, &creds)
	require.NoError(s.T(), err)

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
	require.NoError(s.T(), err)
	s.db = db

	// Create repository
	repo, err := authstorage.NewPostgresAuthRepository(ctx, db)
	require.NoError(s.T(), err)
	s.repo = repo
}

func (s *AuthRepositoryTestSuite) TearDownSuite() {
	if s.container != nil {
		s.container.Terminate(context.Background())
	}
	err := os.Chdir(s.prevDir)
	require.NoError(s.T(), err)
}

func (s *AuthRepositoryTestSuite) createTestMember() *auth.Member {
	memID := uuid.New()
	user, err := auth.CreateMember(
		memID,
		memID.String()[:8],
		memID.String()[:8]+"@example.com",
		[]byte("hashedpassword"),
		time.Now(),
	)
	require.NoError(s.T(), err)
	return user
}

func (s *AuthRepositoryTestSuite) createTestAuthor() *auth.Author {
	member := s.createTestMember()
	author, err := auth.CreateAuthor(
		*member,
		time.Now(),
		true,
		time.Now().Add(-512*time.Hour),
	)
	require.NoError(s.T(), err)
	return author
}

func (s *AuthRepositoryTestSuite) TestCreateMember() {
	ctx := context.Background()
	testMember := s.createTestMember()

	// Test creation
	createdUser, err := s.repo.Create(ctx, testMember)
	require.NoError(s.T(), err)

	// Verify returned user matches input
	createdMember, ok := createdUser.(*auth.Member)
	require.True(s.T(), ok)
	assert.Equal(s.T(), testMember.ID(), createdMember.ID())
	assert.Equal(s.T(), testMember.Name(), createdMember.Name())
	assert.Equal(s.T(), testMember.Email(), createdMember.Email())
	assert.Equal(s.T(), testMember.HashedPassword(), createdMember.HashedPassword())

	// Verify can retrieve
	fetchedUser, err := s.repo.Get(ctx, testMember.ID())
	require.NoError(s.T(), err)
	fetchedMember, ok := fetchedUser.(*auth.Member)
	require.True(s.T(), ok)
	assert.Equal(s.T(), testMember.ID(), fetchedMember.ID())
}

func (s *AuthRepositoryTestSuite) TestCreateDuplicateMember() {
	ctx := context.Background()
	testMember := s.createTestMember()

	// First creation should succeed
	_, err := s.repo.Create(ctx, testMember)
	require.NoError(s.T(), err)

	// Second creation with same ID should fail
	_, err = s.repo.Create(ctx, testMember)
	require.Error(s.T(), err)
}

func (s *AuthRepositoryTestSuite) TestGetMember() {
	ctx := context.Background()
	testMember := s.createTestMember()

	// Create member first
	_, err := s.repo.Create(ctx, testMember)
	require.NoError(s.T(), err)

	// Test retrieval
	fetchedUser, err := s.repo.Get(ctx, testMember.ID())
	require.NoError(s.T(), err)

	// Verify type and fields
	fetchedMember, ok := fetchedUser.(*auth.Member)
	require.True(s.T(), ok)
	assert.Equal(s.T(), testMember.ID(), fetchedMember.ID())
	assert.Equal(s.T(), testMember.Name(), fetchedMember.Name())
	assert.Equal(s.T(), testMember.Email(), fetchedMember.Email())
	assert.Equal(s.T(), testMember.HashedPassword(), fetchedMember.HashedPassword())
}

func (s *AuthRepositoryTestSuite) TestGetNonExistentMember() {
	ctx := context.Background()
	nonExistentID := uuid.New()

	_, err := s.repo.Get(ctx, nonExistentID)
	require.Error(s.T(), err)
	assert.Equal(s.T(), auth.ErrUserNotFound, err)
}

func (s *AuthRepositoryTestSuite) TestGetByName() {
	ctx := context.Background()
	testMember := s.createTestMember()

	// Create member first
	_, err := s.repo.Create(ctx, testMember)
	require.NoError(s.T(), err)

	// Test retrieval by name
	fetchedUser, err := s.repo.GetByName(ctx, testMember.Name())
	require.NoError(s.T(), err)

	// Verify fields
	fetchedMember, ok := fetchedUser.(*auth.Member)
	require.True(s.T(), ok)
	assert.Equal(s.T(), testMember.ID(), fetchedMember.ID())
	assert.Equal(s.T(), testMember.Name(), fetchedMember.Name())
}

func (s *AuthRepositoryTestSuite) TestGetByEmail() {
	ctx := context.Background()
	testMember := s.createTestMember()

	// Create member first
	_, err := s.repo.Create(ctx, testMember)
	require.NoError(s.T(), err)

	// Test retrieval by email
	fetchedUser, err := s.repo.GetByEmail(ctx, testMember.Email())
	require.NoError(s.T(), err)

	// Verify fields
	fetchedMember, ok := fetchedUser.(*auth.Member)
	require.True(s.T(), ok)
	assert.Equal(s.T(), testMember.ID(), fetchedMember.ID())
	assert.Equal(s.T(), testMember.Email(), fetchedMember.Email())
}

func (s *AuthRepositoryTestSuite) TestUpdateMember() {
	ctx := context.Background()
	testMember := s.createTestMember()

	// Create initial member
	_, err := s.repo.Create(ctx, testMember)
	require.NoError(s.T(), err)

	// Test update
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
	require.NoError(s.T(), err)

	// Verify updates
	updatedMember, ok := updatedUser.(*auth.Member)
	require.True(s.T(), ok)
	assert.Equal(s.T(), newName, updatedMember.Name())
	assert.Equal(s.T(), newEmail, updatedMember.Email())
	assert.Equal(s.T(), newPassword, updatedMember.HashedPassword())

	// Verify persistence
	fetchedUser, err := s.repo.Get(ctx, testMember.ID())
	require.NoError(s.T(), err)
	fetchedMember, ok := fetchedUser.(*auth.Member)
	require.True(s.T(), ok)
	assert.Equal(s.T(), newName, fetchedMember.Name())
	assert.Equal(s.T(), newEmail, fetchedMember.Email())
	assert.Equal(s.T(), newPassword, fetchedMember.HashedPassword())
}

func (s *AuthRepositoryTestSuite) TestUpdateNonExistentMember() {
	ctx := context.Background()
	nonExistentID := uuid.New()

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
	require.Error(s.T(), err)
	assert.Equal(s.T(), auth.ErrUserNotFound, err)
}

func (s *AuthRepositoryTestSuite) TestPromoteToAuthor() {
	ctx := context.Background()
	testMember := s.createTestMember()

	// Create initial member
	_, err := s.repo.Create(ctx, testMember)
	require.NoError(s.T(), err)

	// Test promotion to author
	grantTime := time.Now()
	updatedUser, err := s.repo.Update(ctx, testMember.ID(), func(u auth.User) (auth.User, error) {
		member := u.(*auth.Member)
		return auth.CreateAuthor(*member, grantTime, true, grantTime.Add(-512*time.Hour))
	})
	require.NoError(s.T(), err)

	// Verify type and fields
	updatedAuthor, ok := updatedUser.(*auth.Author)
	require.True(s.T(), ok)
	assert.Equal(s.T(), testMember.ID(), updatedAuthor.ID())
	assert.Equal(s.T(), testMember.Name(), updatedAuthor.Name())
	assert.True(s.T(), updatedAuthor.HasRights())

	// Verify persistence
	fetchedUser, err := s.repo.Get(ctx, testMember.ID())
	require.NoError(s.T(), err)
	fetchedAuthor, ok := fetchedUser.(*auth.Author)
	require.True(s.T(), ok)
	assert.Equal(s.T(), testMember.ID(), fetchedAuthor.ID())
	assert.True(s.T(), fetchedAuthor.HasRights())
}

func (s *AuthRepositoryTestSuite) TestRevokeAuthorRights() {
	ctx := context.Background()
	testAuthor := s.createTestAuthor()

	// Create initial author
	_, err := s.repo.Create(ctx, &testAuthor.Member)
	require.NoError(s.T(), err)
	_, err = s.repo.Update(ctx, testAuthor.ID(), func(u auth.User) (auth.User, error) {
		return testAuthor, nil
	})
	require.NoError(s.T(), err)

	// Test revoking rights
	updatedUser, err := s.repo.Update(ctx, testAuthor.ID(), func(u auth.User) (auth.User, error) {
		switch fact := u.(type) {
		case *auth.Author:
			fact.RevokeAuthorRights()
			return fact, nil
		default:
			return nil, fmt.Errorf("unsupported user type: %v", u)
		}
	})
	require.NoError(s.T(), err)

	// Verify updates
	updatedAuthor, ok := updatedUser.(*auth.Author)
	require.True(s.T(), ok)
	assert.False(s.T(), updatedAuthor.HasRights())

	// Verify persistence
	fetchedUser, err := s.repo.Get(ctx, testAuthor.ID())
	require.NoError(s.T(), err)
	fetchedAuthor, ok := fetchedUser.(*auth.Author)
	require.True(s.T(), ok)
	assert.False(s.T(), fetchedAuthor.HasRights())
}

func (s *AuthRepositoryTestSuite) TestGetAuthor() {
	ctx := context.Background()
	testAuthor := s.createTestAuthor()

	// Create author first
	_, err := s.repo.Create(ctx, &testAuthor.Member)
	require.NoError(s.T(), err)
	_, err = s.repo.Update(ctx, testAuthor.ID(), func(u auth.User) (auth.User, error) {
		return testAuthor, nil
	})
	require.NoError(s.T(), err)

	// Test retrieval
	fetchedUser, err := s.repo.Get(ctx, testAuthor.ID())
	require.NoError(s.T(), err)

	// Verify type and fields
	fetchedAuthor, ok := fetchedUser.(*auth.Author)
	require.True(s.T(), ok)
	assert.Equal(s.T(), testAuthor.ID(), fetchedAuthor.ID())
	assert.Equal(s.T(), testAuthor.Name(), fetchedAuthor.Name())
	assert.True(s.T(), fetchedAuthor.HasRights())
}
