package auth

import (
	"fmt"
	"time"
)

type Author struct {
	member     Member
	rights     bool
	giveTime   time.Time
	revokeTime time.Time
}

func CreateAuthor(member Member, giveTime time.Time, rights bool, revokeTime time.Time) (*Author, error) {
	ath := &Author{
		member:     member,
		rights:     true,
		giveTime:   giveTime,
		revokeTime: time.Time{},
	}
	if err := ath.Validate(); err != nil {
		return nil, err
	}

	return ath, nil
}

func (a *Author) Validate() error {
	if err := a.member.Validate(); err != nil {
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
