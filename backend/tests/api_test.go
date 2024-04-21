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
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUsers(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/users", nil)

	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}

	l := logger.New(false)
	v := validatorUtil.New()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	usersAPI := users.New(l, db, v)
	id := uuid.New()
	mockRows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(id, "user1", "email@email.com").
		AddRow(uuid.New(), "user2", "email2@email.com")

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
}

func TestAddUser(t *testing.T) {
	l := logger.New(false)
	v := validatorUtil.New()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	usersAPI := users.New(l, db, v)

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO \"users\" ").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := &users.Form{Name: "name", Email: "email@email.com", Password: "Password@123"}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersAPI.Create)

	body, err := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewReader(body))
	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}

	handler.ServeHTTP(rr, req)
	status := rr.Code
	testUtil.Equal(t, status, http.StatusCreated)
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

	id, err := uuid.Parse(idString)
	testUtil.NoError(t, err)
	usersAPI := users.New(l, db, v)
	mockRows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(id, "user1", "email@email.com")

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

	testUtil.Equal(t, user.Name, "user1")
	testUtil.Equal(t, user.Email, "email@email.com")
}

func TestUpdateUser(t *testing.T) {
	idString := "c50abe98-7f20-4cb9-b4a8-fbef37988e7f"

	l := logger.New(false)
	v := validatorUtil.New()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	usersAPI := users.New(l, db, v)

	id, err := uuid.Parse(idString)
	_ = sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(id, "user1", "email@email.com")
	testUtil.NoError(t, err)
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"users\" SET").
		WithArgs("name", "email2@email.com", id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := &users.UpdateForm{Name: "name", Email: "email2@email.com"}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersAPI.Update)

	body, err := json.Marshal(user)
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

	usersAPI := users.New(l, db, v)

	id, err := uuid.Parse(idString)
	testUtil.NoError(t, err)
	_ = sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(id, "user1", "email@email.com")
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
