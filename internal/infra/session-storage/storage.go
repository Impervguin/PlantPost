package sessionstorage

import (
	authservice "PlantSite/internal/services/auth-service"
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MapSessionStorage struct {
	storage map[uuid.UUID]*authservice.Session
	mutex   sync.RWMutex
}

func NewMapSessionStorage() *MapSessionStorage {
	return &MapSessionStorage{
		storage: make(map[uuid.UUID]*authservice.Session),
		mutex:   sync.RWMutex{},
	}
}

func (storage *MapSessionStorage) Get(ctx context.Context, sid uuid.UUID) (*authservice.Session, error) {
	storage.mutex.RLock()
	session, ok := storage.storage[sid]
	storage.mutex.RUnlock()
	if !ok {
		return nil, authservice.ErrSessionNotFound
	}
	if session.ExpiresAt.Before(time.Now()) {
		storage.Delete(ctx, sid)
		return nil, authservice.ErrSessionExpired
	}
	return session, nil
}

func (storage *MapSessionStorage) Store(ctx context.Context, sid uuid.UUID, session *authservice.Session) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	storage.storage[sid] = session
	return nil
}

func (storage *MapSessionStorage) Delete(ctx context.Context, sid uuid.UUID) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	delete(storage.storage, sid)
	return nil
}

func (storage *MapSessionStorage) ClearExpired(ctx context.Context) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	for sid, session := range storage.storage {
		if session.ExpiresAt.Before(time.Now()) {
			delete(storage.storage, sid)
		}
	}
	return nil
}
