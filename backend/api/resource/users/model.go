package users

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "gorm.io/gorm"
)

var GenerateHash = func(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

type ListResponse struct {
	Users []*UserResponse `json:"users"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type Form struct {
	Name     string `json:"name" form:"required,alpha_space,max=255"`
	Email    string `json:"email" form:"required,email,max=255"`
	Password string `json:"password" form:"required,password,max=255"`
}

type UpdateForm struct {
	Name  string `json:"name" form:"required_without=Email,alpha_space,max=255"`
	Email string `json:"email" form:"required_without=Name,email,max=255"`
}

type User struct {
	ID       uuid.UUID `gorm:"primarykey"`
	Name     string
	Email    string
	Password []byte
}

type Users []*User

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func (users Users) ToResponse() *ListResponse {
	var response []*UserResponse
	for _, u := range users {
		response = append(response, u.ToResponse())
	}
	return &ListResponse{Users: response}
}

func (f *Form) ToModel() *User {
	password, _ := GenerateHash([]byte(f.Password))

	return &User{
		ID:       uuid.New(),
		Name:     f.Name,
		Email:    f.Email,
		Password: password,
	}
}

func (f *UpdateForm) ToModel() *User {
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	return &User{
		ID:       uuid.New(),
		Name:     f.Name,
		Email:    f.Email,
		Password: password,
	}
}
