package model

import (
	"backend_test_RBAC/db"
	"errors"
	"strconv"
)

type User struct {
	ID   int
	Name string
	Role string
}

func (u User) Exists(id int, db *db.DB) (*bool, error) {

	exists, err := db.CheckingForPresenceInTheDb(strconv.Itoa(id), "id")
	if err != nil {
		return nil, err
	}

	return exists, nil
}

func (u User) FindByName(name string, db *db.DB) (*User, error) {

	query := "SELECT id, user_name, role FROM the_users WHERE user_name = " + " " + "'" + name + "'"
	rows, err := db.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.Name, &u.Role); err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if u.Name == "" {
		return nil, errors.New("ПОЛЬЗОВАТЕЛЬ НЕ НАЙДЕН")
	}

	return &u, nil
}
