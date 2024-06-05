package users

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend/utils/pagination"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List(p pagination.Pagination) (*pagination.Pagination, error) {
	var users Users

	r.db.Scopes(pagination.Paginate(users, &p, r.db)).Find(&users)

	if len(users) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	p.Rows = users

	return &p, nil
}

func (r *Repository) Create(user *User) (*User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) Read(id uuid.UUID) (*User, error) {
	user := &User{}
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetByEmail(email string) (*User, error) {
	user := &User{}
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetByRole(role string) (Users, error) {
	users := make([]*User, 0)
	if err := r.db.Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) Update(user *User) (int64, error) {
	result := r.db.Model(&User{}).
		Select("name", "email").
		Where("id = ?", user.ID).
		Updates(user)

	return result.RowsAffected, result.Error
}

func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&User{})

	return result.RowsAffected, result.Error
}
