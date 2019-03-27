package postgresql

import "golang.org/x/crypto/bcrypt"

func HashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	switch err {
	case bcrypt.ErrMismatchedHashAndPassword:
		{
			return false, nil
		}
	case nil:
		{
			return true, nil
		}
	default:
		{
			return false, err
		}
	}
}
