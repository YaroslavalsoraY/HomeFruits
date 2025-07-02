package hashfunc

import "golang.org/x/crypto/bcrypt"

func HashingPassword(rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}

	return string(hashedPassword), nil
}

func HashCompareWithPassw(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}

	return true
}