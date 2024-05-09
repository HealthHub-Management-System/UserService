package users

import (
	"slices"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "gorm.io/gorm" // nolint
)

var GenerateHash = func(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

type Role string

const (
	Unknown Role = "unknown"
	Patient Role = "patient"
	Doctor  Role = "doctor"
	Admin   Role = "admin"
)

func (r Role) ToString() string {
	return string(r)
}

func ToRole(s string) Role {
	roles := []string{"unknown", "patient", "doctor", "admin"}
	if slices.Contains(roles, s) {
		return Role(s)
	}
	return Unknown
}

type ListResponse struct {
	Users         []*UserResponse `json:"users"`
	TotalItems    int64           `json:"total"`
	NumberOfPages int             `json:"pages"`
	CurrentPage   int             `json:"currentPage"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

type Form struct {
	Name     string `json:"name" form:"required,alpha_space,max=255"`
	Email    string `json:"email" form:"required,email,max=255"`
	Password string `json:"password" form:"required,password,max=255"`
	Role     Role   `json:"role" form:"required,role"`
}

type UpdateForm struct {
	Name  string `json:"name" form:"required_without=Email,alpha_space,max=255"`
	Email string `json:"email" form:"required_without=Name,email,max=255"`
}

type LoginForm struct {
	Email    string `json:"email" form:"email,max=255"`
	Password string `json:"password" form:"required,password,max=255"`
}

type User struct {
	ID       uuid.UUID `gorm:"primarykey"`
	Name     string
	Email    string
	Password []byte
	Role     Role `gorm:"type:Role, default:unknown"`
}

type Users []*User

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Role:  u.Role.ToString(),
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
		Role:     f.Role,
	}
}

func (f *UpdateForm) ToModel() *User {
	return &User{
		ID:       uuid.New(),
		Name:     f.Name,
		Email:    f.Email,
		Password: []byte("password"), // Not allowed to update password
		Role:     Unknown,            // Not allowed to update role
	}
}
