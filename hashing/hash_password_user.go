package hashing

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"
	"regexp"
	"strings"
)

func ValidateAndHashPasswordUser(imputpPassword string) (hashPassword string, err error) {
	password := strings.TrimSpace(imputpPassword)

	if password == "" {
		return "", errors.New("Не задан пароль")
	}

	pattern := `^[A-Za-z0-9_-]{6,25}$`
	matched, err := regexp.Match(pattern, []byte(password))
	if matched == false || err != nil {
		return "", errors.New("Пароль должен состоять хотя бы из 6 символов, может содержать буквы, цифры, знаки -, _")
	}

	salt := md5.Sum([]byte(os.Getenv("SALT")))
	saltHash := hex.EncodeToString(salt[:])
	hashPasswordSalt := password + saltHash

	hash := md5.Sum([]byte(hashPasswordSalt))
	return hex.EncodeToString(hash[:]), nil
}
