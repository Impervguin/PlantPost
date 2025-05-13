package auth

import (
	"fmt"

	"github.com/google/uuid"
)

var _ User = (*Admin)(nil)

type Admin struct {
	id           uuid.UUID
	login        string
	hashPassword []byte
}

func CreateAdmin(id uuid.UUID, login string, hashPassword []byte) (*Admin, error) {
	admin := &Admin{
		id:           id,
		login:        login,
		hashPassword: hashPassword,
	}
	if err := admin.Validate(); err != nil {
		return nil, err
	}
	return admin, nil
}

func (a *Admin) Validate() error {
	if a.id == uuid.Nil {
		return fmt.Errorf("id should not be nil")
	}
	if a.login == "" {
		return fmt.Errorf("login should not be empty")
	}
	if a.hashPassword == nil {
		return fmt.Errorf("hash_password should not be empty")
	}
	return nil
}

func NewAdmin(login string, hashPassword []byte) (*Admin, error) {
	id := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(login))
	return CreateAdmin(id, login, hashPassword)
}

func (a *Admin) HasAuthorRights() bool {
	return true
}

func (a *Admin) HasMemberRights() bool {
	return true
}

func (a *Admin) Auth(passwd []byte, authFunc func(hashPasswd []byte, plainPasswd []byte) (bool, error)) bool {
	res, err := authFunc([]byte(a.hashPassword), passwd)
	if err != nil {
		return false
	}
	return res
}

func (a *Admin) ID() uuid.UUID {
	return a.id
}

func (a *Admin) Login() string {
	return a.login
}

func (a *Admin) HashedPassword() []byte {
	return []byte(a.hashPassword)
}

func (a *Admin) IsAuthenticated() bool {
	return true
}

func (a *Admin) Username() string {
	return a.login
}
