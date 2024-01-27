package util

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword は、与えられたパスワードを bcrypt を使ってハッシュ化する
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHashは、生のパスワードとそのハッシュ化されたバージョンを比較する
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
