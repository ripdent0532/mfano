package secret

import "golang.org/x/crypto/bcrypt"

func EncodePassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	encodePassword := string(hash)
	return encodePassword
}

func CheckPassword(encodePassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodePassword), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}
