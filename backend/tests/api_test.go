package tests

import (
	"backend/api/resource/users"
	"backend/utils/logger"
	mockDB "backend/utils/mock"
	testUtil "backend/utils/test"
	validatorUtil "backend/utils/validator"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPostsHandler(t *testing.T) {
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
		AddRow(id, "user1", "email@email.com")

	mock.ExpectQuery("^SELECT (.*) FROM \"users\"").WillReturnRows(mockRows)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersAPI.List)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusOK, status)
	}
}
