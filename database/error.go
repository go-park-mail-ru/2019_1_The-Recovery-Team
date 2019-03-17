package database

import "errors"

var (
	ErrConnRefused           = errors.New("ConnectionRefused")
	ErrEmailAlreadyExists    = errors.New("EmailAlreadyExists")
	ErrNicknameAlreadyExists = errors.New("NicknameAlreadyExists")
	ErrIncorrectPassword     = errors.New("IncorrectProfilePassword")
)
