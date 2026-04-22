package postgress

import (
	"github.com/pav-dev98/pm-auth-svc/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
	"errors"
)

type AuthCredentialModel struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex"`
	Password string
	Role     string
	IsActive bool
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

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(dsn string) (*AuthRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate crea la tabla si no existe
	db.AutoMigrate(&AuthCredentialModel{})

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
