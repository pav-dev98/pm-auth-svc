package usecase_test

import (
	"errors"
	"testing"

	"github.com/pav-dev98/pm-auth-svc/internal/application/ports/mocks"
	"github.com/pav-dev98/pm-auth-svc/internal/application/usecase"
	"github.com/pav-dev98/pm-auth-svc/internal/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLoginCredential_Execute(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repository failure")
	compareErr := errors.New("compare failure")
	tokenErr := errors.New("token failure")

	tests := []struct {
		name          string
		email         string
		password      string
		setupMocks    func(*mocks.MockAuthRepository, *mocks.MockPasswordHasher, *mocks.MockTokenService)
		expectedToken string
		expectedErr   error
	}{
		{
			name:     "returns token when credentials are valid",
			email:    "user@example.com",
			password: "plain-password",
			setupMocks: func(repo *mocks.MockAuthRepository, hasher *mocks.MockPasswordHasher, tokenService *mocks.MockTokenService) {
				cred := &domain.AuthCredential{ID: 42, Email: "user@example.com", Password: "hashed-password"}

				repo.EXPECT().FindByEmail("user@example.com").Return(cred, nil)
				hasher.EXPECT().Compare("hashed-password", "plain-password").Return(nil)
				tokenService.EXPECT().GenerateToken(uint(42), "user@example.com").Return("jwt-token", nil)
			},
			expectedToken: "jwt-token",
		},
		{
			name:     "returns repository error when find by email fails",
			email:    "user@example.com",
			password: "plain-password",
			setupMocks: func(repo *mocks.MockAuthRepository, hasher *mocks.MockPasswordHasher, tokenService *mocks.MockTokenService) {
				repo.EXPECT().FindByEmail("user@example.com").Return(nil, repoErr)
			},
			expectedErr: repoErr,
		},
		{
			name:     "returns invalid credentials when credential is not found",
			email:    "missing@example.com",
			password: "plain-password",
			setupMocks: func(repo *mocks.MockAuthRepository, hasher *mocks.MockPasswordHasher, tokenService *mocks.MockTokenService) {
				repo.EXPECT().FindByEmail("missing@example.com").Return(nil, nil)
			},
			expectedErr: errors.New("invalid credentials"),
		},
		{
			name:     "returns invalid credentials when password comparison fails",
			email:    "user@example.com",
			password: "wrong-password",
			setupMocks: func(repo *mocks.MockAuthRepository, hasher *mocks.MockPasswordHasher, tokenService *mocks.MockTokenService) {
				cred := &domain.AuthCredential{ID: 42, Email: "user@example.com", Password: "hashed-password"}

				repo.EXPECT().FindByEmail("user@example.com").Return(cred, nil)
				hasher.EXPECT().Compare("hashed-password", "wrong-password").Return(compareErr)
			},
			expectedErr: errors.New("invalid credentials"),
		},
		{
			name:     "returns token service error when token generation fails",
			email:    "user@example.com",
			password: "plain-password",
			setupMocks: func(repo *mocks.MockAuthRepository, hasher *mocks.MockPasswordHasher, tokenService *mocks.MockTokenService) {
				cred := &domain.AuthCredential{ID: 42, Email: "user@example.com", Password: "hashed-password"}

				repo.EXPECT().FindByEmail("user@example.com").Return(cred, nil)
				hasher.EXPECT().Compare("hashed-password", "plain-password").Return(nil)
				tokenService.EXPECT().GenerateToken(uint(42), "user@example.com").Return("", tokenErr)
			},
			expectedErr: tokenErr,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mocks.NewMockAuthRepository(ctrl)
			hasher := mocks.NewMockPasswordHasher(ctrl)
			tokenService := mocks.NewMockTokenService(ctrl)

			tt.setupMocks(repo, hasher, tokenService)

			uc := usecase.NewLoginCredential(repo, hasher, tokenService)

			token, err := uc.Execute(tt.email, tt.password)

			assert.Equal(t, tt.expectedToken, token)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
