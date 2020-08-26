package validation

import (
	"backend_test_RBAC/db"
	"errors"
	"regexp"
	"strings"
)

func ValidateNameUser(UserName string, db *db.DB) (resultNameUser string, err error) {

	addUserName := prepareName(UserName)

	if addUserName == "" {
		return "", errors.New("Введите имя")
	}

	pattern := `^[A-Za-z]+$`
	matched, err := regexp.Match(pattern, []byte(addUserName))
	if matched == false || err != nil {
		return "", errors.New("Имя не может содержать цифры, знаки пунктуации, символы, пробелы")
	}

	exists, err := db.CheckingForPresenceInTheDb(addUserName, "user_name")
	if err != nil {
		return "", err
	}

	if *exists == false {
		return addUserName, nil
	} else {
		return "", errors.New("Пользователь с именем" + " " + addUserName + " " + "уже существует")
	}
}

func prepareName(imputUserName string) (outputUserName string) {
	nameSpaceRemoval := strings.TrimSpace(imputUserName)
	nameLowerCase := strings.ToLower(nameSpaceRemoval)
	formattedUserName := strings.Title(nameLowerCase)
	return formattedUserName
}
