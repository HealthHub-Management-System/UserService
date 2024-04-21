package tests

import (
	"backend/api/resource/users"
	"backend/utils/logger"
	mockDB "backend/utils/mock"
	testUtil "backend/utils/test"
	validatorUtil "backend/utils/validator"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
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
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusOK, status)
	}

	var response users.ListResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	testUtil.NoError(t, err)
	responseUsers := response.Users

	testUtil.Equal(t, len(responseUsers), 2)
	testUtil.Equal(t, responseUsers[0].Name, "user1")
	testUtil.Equal(t, responseUsers[1].Name, "user2")
}
