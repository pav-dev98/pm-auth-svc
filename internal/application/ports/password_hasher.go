package ports

//go:generate mockgen -destination=./mocks/password_hasher.go -package=mocks . PasswordHasher
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed string, plain string) error
}