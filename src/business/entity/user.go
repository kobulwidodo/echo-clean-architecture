package entity

import (
	"go-clean/src/lib/auth"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string `json:"-"`
	Name     string
	IsAdmin  bool
}

type CreateUserParam struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
	Name     string `validate:"required"`
}

type LoginUserParam struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

func (u *User) ConvertToAuthUser() auth.User {
	return auth.User{
		ID:       u.ID,
		Username: u.Username,
		IsAdmin:  u.IsAdmin,
		Name:     u.Name,
	}
}
