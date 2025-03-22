package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

func NewMember(name, email string, password string) (*Member, error) {
	id := uuid.New()
	hashPasswd, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return nil, err
	}
	createdAt := time.Now()
	return CreateMember(id, name, email, hashPasswd, createdAt)
}

func (m *Member) HasAuthorRights() bool {
	return false
}

func (m *Member) HasMemberRights() bool {
	return true
}
