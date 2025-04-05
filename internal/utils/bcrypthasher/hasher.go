package bcrypthasher

import (
	authservice "PlantSite/internal/services/auth-service"

	"golang.org/x/crypto/bcrypt"
)

type bcrypthasher struct {
	cost int
}

func NewBcryptHasher(cost int) authservice.PasswdHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		panic("bcrypt cost out of range")
	}
	return &bcrypthasher{}
}

func (b *bcrypthasher) Hash(passwd []byte) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword(passwd, b.cost)
	return hashed, err
}

func (b *bcrypthasher) Compare(hashedPasswd, plainPasswd []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPasswd, plainPasswd)
	return err == nil, err
}
