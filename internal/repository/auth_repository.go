package repository

import (
	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// AuthCredantial modelo de la tabla
type AuthCredantial struct {
	ID                    uint      `gorm:"primaryKey;autoIncrement"`
	Email                 string    `gorm:"uniqueIndex;not null" validate:"required,email"`
	Password              string    `gorm:"not null" validate:"required,min=8"`
	Role                  string    `gorm:"not null;default:'user'" validate:"required,oneof=user admin"`
	IsOnboardingComplete  bool      `gorm:"not null;default:false"`
	IsActive              bool      `gorm:"not null;default:true"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
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
	db.AutoMigrate(&AuthCredantial{})

	return &AuthRepository{db: db}, nil
}

// FindByEmail busca una credencial por email
func (r *AuthRepository) FindByEmail(email string) (*AuthCredantial, error) {
	var authCredantial AuthCredantial
	result := r.db.Where("email = ?", email).First(&authCredantial)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // no existe
		}
		return nil, result.Error
	}
	return &authCredantial, nil
}

// Create guarda una nueva credencial
func (r *AuthRepository) Create(email, hashedPassword string) (*AuthCredantial, error) {
	authCredantial := &AuthCredantial{
		Email:                email,
		Password:             hashedPassword,
		Role:                 "user",
		IsOnboardingComplete: false,
		IsActive:             true,
	}
	result := r.db.Create(authCredantial)
	return authCredantial, result.Error
}