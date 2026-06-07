package ports

import "github.com/pav-dev98/pm-auth-svc/internal/domain"

//go:generate mockgen -destination=./mocks/auth_repository.go -package=mocks . AuthRepository
type AuthRepository interface {
	Create(cred *domain.AuthCredential) error
	FindByEmail(email string) (*domain.AuthCredential, error)
	SaveSession(session *domain.Session) error
	FindSession(token string) (*domain.Session, error)
	RevokeSession(token string) error
}