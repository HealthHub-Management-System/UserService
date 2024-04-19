package users

import (
	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type Form struct {
	Name     string `json:"name" validate:"required,alpha_space,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=255"`
}

type User struct {
	ID       uuid.UUID `gorm:"primarykey"`
	Name     string
	Email    string
	Password string
}

type Users []*User

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func (users Users) ToResponse() []*UserResponse {
	var response []*UserResponse
	for _, u := range users {
		response = append(response, u.ToResponse())
	}
	return response
}

func (f *Form) ToModel() *User {
	return &User{
		ID:       uuid.New(),
		Name:     f.Name,
		Email:    f.Email,
		Password: f.Password,
	}
}
