package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var _ User = (*Member)(nil)

type Member struct {
	id         uuid.UUID
	name       string
	email      string
	hashPasswd []byte
	createdAt  time.Time
}

func CreateMember(id uuid.UUID, name string, email string, hashPasswd []byte, createdAt time.Time) (*Member, error) {
	member := &Member{
		id:         id,
		name:       name,
		email:      email,
		hashPasswd: hashPasswd,
		createdAt:  createdAt,
	}

	if err := member.Validate(); err != nil {
		return nil, err
	}
	return member, nil
}

func (m *Member) Validate() error {
	if m.name == "" {
		return fmt.Errorf("name should not be empty")
	}
	if m.email == "" {
		return fmt.Errorf("email should not be empty")
	}
	if len(m.name) > MaximumNameLength {
		return fmt.Errorf("name should not exceed %d characters", MaximumNameLength)
	}
	if err := validateEmail(m.email); err != nil {
		return fmt.Errorf("invalid email address: %v", err)
	}
	if m.createdAt.IsZero() {
		return fmt.Errorf("created_at should not be zero")
	}
	if m.hashPasswd == nil {
		return fmt.Errorf("hash_passwd should not be nil")
	}

	return nil
}

func NewMember(name, email string, hashPasswd []byte) (*Member, error) {
	id := uuid.New()
	createdAt := time.Now()
	return CreateMember(id, name, email, hashPasswd, createdAt)
}

func (m *Member) HasAuthorRights() bool {
	return false
}

func (m *Member) HasMemberRights() bool {
	return true
}

func (m *Member) Auth(passwd []byte, authFunc func(hashPasswd []byte, plainPasswd []byte) (bool, error)) bool {
	res, err := authFunc(m.hashPasswd, passwd)
	if err != nil {
		return false
	}
	return res
}

func (m *Member) ID() uuid.UUID {
	return m.id
}

func (m *Member) Name() string {
	return m.name
}

func (m *Member) Email() string {
	return m.email
}

func (m *Member) HashedPassword() []byte {
	return m.hashPasswd
}

func (m *Member) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Member) UpdateName(name string) error {
	previousName := m.name
	m.name = name
	if err := m.Validate(); err != nil {
		m.name = previousName
		return err
	}
	return nil
}

func (m *Member) UpdateEmail(email string) error {
	previousEmail := m.email
	m.email = email
	if err := m.Validate(); err != nil {
		m.email = previousEmail
		return err
	}
	return nil
}

func (m *Member) UpdateHashedPassword(hashPasswd []byte) error {
	previousHashPasswd := m.hashPasswd
	m.hashPasswd = hashPasswd
	if err := m.Validate(); err != nil {
		m.hashPasswd = previousHashPasswd
		return err
	}
	return nil
}

func (m *Member) IsAuthenticated() bool {
	return true
}
