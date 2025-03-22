package auth

type User interface {
	HasAuthorRights() bool
	HasMemberRights() bool
}
