package ports

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed string, plain string) error
}