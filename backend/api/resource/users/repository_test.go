package users_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"backend/api/resource/users"
	mockDB "backend/utils/mock"
	testUtil "backend/utils/test"
)

func TestRepository_List(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := users.NewRepository(db)

	mockRows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(uuid.New(), "user1", "email1@email.com").
		AddRow(uuid.New(), "user2", "email2@email.com")

	mock.ExpectQuery("^SELECT (.+) FROM \"users\"").WillReturnRows(mockRows)

	users, err := repo.List()
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(users), 2)
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := users.NewRepository(db)

	id := uuid.New()
	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO \"users\" ").
		WithArgs(id, "name", "email", "password").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := &users.User{ID: id, Name: "name", Email: "email", Password: "password"}
	_, err = repo.Create(user)
	testUtil.NoError(t, err)
}

func TestRepository_Read(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := users.NewRepository(db)

	id := uuid.New()
	mockRows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(id, "user1", "email@email.com")

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
	_ = sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(id, "user1", "email@email.com")

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"users\" SET").
		WithArgs("name", "email", id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := &users.User{ID: id, Name: "name", Email: "email", Password: "password"}
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
	_ = sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(id, "user1", "email@email.com")

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM \"users\" WHERE").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rows, err := repo.Delete(id)
	testUtil.NoError(t, err)
	testUtil.Equal(t, 1, rows)
}
