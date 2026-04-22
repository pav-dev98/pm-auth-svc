package usecase

import (
	"github.com/pav-dev98/pm-auth-svc/internal/application/ports"
	"errors"
)

type LoginCredential struct {
	repo            ports.AuthRepository
	passwordChecker ports.PasswordHasher
	tokenService    ports.TokenService
}

func NewLoginCredential(
	repo ports.AuthRepository,
	passwordChecker ports.PasswordHasher,
	tokenService ports.TokenService,
) *LoginCredential {
	return &LoginCredential{
		repo:            repo,
		passwordChecker: passwordChecker,
		tokenService:    tokenService,
	}
}

func (uc *LoginCredential) Execute(email, password string) (string, error) {

	cred, err := uc.repo.FindByEmail(email)
	
	if err != nil {
		return "", err
	}
	if cred == nil {
		return "", errors.New("invalid credentials")
	}

	err = uc.passwordChecker.Compare(cred.Password, password)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := uc.tokenService.GenerateToken(cred.ID,cred.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}