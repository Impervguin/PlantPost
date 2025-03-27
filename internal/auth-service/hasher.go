package authservice

type PasswdHasher interface {
	Hash(passwd []byte) ([]byte, error)
	Compare(hashedPasswd, plainPasswd []byte) (bool, error)
}
