package models

import (
	"time"

	"github.com/pgvector/pgvector-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           uint            `json:"id" gorm:"primaryKey;autoIncrement;index" swaggerignore:"true"`
	Username     string          `json:"username" gorm:"unique;not null"`
	Email        string          `json:"email" gorm:"unique;not null"`
	Password     string          `json:"-" gorm:"not null" query:"-"`
	PasswordTemp string          `json:"password" gorm:"-" query:"-"`
	Embedding    pgvector.Vector `json:"-" gorm:"type:vector(768)"`
	CreatedAt    time.Time       `json:"created_at" gorm:"autoCreateTime" swaggerignore:"true"`
	UpdatedAt    time.Time       `json:"updated_at" gorm:"autoUpdateTime" swaggerignore:"true"`
	DeletedAt    gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.PasswordTemp != "" {
		hashedPassword, err := HashPassword(u.PasswordTemp)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
		u.PasswordTemp = ""
	}
	return nil
}
func (u *User) CheckPassword(password string) bool {
	return CheckPasswordHash(password, u.Password)
}
