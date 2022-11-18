package password

import "golang.org/x/crypto/bcrypt"

func EncodePassword(rawPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidatePassword(encodePassword, raw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodePassword), []byte(raw))
	return err == nil
}
