package crypto

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const defaultPasswordLength = 6

func IsGoodPassword(rawPassword string) bool {
	return len(rawPassword) >= defaultPasswordLength
}

func HashAndSalt(pwd string) string {
	pwdBytes := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(pwdBytes, bcrypt.MinCost)
	if err != nil {
		logrus.Error(err)
	}
	return string(hash)
}

func ComparePasswords(hashedPwd, plainPwd string) bool {
	pwdBytes := []byte(plainPwd)

	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, pwdBytes)
	if err != nil {
		logrus.Println(err)
		return false
	}

	return true
}
