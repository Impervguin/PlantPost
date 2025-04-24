package admins

import (
	"PlantSite/internal/models/auth"
	"PlantSite/internal/repositories/authrepo"
	"context"
	"sync"

	"github.com/google/uuid"
)

var _ authrepo.AdminRepository = (*AdminMap)(nil)

type AdminMap struct {
	ml map[string]auth.Admin
	mu map[uuid.UUID]auth.Admin
	s  sync.RWMutex
}

func NewAdminMap() *AdminMap {
	return &AdminMap{
		ml: make(map[string]auth.Admin),
		mu: make(map[uuid.UUID]auth.Admin),
		s:  sync.RWMutex{},
	}
}

func (m *AdminMap) Add(ctx context.Context, admin auth.Admin) error {
	m.s.Lock()
	defer m.s.Unlock()
	m.ml[admin.Login()] = admin
	m.mu[admin.ID()] = admin
	return nil
}

func (m *AdminMap) Get(ctx context.Context, login string) (auth.Admin, bool) {
	m.s.RLock()
	defer m.s.RUnlock()
	admin, exists := m.ml[login]
	return admin, exists
}

func (m *AdminMap) GetByID(ctx context.Context, id uuid.UUID) (auth.Admin, bool) {
	m.s.RLock()
	defer m.s.RUnlock()
	admin, exists := m.mu[id]
	return admin, exists
}

func (m *AdminMap) Delete(ctx context.Context, login string) error {
	m.s.Lock()
	defer m.s.Unlock()
	admin, exists := m.ml[login]
	if exists {
		delete(m.ml, login)
		delete(m.mu, admin.ID())
	}
	return nil
}

func (m *AdminMap) Iterate(ctx context.Context, fn func(admin auth.Admin) error) error {
	m.s.RLock()
	defer m.s.RUnlock()
	for _, admin := range m.ml {
		if err := fn(admin); err != nil {
			return err
		}
	}
	return nil
}
