package ports

//go:generate mockgen -destination=./mocks/token_service.go -package=mocks . TokenService
type TokenService interface {
	GenerateToken(ID uint, email string) (string, error)
}