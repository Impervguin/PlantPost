package auth

import "net/mail"

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
