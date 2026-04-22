package ports

type TokenService interface {
	GenerateToken(ID uint, email string) (string, error)
}