package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// return the bcrypt hash string of input password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// check if correct password is given during authentication
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(hashedPassword))
}
