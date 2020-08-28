package verification

import (
	"backend_test_RBAC/db"
	"backend_test_RBAC/hashing"
	"errors"
)

func VerificationLogin(login, passUser string, db *db.DB) error {

	pass, err := db.PasswordCheck(login)
	if err != nil {
		return err
	}
	hashPass, err := hashing.ValidateAndHashPasswordUser(passUser)
	if err != nil {
		return err
	}

	if *pass != hashPass {
		return errors.New("Введен неверный логин или пароль")
	}
	return nil
}
