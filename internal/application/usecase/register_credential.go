package usecase

import (
	"errors"

	"github.com/pav-dev98/pm-auth-svc/internal/application/ports"
	"github.com/pav-dev98/pm-auth-svc/internal/domain"
)

type RegisterCredential struct {
	repo           ports.AuthRepository
	passwordHasher ports.PasswordHasher
	tokenService   ports.TokenService
}

func NewRegisterCredential(
	repo ports.AuthRepository,
	passwordHasher ports.PasswordHasher,
	tokenService ports.TokenService,
) *RegisterCredential {
	return &RegisterCredential{
		repo:           repo,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
	}
}

func (uc *RegisterCredential) Execute(email, password string) (string, string, error) {

	existing, err := uc.repo.FindByEmail(email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return "", "", domain.ErrDatabase
	}
	if existing != nil {
		return "", "", domain.ErrDuplicateEmail
	}

	hashedPassword, err := uc.passwordHasher.Hash(password)
	if err != nil {
		return "", "", domain.ErrPasswordHash
	}

	cred := domain.NewAuthCredential(email, hashedPassword)

	err = uc.repo.Create(cred)
	if err != nil {
		return "", "", err
	}

	token, err := uc.tokenService.GenerateToken(cred.ID, cred.Email)
	if err != nil {
		return "", "", err
	}

	return token, "", nil
}
