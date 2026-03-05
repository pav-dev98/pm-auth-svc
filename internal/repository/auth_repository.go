package repository

import (
	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User modelo de la tabla
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
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
	db.AutoMigrate(&User{})

	return &AuthRepository{db: db}, nil
}

// FindByEmail busca un usuario por email
func (r *AuthRepository) FindByEmail(email string) (*User, error) {
	var user User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // no existe
		}
		return nil, result.Error
	}
	return &user, nil
}

// Create guarda un nuevo usuario
func (r *AuthRepository) Create(email, hashedPassword string) (*User, error) {
	user := &User{
		Email:    email,
		Password: hashedPassword,
	}
	result := r.db.Create(user)
	return user, result.Error
}