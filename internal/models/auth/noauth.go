package auth

type NoAuthUser struct{}

func (u NoAuthUser) HasAuthorRights() bool {
	return false
}

func (u NoAuthUser) HasMemberRights() bool {
	return false
}

func NewNoAuthUser() User {
    return &NoAuthUser{}
}
