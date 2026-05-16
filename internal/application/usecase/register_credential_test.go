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

func TestRegisterCredential_Execute(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repository failure")
	createErr := errors.New("create failure")
	hashErr := errors.New("hash failure")
	tokenErr := errors.New("token failure")
	refreshTokenErr := errors.New("refresh token failure")

	tests := []struct {
		name            string
		email           string
		password        string
		setupMocks      func(*mocks.MockAuthRepository, *mocks.MockPasswordHasher, *mocks.MockTokenService)
		expectedToken   string
		expectedRefresh string
		expectedErr     error
	}{
		{
			name:     "returns access and refresh token when registration succeeds",
			email:    "user@example.com",
			password: "plain-password",
			setupMocks: func(repo *mocks.MockAuthRepository, hasher *mocks.MockPasswordHasher, tokenService *mocks.MockTokenService) {
				repo.EXPECT().FindByEmail("user@example.com").Return(nil, domain.ErrNotFound)
				hasher.EXPECT().Hash("plain-password").Return("hashed-password", nil)
				repo.EXPECT().Create(gomock.AssignableToTypeOf(&domain.AuthCredential{})).DoAndReturn(func(cred *domain.AuthCredential) error {
					assert.Equal(t, "user@example.com", cred.Email)
					assert.Equal(t, "hashed-password", cred.Password)
					assert.Equal(t, "user", cred.Role)
					assert.True(t, cred.IsActive)
					return nil
				})
				tokenService.EXPECT().GenerateToken(uint(0), "user@example.com").Return("jwt-access-token", nil)
				tokenService.EXPECT().GenerateRefreshToken(uint(0)).Return("jwt-refresh-token", nil)
			},
			expectedToken:   "jwt-access-token",
			expectedRefresh: "jwt-refresh-token",
		},
		{
			name:     "returns database error when find by email fails unexpectedly",
			email:    "user@example.com",
			password: "plain-password",
			setupMocks: func(repo *mocks.MockAuthRepository, hasher *mocks.MockPasswordHasher, tokenService *mocks.MockTokenService) {
				repo.EXPECT().FindByEmail("user@example.com").Return(nil, repoErr)
			},
			expectedErr: domain.ErrDatabase,
		},
		{
			name:     "returns duplicate email error when credential already exists",
			email:    "user@example.com",
			password: "plain-password",
			setupMocks: func(repo *mocks.MockAuthRepository, hasher *mocks.MockPasswordHasher, tokenService *mocks.MockTokenService) {
				existing := &domain.AuthCredential{ID: 1, Email: "user@example.com"}
				repo.EXPECT().FindByEmail("user@example.com").Return(existing, nil)
			},
			expectedErr: domain.ErrDuplicateEmail,
		},
		{
			name:     "returns password hash error when hashing fails",
			email:    "user@example.com",
			password: "plain-password",
			setupMocks: func(repo *mocks.MockAuthRepository, hasher *mocks.MockPasswordHasher, tokenService *mocks.MockTokenService) {
				repo.EXPECT().FindByEmail("user@example.com").Return(nil, domain.ErrNotFound)
				hasher.EXPECT().Hash("plain-password").Return("", hashErr)
			},
			expectedErr: domain.ErrPasswordHash,
		},
		{
			name:     "returns create error when repository create fails",
			email:    "user@example.com",
			password: "plain-password",
			setupMocks: func(repo *mocks.MockAuthRepository, hasher *mocks.MockPasswordHasher, tokenService *mocks.MockTokenService) {
				repo.EXPECT().FindByEmail("user@example.com").Return(nil, domain.ErrNotFound)
				hasher.EXPECT().Hash("plain-password").Return("hashed-password", nil)
				repo.EXPECT().Create(gomock.AssignableToTypeOf(&domain.AuthCredential{})).Return(createErr)
			},
			expectedErr: createErr,
		},
		{
			name:     "returns error when access token generation fails",
			email:    "user@example.com",
			password: "plain-password",
			setupMocks: func(repo *mocks.MockAuthRepository, hasher *mocks.MockPasswordHasher, tokenService *mocks.MockTokenService) {
				repo.EXPECT().FindByEmail("user@example.com").Return(nil, domain.ErrNotFound)
				hasher.EXPECT().Hash("plain-password").Return("hashed-password", nil)
				repo.EXPECT().Create(gomock.AssignableToTypeOf(&domain.AuthCredential{})).Return(nil)
				tokenService.EXPECT().GenerateToken(uint(0), "user@example.com").Return("", tokenErr)
			},
			expectedErr: tokenErr,
		},
		{
			name:     "returns error when refresh token generation fails",
			email:    "user@example.com",
			password: "plain-password",
			setupMocks: func(repo *mocks.MockAuthRepository, hasher *mocks.MockPasswordHasher, tokenService *mocks.MockTokenService) {
				repo.EXPECT().FindByEmail("user@example.com").Return(nil, domain.ErrNotFound)
				hasher.EXPECT().Hash("plain-password").Return("hashed-password", nil)
				repo.EXPECT().Create(gomock.AssignableToTypeOf(&domain.AuthCredential{})).Return(nil)
				tokenService.EXPECT().GenerateToken(uint(0), "user@example.com").Return("jwt-access-token", nil)
				tokenService.EXPECT().GenerateRefreshToken(uint(0)).Return("", refreshTokenErr)
			},
			expectedErr: refreshTokenErr,
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

			uc := usecase.NewRegisterCredential(repo, hasher, tokenService)

			token, refreshToken, err := uc.Execute(tt.email, tt.password)

			assert.Equal(t, tt.expectedToken, token)
			assert.Equal(t, tt.expectedRefresh, refreshToken)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
