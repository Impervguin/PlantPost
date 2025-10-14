//go:build unit

package sessionstorage

import (
	authservice "PlantSite/internal/services/auth-service"
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MapSessionStorageTestSuite struct {
	suite.Suite
	storage *MapSessionStorage
}

func (s *MapSessionStorageTestSuite) BeforeEach(t provider.T) {
	t.Epic("Authentication")
	t.Feature("Session Storage")

	s.storage = NewMapSessionStorage()
}

func (s *MapSessionStorageTestSuite) TestNewMapSessionStorage(t provider.T) {
	t.Tags("creation", "initialization")
	t.Description("Test creation of new MapSessionStorage instance")

	storage := NewMapSessionStorage()

	t.WithNewStep("Verify storage initialization", func(ctx provider.StepCtx) {
		assert.NotNil(t, storage)
		assert.NotNil(t, storage.storage)
		assert.Empty(t, storage.storage)
	})
}

func (s *MapSessionStorageTestSuite) TestStoreAndGetSession(t provider.T) {
	t.Tags("functionality", "storage")
	t.Description("Test storing and retrieving session from storage")

	sessionID := uuid.New()
	session := &authservice.Session{
		MemberID:  uuid.New(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	t.WithNewStep("Store session in storage", func(ctx provider.StepCtx) {
		err := s.storage.Store(context.Background(), sessionID, session)
		require.NoError(t, err)
	})

	var retrievedSession *authservice.Session
	var err error
	t.WithNewStep("Retrieve session from storage", func(ctx provider.StepCtx) {
		retrievedSession, err = s.storage.Get(context.Background(), sessionID)
	})

	t.WithNewStep("Verify retrieved session", func(ctx provider.StepCtx) {
		require.NoError(t, err)
		assert.Equal(t, session.MemberID, retrievedSession.MemberID)
		assert.Equal(t, session.ExpiresAt, retrievedSession.ExpiresAt)
	})
}

func (s *MapSessionStorageTestSuite) TestGetNonExistentSession(t provider.T) {
	t.Tags("negative", "retrieval")
	t.Description("Test retrieving non-existent session from storage")

	nonExistentID := uuid.New()

	var session *authservice.Session
	var err error
	t.WithNewStep("Attempt to retrieve non-existent session", func(ctx provider.StepCtx) {
		session, err = s.storage.Get(context.Background(), nonExistentID)
	})

	t.WithNewStep("Verify session not found error", func(ctx provider.StepCtx) {
		assert.Nil(t, session)
		assert.Error(t, err)
		assert.Equal(t, authservice.ErrSessionNotFound, err)
	})
}

func (s *MapSessionStorageTestSuite) TestGetExpiredSession(t provider.T) {
	t.Tags("negative", "expiration")
	t.Description("Test retrieving expired session from storage")

	sessionID := uuid.New()
	expiredSession := &authservice.Session{
		MemberID:  uuid.New(),
		ExpiresAt: time.Now().Add(-time.Hour),
	}

	t.WithNewStep("Store expired session", func(ctx provider.StepCtx) {
		err := s.storage.Store(context.Background(), sessionID, expiredSession)
		require.NoError(t, err)
	})

	var session *authservice.Session
	var err error
	t.WithNewStep("Attempt to retrieve expired session", func(ctx provider.StepCtx) {
		session, err = s.storage.Get(context.Background(), sessionID)
	})

	t.WithNewStep("Verify session expired error", func(ctx provider.StepCtx) {
		assert.Nil(t, session)
		assert.Error(t, err)
		assert.Equal(t, authservice.ErrSessionExpired, err)
	})
}

func (s *MapSessionStorageTestSuite) TestDeleteSession(t provider.T) {
	t.Tags("functionality", "deletion")
	t.Description("Test session deletion from storage")

	sessionID := uuid.New()
	session := &authservice.Session{
		MemberID:  uuid.New(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	t.WithNewStep("Store session", func(ctx provider.StepCtx) {
		err := s.storage.Store(context.Background(), sessionID, session)
		require.NoError(t, err)
	})

	t.WithNewStep("Verify session exists before deletion", func(ctx provider.StepCtx) {
		retrieved, err := s.storage.Get(context.Background(), sessionID)
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
	})

	t.WithNewStep("Delete session", func(ctx provider.StepCtx) {
		err := s.storage.Delete(context.Background(), sessionID)
		require.NoError(t, err)
	})

	t.WithNewStep("Verify session deleted", func(ctx provider.StepCtx) {
		retrieved, err := s.storage.Get(context.Background(), sessionID)
		assert.Nil(t, retrieved)
		assert.Error(t, err)
		assert.Equal(t, authservice.ErrSessionNotFound, err)
	})
}

func (s *MapSessionStorageTestSuite) TestDeleteNonExistentSession(t provider.T) {
	t.Tags("functionality", "deletion")
	t.Description("Test deletion of non-existent session")

	nonExistentID := uuid.New()

	t.WithNewStep("Delete non-existent session", func(ctx provider.StepCtx) {
		err := s.storage.Delete(context.Background(), nonExistentID)
		require.NoError(t, err)
	})
}

func (s *MapSessionStorageTestSuite) TestClearExpiredSessions(t provider.T) {
	t.Tags("functionality", "cleanup")
	t.Description("Test clearing expired sessions from storage")

	validSessionID := uuid.New()
	validSession := &authservice.Session{
		MemberID:  uuid.New(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	expiredSessionID := uuid.New()
	expiredSession := &authservice.Session{
		MemberID:  uuid.New(),
		ExpiresAt: time.Now().Add(-time.Hour),
	}

	t.WithNewStep("Store both valid and expired sessions", func(ctx provider.StepCtx) {
		err := s.storage.Store(context.Background(), validSessionID, validSession)
		require.NoError(t, err)
		err = s.storage.Store(context.Background(), expiredSessionID, expiredSession)
		require.NoError(t, err)
	})

	t.WithNewStep("Clear expired sessions", func(ctx provider.StepCtx) {
		err := s.storage.ClearExpired(context.Background())
		require.NoError(t, err)
	})

	t.WithNewStep("Verify expired session removed", func(ctx provider.StepCtx) {
		session, err := s.storage.Get(context.Background(), expiredSessionID)
		assert.Nil(t, session)
		assert.Error(t, err)
	})

	t.WithNewStep("Verify valid session remains", func(ctx provider.StepCtx) {
		session, err := s.storage.Get(context.Background(), validSessionID)
		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, validSession.MemberID, session.MemberID)
	})
}

func (s *MapSessionStorageTestSuite) TestConcurrentAccess(t provider.T) {
	t.Tags("concurrency", "thread-safety")
	t.Description("Test concurrent access to session storage")

	sessionID := uuid.New()
	session := &authservice.Session{
		MemberID:  uuid.New(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	t.WithNewStep("Store session", func(ctx provider.StepCtx) {
		err := s.storage.Store(context.Background(), sessionID, session)
		require.NoError(t, err)
	})

	t.WithNewStep("Perform concurrent reads", func(ctx provider.StepCtx) {
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				retrieved, err := s.storage.Get(context.Background(), sessionID)
				require.NoError(t, err)
				assert.Equal(t, session.MemberID, retrieved.MemberID)
			}()
		}
		wg.Wait()
	})
}

func TestMapSessionStorage(t *testing.T) {
	suite.RunSuite(t, new(MapSessionStorageTestSuite))
}
