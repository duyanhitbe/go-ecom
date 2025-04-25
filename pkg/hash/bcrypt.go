package hash

import (
	"github.com/duyanhitbe/go-ecom/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct {
}

func NewBcrypt() *Bcrypt {
	return &Bcrypt{}
}

func (b *Bcrypt) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (b *Bcrypt) Verify(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if utils.IsErrMismatchedPassword(err) {
			return false, nil
		}

		return false, err
	}
	return true, nil
}
