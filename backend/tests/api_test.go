package tests

import (
	"backend/api/resource/users"
	"backend/utils/logger"
	mockDB "backend/utils/mock"
	testUtil "backend/utils/test"
	validatorUtil "backend/utils/validator"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wader/gormstore/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func mockGormStoreRequests(mock sqlmock.Sqlmock) {
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM (.*)").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("SELECT CURRENT_DATABASE\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	expectedRows := sqlmock.NewRows([]string{"column_name", "is_nullable", "udt_name", "character_maximum_length", "numeric_precision", "numeric_precision_radix", "numeric_scale", "datetime_precision", "typlen", "column_default", "description", "identity_increment"}).
		AddRow("value1", "value2", "value3", 10, 5, 2, 1, 6, 16, "default_value", "description_value", 1)
	mock.ExpectQuery("SELECT (.*) FROM (.*)").WithArgs("1", "sessions").WillReturnRows(expectedRows)
}

func TestGetUsers(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/users", nil)

	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}

	l := logger.New(false)
	v := validatorUtil.New()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	mockGormStoreRequests(mock)
	s := gormstore.New(db, []byte("secret"))

	usersAPI := users.New(l, db, v, s, "APIKey")
	id := uuid.New()
	mockRows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(id, "user1", "email@email.com", "patient").
		AddRow(uuid.New(), "user2", "email2@email.com", "doctor")

	mock.ExpectQuery("^SELECT count\\(\\*\\) FROM \"users\"").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mock.ExpectQuery("^SELECT (.*) FROM \"users\"").WillReturnRows(mockRows)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersAPI.List)

	handler.ServeHTTP(rr, req)
	status := rr.Code
	testUtil.Equal(t, status, http.StatusOK)

	var response users.ListResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	testUtil.NoError(t, err)
	responseUsers := response.Users

	testUtil.Equal(t, len(responseUsers), 2)
	testUtil.Equal(t, responseUsers[0].Name, "user1")
	testUtil.Equal(t, responseUsers[1].Name, "user2")
	testUtil.Equal(t, responseUsers[0].Role, "patient")
	testUtil.Equal(t, responseUsers[1].Role, "doctor")
}

func TestAddUser(t *testing.T) {
	l := logger.New(false)
	v := validatorUtil.New()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	mockGormStoreRequests(mock)
	s := gormstore.New(db, []byte("secret"))

	usersAPI := users.New(l, db, v, s, "APIKey")
	old := users.GetUUID
	defer func() { users.GetUUID = old }()
	users.GetUUID = func() uuid.UUID {
		return uuid.UUID{
			0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
			0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
		}
	}
	oldHash := users.GenerateHash
	defer func() { users.GenerateHash = oldHash }()
	users.GenerateHash = func(password []byte) ([]byte, error) {
		hash := []byte{36, 50, 97, 36, 49, 48, 36, 76, 74, 53, 49, 56, 116, 87, 119, 86, 65, 74, 98, 87, 104, 49,
			106, 86, 72, 97, 51, 89, 117, 101, 80, 98, 83, 104, 118, 54, 74, 97, 56, 48, 98, 68, 78, 101, 71, 104, 50,
			73, 84, 109, 73, 100, 101, 112, 47, 69, 84, 114, 70, 117}
		return hash, nil
	}

	mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").
		WithArgs("email@email.com", 1).
		WillReturnRows(&sqlmock.Rows{})

	password, _ := users.GenerateHash([]byte("password"))
	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO \"users\" ").
		WithArgs(users.GetUUID(), "name", "email@email.com", password, "patient").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := &users.Form{Name: "name", Email: "email@email.com", Password: "Password@123", Role: "patient"}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersAPI.Create)

	body, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewReader(body))
	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}

	handler.ServeHTTP(rr, req)
	status := rr.Code
	testUtil.Equal(t, status, http.StatusCreated)
}

func TestCurrentUser(t *testing.T) {
	idString := "c50abe98-7f20-4cb9-b4a8-fbef37988e7f"
	req, err := http.NewRequest("GET", "/api/v1/users/current", nil)
	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}

	l := logger.New(false)
	v := validatorUtil.New()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	mockGormStoreRequests(mock)
	s := gormstore.New(db, []byte("secret"))
	session, _ := s.Get(req, "session")
	session.Values["id"] = idString

	id, err := uuid.Parse(idString)
	testUtil.NoError(t, err)
	usersAPI := users.New(l, db, v, s, "APIKey")
	mockRows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(id, "user1", "email@email.com", "admin")

	mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").
		WithArgs(id, 1).
		WillReturnRows(mockRows)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersAPI.Current)

	handler.ServeHTTP(rr, req)
	status := rr.Code
	testUtil.Equal(t, status, http.StatusOK)

	var user users.User
	err = json.NewDecoder(rr.Body).Decode(&user)
	testUtil.NoError(t, err)

	testUtil.Equal(t, user.ID, id)
	testUtil.Equal(t, user.Name, "user1")
	testUtil.Equal(t, user.Email, "email@email.com")
	testUtil.Equal(t, user.Role, "admin")

}

func TestGetUser(t *testing.T) {
	idString := "c50abe98-7f20-4cb9-b4a8-fbef37988e7f"
	req, err := http.NewRequest("GET", "/api/v1/users/{id}", nil)
	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", idString)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	l := logger.New(false)
	v := validatorUtil.New()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	mockGormStoreRequests(mock)
	s := gormstore.New(db, []byte("secret"))

	id, err := uuid.Parse(idString)
	testUtil.NoError(t, err)
	usersAPI := users.New(l, db, v, s, "APIKey")
	mockRows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(id, "user1", "email@email.com", "admin")

	mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").
		WithArgs(id, 1).
		WillReturnRows(mockRows)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersAPI.Read)

	handler.ServeHTTP(rr, req)
	status := rr.Code
	testUtil.Equal(t, status, http.StatusOK)

	var user users.User
	err = json.NewDecoder(rr.Body).Decode(&user)
	testUtil.NoError(t, err)

	testUtil.Equal(t, user.ID, id)
	testUtil.Equal(t, user.Name, "user1")
	testUtil.Equal(t, user.Email, "email@email.com")
	testUtil.Equal(t, user.Role, "admin")
}

func TestUpdateUser(t *testing.T) {
	idString := "c50abe98-7f20-4cb9-b4a8-fbef37988e7f"

	l := logger.New(false)
	v := validatorUtil.New()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	mockGormStoreRequests(mock)
	s := gormstore.New(db, []byte("secret"))

	usersAPI := users.New(l, db, v, s, "APIKey")

	id, err := uuid.Parse(idString)
	_ = sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(id, "user1", "email@email.com", "patient")
	testUtil.NoError(t, err)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"users\" SET").
		WithArgs("name", "email2@email.com", id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mockRows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(id, "user1", "email@email.com", "admin")

	mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").
		WithArgs(id, 1).
		WillReturnRows(mockRows)

	user := &users.UpdateForm{Name: "name", Email: "email2@email.com"}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersAPI.Update)

	body, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/v1/users/{id}", bytes.NewReader(body))
	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", idString)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.ServeHTTP(rr, req)
	status := rr.Code
	testUtil.Equal(t, status, http.StatusOK)
}

func TestDeleteUser(t *testing.T) {
	idString := "c50abe98-7f20-4cb9-b4a8-fbef37988e7f"

	l := logger.New(false)
	v := validatorUtil.New()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	mockGormStoreRequests(mock)
	s := gormstore.New(db, []byte("secret"))

	usersAPI := users.New(l, db, v, s, "APIKey")

	id, err := uuid.Parse(idString)
	testUtil.NoError(t, err)
	_ = sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(id, "user1", "email@email.com", "patient")
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM \"users\" WHERE").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersAPI.Delete)

	req, err := http.NewRequest("DELETE", "/api/v1/users/{id}", nil)
	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", idString)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.ServeHTTP(rr, req)
	status := rr.Code
	testUtil.Equal(t, status, http.StatusOK)
}

func TestLogin(t *testing.T) {
	l := logger.New(false)
	v := validatorUtil.New()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	mockGormStoreRequests(mock)
	s := gormstore.New(db, []byte("secret"))

	usersAPI := users.New(l, db, v, s, "APIKey")

	password := "Password@123"
	pass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	email := "email@email.com"
	id := uuid.New()

	mockRows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).
		AddRow(id, "user1", email, pass, "patient")

	mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").
		WithArgs(email, 1).
		WillReturnRows(mockRows)

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO \"sessions\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersAPI.Login)

	user := &users.LoginForm{Email: email, Password: password}
	body, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/v1/users/login", bytes.NewReader(body))
	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}

	handler.ServeHTTP(rr, req)
	status := rr.Code
	returnedUser := &users.UserResponse{}
	err = json.NewDecoder(rr.Body).Decode(returnedUser)
	testUtil.Equal(t, status, http.StatusOK)
	testUtil.NoError(t, err)
	testUtil.Equal(t, returnedUser.Name, "user1")
	testUtil.Equal(t, returnedUser.Email, email)
	testUtil.Equal(t, returnedUser.ID, id)
	testUtil.Equal(t, returnedUser.Role, "patient")
}

func TestLogout(t *testing.T) {
	l := logger.New(false)
	v := validatorUtil.New()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	mockGormStoreRequests(mock)
	s := gormstore.New(db, []byte("secret"))

	usersAPI := users.New(l, db, v, s, "APIKey")

	password := "Password@123"
	pass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	email := "email@email.com"

	mockRows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).
		AddRow(uuid.New(), "user1", email, pass, "patient")

	mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE (.+)").
		WithArgs(email, 1).
		WillReturnRows(mockRows)

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO \"sessions\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersAPI.Login)

	user := &users.LoginForm{Email: email, Password: password}
	body, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/v1/users/logout", bytes.NewReader(body))
	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}

	handler.ServeHTTP(rr, req)
	status := rr.Code
	testUtil.Equal(t, status, http.StatusOK)
}
