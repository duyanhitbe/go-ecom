package utils

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func IsErrNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func IsErrMismatchedPassword(err error) bool {
	return errors.Is(err, bcrypt.ErrMismatchedHashAndPassword)
}
