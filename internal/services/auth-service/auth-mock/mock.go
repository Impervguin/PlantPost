package authmock

import (
	"PlantSite/internal/models/auth"
	authservice "PlantSite/internal/services/auth-service"
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockSessionStorage implements SessionStorage interface
type MockSessionStorage struct {
	mock.Mock
}

func (m *MockSessionStorage) Store(ctx context.Context, sid uuid.UUID, session *authservice.Session) error {
	args := m.Called(ctx, sid, session)
	return args.Error(0)
}

func (m *MockSessionStorage) Get(ctx context.Context, sid uuid.UUID) (*authservice.Session, error) {
	args := m.Called(ctx, sid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authservice.Session), args.Error(1)
}

func (m *MockSessionStorage) Delete(ctx context.Context, sid uuid.UUID) error {
	args := m.Called(ctx, sid)
	return args.Error(0)
}

func (m *MockSessionStorage) ClearExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockAuthRepository implements auth.AuthRepository interface
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) Get(ctx context.Context, id uuid.UUID) (auth.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(auth.User), args.Error(1)
}

func (m *MockAuthRepository) Create(ctx context.Context, mem *auth.Member) (auth.User, error) {
	args := m.Called(ctx, mem)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(auth.User), args.Error(1)
}

func (m *MockAuthRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(auth.User) (auth.User, error)) (auth.User, error) {
	args := m.Called(ctx, id, updateFn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(auth.User), args.Error(1)
}

func (m *MockAuthRepository) GetByName(ctx context.Context, name string) (auth.User, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(auth.User), args.Error(1)
}

func (m *MockAuthRepository) GetByEmail(ctx context.Context, email string) (auth.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(auth.User), args.Error(1)
}

// MockPasswdHasher implements PasswdHasher interface
type MockPasswdHasher struct {
	mock.Mock
}

func (m *MockPasswdHasher) Hash(password []byte) ([]byte, error) {
	args := m.Called(password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockPasswdHasher) Compare(hashPasswd, plainPasswd []byte) (bool, error) {
	args := m.Called(hashPasswd, plainPasswd)
	return args.Bool(0), args.Error(1)
}

// MockUser implements auth.User interface
type MockUser struct {
	mock.Mock
}

var _ auth.User = &MockUser{}

func (m *MockUser) ID() uuid.UUID {
	return m.Called().Get(0).(uuid.UUID)
}

func (m *MockUser) Name() string {
	return m.Called().String(0)
}

func (m *MockUser) Email() string {
	return m.Called().String(0)
}

func (m *MockUser) HashedPassword() []byte {
	return m.Called().Get(0).([]byte)
}

func (m *MockUser) Auth(password []byte, compareFn func(hashPasswd, plainPasswd []byte) (bool, error)) bool {
	args := m.Called(password, compareFn)
	return args.Bool(0)
}

func (m *MockUser) HasMemberRights() bool {
	return m.Called().Bool(0)
}

func (m *MockUser) HasAuthorRights() bool {
	return m.Called().Bool(0)
}

func (m *MockUser) IsAuthenticated() bool {
	return m.Called().Bool(0)
}

func (m *MockUser) Username() string {
	return m.Called().String(0)
}
