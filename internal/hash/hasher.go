package hash

type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hashedPassword string) bool
}
