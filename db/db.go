package db

import (
	"backend_test_RBAC/hashing"
	"fmt"
	_ "github.com/lib/pq"
)

func (db *DB) СreateTable() error {

	ins := "CREATE TABLE IF NOT EXISTS the_users (id SERIAL, user_name VARCHAR, role VARCHAR, password VARCHAR)"
	_, err := db.Connection.Exec(ins)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateAdministrator(adminName, role, adminPass string) error {

	password, err := hashing.ValidateAndHashPasswordUser(adminPass)
	if err != nil {
		return err
	}

	ins := `INSERT INTO the_users (user_name, role, password )
	SELECT user_name, role, password FROM the_users UNION  VALUES ($1, $2, $3)
	except SELECT user_name, role, password FROM the_users`

	_, err = db.Connection.Exec(ins, adminName, role, password)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CheckingForPresenceInTheDb(inputParameter, nameColumn string) (*bool, error) {
	rows, err := db.Connection.Query(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM the_users WHERE %s = '%s')", nameColumn, inputParameter))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exists bool
	for rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &exists, nil
}

func (db *DB) AddNewUser(userName, role, PasswordUser string) (err error) {

	ins := "INSERT INTO the_users (user_name, role, password) VALUES ($1, $2, $3)"
	_, err = db.Connection.Exec(ins, userName, role, PasswordUser)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) PasswordCheck(login string) (string, error) {

	rows, err := db.Connection.Query(fmt.Sprintf("SELECT password  FROM the_users WHERE user_name = '%s'", login))
	if err != nil {
		return "", err
	}

	var password string
	for rows.Next() {
		err := rows.Scan(&password)
		if err != nil {
			return "", err
		}
	}
	return password, nil
}
