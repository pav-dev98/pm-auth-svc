package usecase

import (
	"github.com/pav-dev98/pm-auth-svc/internal/application/ports"
	"github.com/pav-dev98/pm-auth-svc/internal/domain"
)

type RefreshCredential struct {
	authRepository ports.AuthRepository
	tokenService   ports.TokenService
}

func NewRefreshCredential(
	authRepository ports.AuthRepository,
	tokenService ports.TokenService,
) *RefreshCredential {
	return &RefreshCredential{
		authRepository: authRepository,
		tokenService:   tokenService,
	}
}

func (uc *RefreshCredential) Execute(refreshToken string) (string, error) {
	// 1. Buscar el refresh token en DB
	session, err := uc.authRepository.FindSession(refreshToken)
	if err != nil {
		return "", domain.ErrNotFound
	}

	// 2. Validar que no esté expirado ni revocado
	if !session.IsValid() {
		return "", domain.ErrInvalidToken
	}

	// 3. Generar nuevo access token
	newAccessToken, err := uc.tokenService.GenerateToken(session.CredentialID, "")
	if err != nil {
		return "", domain.ErrInternal
	}

	return newAccessToken, nil
}
