package hash

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// Warning: High levels can be slooooowwwwwww!
	// (2 ^ cost time)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
