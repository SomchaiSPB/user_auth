package hash

import "golang.org/x/crypto/bcrypt"

type StdHasher struct{}

func NewHasher() StdHasher {
	return StdHasher{}
}

func (h StdHasher) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (h StdHasher) CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
