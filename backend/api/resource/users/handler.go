package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/wader/gormstore/v2"
	"gorm.io/gorm"

	e "backend/api/resource/common/error"
	"backend/utils/pagination"
	validatorUtil "backend/utils/validator"
)

var GetUUID = uuid.New

type API struct {
	repository *Repository
	validator  *validator.Validate
	logger     *zerolog.Logger
	store      *gormstore.Store
}

func New(l *zerolog.Logger, db *gorm.DB, v *validator.Validate, s *gormstore.Store) *API {
	return &API{
		repository: NewRepository(db),
		validator:  v,
		logger:     l,
		store:      s,
	}
}

// List godoc
//
//	@summary		List users
//	@description	List users
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			page	query	int	false	"Page number"
//	@param			limit	query	int	false	"Number of items per page"
//	@success		200	{object}	ListResponse
//	@failure		500	{object}	error.Error
//	@router			/users [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	pagination := &pagination.Pagination{}
	pagination.Parse(r.URL.Query())
	if err := a.validator.Struct(pagination); err != nil {
		a.logger.Error().Err(err).Msg("List users failed")
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			e.ServerError(w, e.RespJSONEncodeFailure)
			return
		}

		e.ValidationErrors(w, respBody)
		return
	}

	pagination, err := a.repository.List(*pagination)
	if err != nil {
		a.logger.Error().Err(err).Msg("List users failed")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	if users, ok := pagination.Rows.(Users); ok {
		response := users.ToResponse()
		response.TotalItems = pagination.TotalRows
		response.NumberOfPages = pagination.TotalPages
		response.CurrentPage = pagination.Page

		if err := json.NewEncoder(w).Encode(response); err != nil {
			a.logger.Error().Err(err).Msg("List users failed")
			e.ServerError(w, e.RespJSONEncodeFailure)
			return
		}
	} else {
		a.logger.Error().Msg("List users failed")
		e.ServerError(w, e.RespJSONEncodeFailure)
	}
}

// Create godoc
//
//	@summary		Create user
//	@description	Create user
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			body	body	Form	true	"User form"
//	@success		201 {object}	UserResponse
//	@failure		400	{object}	error.Error
//	@failure		422	{object}	error.Errors
//	@failure		500	{object}	error.Error
//	@router			/users [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		a.logger.Error().Err(err).Msg("Create user failed")
		e.ServerError(w, e.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		a.logger.Error().Err(err).Msg("Create user failed")
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			e.ServerError(w, e.RespJSONEncodeFailure)
			return
		}

		e.ValidationErrors(w, respBody)
		return
	}

	if form.Role == Doctor {
		session, err := a.store.Get(r, "session")
		if value, ok := session.Values["role"].(string); !(ok && err == nil && value == Admin.ToString()) {
			a.logger.Error().Err(err).Msg("Not admin tried to add doctor")
			http.Error(w, "Doctor can only be added by admin!", http.StatusUnauthorized)
			return
		}
	}

	newUser := form.ToModel()
	newUser.ID = GetUUID()

	_, err := a.repository.Create(newUser)
	if err != nil {
		a.logger.Error().Err(err).Msg("Create user failed")
		e.ServerError(w, e.RespDBDataInsertFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := newUser.ToResponse()
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logger.Error().Err(err).Msg("Create user failed")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

// Read godoc
//
//	@summary		Read user
//	@description	Read user
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"User ID"
//	@success		200	{object}	UserResponse
//	@failure		400	{object}	error.Error
//	@failure		404
//	@failure		500	{object}	error.Error
//	@router			/users/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	fmt.Println(id)

	if err != nil {
		a.logger.Error().Err(err).Msg("Read user failed")
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	user, err := a.repository.Read(id)
	if err != nil {
		a.logger.Error().Err(err).Msg("Read user failed")
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	response := user.ToResponse()
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logger.Error().Err(err).Msg("Read user failed")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

// Update godoc
//
//	@summary		Update user
//	@description	Update user
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			id		path	string	true	"User ID"
//	@param			body	body	Form	true	"User form"
//	@success		200 {object}	UserResponse
//	@failure		400	{object}	error.Error
//	@failure		404
//	@failure		422	{object}	error.Errors
//	@failure		500	{object}	error.Error
//	@router			/users/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		a.logger.Error().Err(err).Msg("Update user failed")
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	form := &UpdateForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		a.logger.Error().Err(err).Msg("Update user failed")
		e.ServerError(w, e.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		a.logger.Error().Err(err).Msg("Update user failed")
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			e.ServerError(w, e.RespJSONEncodeFailure)
			return
		}

		e.ValidationErrors(w, respBody)
		return
	}

	user := form.ToModel()
	user.ID = id

	rows, err := a.repository.Update(user)
	if err != nil {
		a.logger.Error().Err(err).Msg("Update user failed")
		e.ServerError(w, e.RespDBDataUpdateFailure)
		return
	}
	if rows == 0 {
		e.NotFound(w)
		return
	}

	updatedUser, err := a.repository.Read(id)
	if err != nil {
		a.logger.Error().Err(err).Msg("Update user failed")
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	response := updatedUser.ToResponse()
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logger.Error().Err(err).Msg("Update user failed")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

// Delete godoc
//
//	@summary		Delete user
//	@description	Delete user
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			id	path	string	true	"User ID"
//	@success		200
//	@failure		400	{object}	error.Error
//	@failure		404
//	@failure		500	{object}	error.Error
//	@router			/users/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		a.logger.Error().Err(err).Msg("Delete user failed")
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	rows, err := a.repository.Delete(id)
	if err != nil {
		a.logger.Error().Err(err).Msg("Delete user failed")
		e.BadRequest(w, e.RespDBDataRemoveFailure)
		return
	}
	if rows == 0 {
		e.NotFound(w)
		return
	}
}

// Login godoc
//
//	@summary		Login user
//	@description	Login user
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			body	body	Form	true	"Login form"
//	@success		200
//	@failure		401	{object}	error.Error
//	@failure		422	{object}	error.Errors
//	@failure		500	{object}	error.Error
//	@router			/users/login [post]
func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	session, err := a.store.Get(r, "session")
	if value, ok := session.Values["email"].(string); ok && err == nil {
		if len(value) != 0 {
			a.logger.Error().Err(err).Msg("User already logged in!")
			return
		}
	}

	form := &LoginForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		a.logger.Error().Err(err).Msg("Login user failed")
		e.ServerError(w, e.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		a.logger.Error().Err(err).Msg("Login user failed")
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			e.ServerError(w, e.RespJSONEncodeFailure)
			return
		}

		e.ValidationErrors(w, respBody)
		return
	}

	user, err := a.repository.GetByEmail(form.Email)
	if err != nil || user == nil {
		a.logger.Error().Err(err).Msg("Login user failed")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(form.Password))
	if err != nil {
		a.logger.Error().Err(err).Msg("Login user failed")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	session.Values["email"] = user.Email
	session.Values["role"] = user.Role.ToString()
	err = session.Save(r, w)
	if err != nil {
		a.logger.Error().Err(err).Msg("Login user failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Logout godoc
//
//	@summary		Login user
//	@description	Login user
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			body	body	Form	true	"Login form"
//	@success		200
//	@failure		401	{object}	error.Error
//	@failure		422	{object}	error.Errors
//	@failure		500	{object}	error.Error
//	@router			/users/logout [post]
func (a *API) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := a.store.Get(r, "session")
	if err != nil {
		a.logger.Error().Err(err).Msg("Logout user failed")
	}

	session.Values["email"] = nil
	session.Values["role"] = nil
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		a.logger.Error().Err(err).Msg("Logout user failed")
	}
}
