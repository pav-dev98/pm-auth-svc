package postgress

import (
	"errors"
	"fmt"
	"time"

	"github.com/pav-dev98/pm-auth-svc/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AuthCredentialModel struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex"`
	Password string
	Role     string
	IsActive bool
}

type SessionModel struct {
	ID 		 		uint `gorm:"primaryKey"`
	Token 	 		string `gorm:"uniqueIndex"`
	CredentialID 	uint `gorm:"not null"`
	ExpiresAt    	time.Time `gorm:"not null"`
	Revoked       	bool `gorm:"not null"`
}

func toDomain(m *AuthCredentialModel) *domain.AuthCredential {
	return &domain.AuthCredential{
		Email:    m.Email,
		Password: m.Password,
		Role:     m.Role,
		IsActive: m.IsActive,
	}
}

func toModel(d *domain.AuthCredential) *AuthCredentialModel {
	return &AuthCredentialModel{
		Email:    d.Email,
		Password: d.Password,
		Role:     d.Role,
		IsActive: d.IsActive,
	}
}

func toDomainSession(m *SessionModel) *domain.Session {
	return &domain.Session{
		ID: m.ID,
		Token: m.Token,
		CredentialID: m.CredentialID,
		ExpiresAt: m.ExpiresAt,
		Revoked: m.Revoked,
	}
}

func toModelSession(d *domain.Session) *SessionModel {
	return &SessionModel{
		Token: d.Token,
		CredentialID: d.CredentialID,
		ExpiresAt: d.ExpiresAt,
		Revoked: d.Revoked,
	}
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(dsn string) (*AuthRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate crea la tabla si no existe
	db.AutoMigrate(&AuthCredentialModel{}, &SessionModel{})

	return &AuthRepository{db: db}, nil
}

func (r *AuthRepository) Create(cred *domain.AuthCredential) error {
	model := toModel(cred)
	return r.db.Create(&model).Error
}

func (r *AuthRepository) FindByEmail(email string) (*domain.AuthCredential, error) {
	var model AuthCredentialModel

	err := r.db.Where("email = ?", email).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, domain.ErrDatabase
	}

	fmt.Println("User found:", model.Email)

	return toDomain(&model), nil
}

func (r * AuthRepository) SaveSession(session *domain.Session) error {
	model := toModelSession(session)
	return r.db.Create(&model).Error
} 

func (r * AuthRepository) FindSession(token string) (*domain.Session , error) {
	var model SessionModel

	err := r.db.Where("token = ?",token).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil, domain.ErrNotFound
		}
		return nil, domain.ErrDatabase
	}
	fmt.Println("Token found", model.Token)
	return toDomainSession(&model),nil
}

func (r *AuthRepository) RevokeSession(token string) error {
	result := r.db.Model(&SessionModel{}).
		Where("token = ?", token).
		Update("revoked", true)

	if result.Error != nil {
		return domain.ErrDatabase
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}