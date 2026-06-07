package usecase

import (
	"errors"

	"github.com/pav-dev98/pm-auth-svc/internal/application/ports"
	"github.com/pav-dev98/pm-auth-svc/internal/domain"
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

func (uc *LoginCredential) Execute(email, password string) (string, string, error) {
	cred, err := uc.repo.FindByEmail(email)
	if err != nil {
		return "", "", err
	}
	if cred == nil {
		return "", "", errors.New("invalid credentials")
	}

	err = uc.passwordChecker.Compare(cred.Password, password)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := uc.tokenService.GenerateToken(cred.ID, cred.Email) // falta agregar roles , permisos
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uc.tokenService.GenerateRefreshToken(cred.ID)
	if err != nil {
		return "", "", err
	}
	session := domain.NewSession(refreshToken, cred.ID)

	err = uc.repo.SaveSession(session)
	if err != nil {
		return "", "", domain.ErrDatabase
	}

	return accessToken, refreshToken, nil
}

