package ports

import "github.com/pav-dev98/pm-auth-svc/internal/domain"

type AuthRepository interface {
	Create(cred *domain.AuthCredential) error
	FindByEmail(email string) (*domain.AuthCredential, error)
}