package tests_test

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"backend/api/resource/users"
	mockDB "backend/utils/mock"
	"backend/utils/pagination"
	testUtil "backend/utils/test"
)

func TestRepository_List(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := users.NewRepository(db)

	mockRows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(uuid.New(), "user1", "email1@email.com", "patient").
		AddRow(uuid.New(), "user2", "email2@email.com", "patient")

	mock.ExpectQuery("^SELECT count\\(\\*\\) FROM \"users\"").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mock.ExpectQuery("^SELECT (.+) FROM \"users\"").WillReturnRows(mockRows)

	pagination := &pagination.Pagination{
		Page:  1,
		Limit: 10,
		Role:  nil,
	}

	result := repo.List(*pagination)
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(result.Rows.(users.Users)), 2)
	testUtil.Equal(t, result.Page, 1)
	testUtil.Equal(t, result.Limit, 10)
	testUtil.Equal(t, result.TotalPages, 1)
	testUtil.Equal(t, result.TotalRows, 2)
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := users.NewRepository(db)

	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	id := uuid.New()
	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO \"users\" ").
		WithArgs(id, "name", "email", password, "patient").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := &users.User{ID: id, Name: "name", Email: "email", Password: password, Role: users.Patient}
	_, err = repo.Create(user)
	testUtil.NoError(t, err)
}

func TestRepository_Read(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := users.NewRepository(db)

	id := uuid.New()
	mockRows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(id, "user1", "email@email.com", "patient")

	mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").
		WithArgs(id, 1).
		WillReturnRows(mockRows)

	user, err := repo.Read(id)
	testUtil.NoError(t, err)
	testUtil.Equal(t, "user1", user.Name)
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := users.NewRepository(db)

	id := uuid.New()
	_ = sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(id, "user1", "email@email.com", "patient")

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"users\" SET").
		WithArgs("name", "email", id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	user := &users.User{ID: id, Name: "name", Email: "email", Password: password, Role: users.Patient}
	rows, err := repo.Update(user)
	testUtil.NoError(t, err)
	testUtil.Equal(t, 1, rows)
}

func TestRepository_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := users.NewRepository(db)

	id := uuid.New()
	_ = sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(id, "user1", "email@email.com", "patient")

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM \"users\" WHERE").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rows, err := repo.Delete(id)
	testUtil.NoError(t, err)
	testUtil.Equal(t, 1, rows)
}

func TestRepository_GetByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := users.NewRepository(db)

	email := "email@email.com"
	mockRows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(uuid.New(), "user1", email, "patient")

	mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").
		WithArgs(email, 1).
		WillReturnRows(mockRows)

	user, err := repo.GetByEmail(email)
	testUtil.NoError(t, err)
	testUtil.Equal(t, "user1", user.Name)
}
