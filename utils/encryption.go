package utils

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (string, error) {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(encryptedPass), err
}

func ComparePassword(password string, encryptedPass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPass), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
