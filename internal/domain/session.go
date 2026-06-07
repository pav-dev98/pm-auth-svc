package domain

import "time"

type Session struct {
	ID           uint
	Token        string
	CredentialID uint
	ExpiresAt    time.Time
	Revoked      bool
}

func NewSession(Token string, CredentialID uint) *Session {
	expiresAt := time.Now().AddDate(0, 1, 0)
	return &Session{
		ID:           0,
		Token:        Token,
		CredentialID: CredentialID,
		ExpiresAt:    expiresAt,
		Revoked:      false,
	}
}

func (r *Session) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

func (r *Session) IsValid() bool {
	return !r.Revoked && !r.IsExpired()
}