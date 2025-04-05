package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Author struct {
	Member
	rights     bool
	giveTime   time.Time
	revokeTime time.Time
}

func CreateAuthor(member Member, giveTime time.Time, rights bool, revokeTime time.Time) (*Author, error) {
	ath := &Author{
		Member:     member,
		rights:     true,
		giveTime:   giveTime,
		revokeTime: revokeTime,
	}
	if err := ath.Validate(); err != nil {
		return nil, err
	}

	return ath, nil
}

func (a *Author) Validate() error {
	if err := a.Member.Validate(); err != nil {
		return err
	}
	if a.giveTime.IsZero() {
		return fmt.Errorf("give time must be non-zero")
	}
	if a.revokeTime.IsZero() {
		return fmt.Errorf("revoke time must be non-zero")
	}
	if a.revokeTime.Before(a.giveTime) {
		return fmt.Errorf("revoke time must be after give time")
	}

	return nil
}

func (a *Author) RevokeAuthorRights() {
	a.rights = false
	a.revokeTime = time.Now()
}

func (a *Author) HasAuthorRights() bool {
	return a.rights
}

func (a *Author) HasMemberRights() bool {
	return true
}

func (a *Author) Auth(passwd []byte, authFunc func(hashPasswd []byte, plainPasswd []byte) (bool, error)) bool {
	return a.Member.Auth(passwd, authFunc)
}

func (a *Author) ID() uuid.UUID {
	return a.Member.ID()
}

func (a *Author) HasRights() bool {
	return a.rights
}

func (a *Author) GiveTime() time.Time {
	return a.giveTime
}

func (a *Author) RevokeTime() time.Time {
	return a.revokeTime
}
