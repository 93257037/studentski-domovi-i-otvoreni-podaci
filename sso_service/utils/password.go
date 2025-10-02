package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// hesuje lozinku koristeci bcrypt algoritam sa default jacinom
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// proverava da li se lozinka slaze sa hesovanom verzijom
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

