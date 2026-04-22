package domain

type AuthCredential struct {
	ID uint
	Email string
	Password string
	Role string
	IsActive bool
}

func NewAuthCredential(email,password string) *AuthCredential{
	return &AuthCredential{
		ID:       0,
		Email:    email,
		Password: password,
		Role:     "user",
		IsActive: true,
	}
}