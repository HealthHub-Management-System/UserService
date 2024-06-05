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
	apiKey     string
}

func New(l *zerolog.Logger, db *gorm.DB, v *validator.Validate, s *gormstore.Store, apiKey string) *API {
	return &API{
		repository: NewRepository(db),
		validator:  v,
		logger:     l,
		store:      s,
		apiKey:     apiKey,
	}
}

// List godoc
//
//	@summary		List users
//	@description	List users
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			page	query	int	false		"Page number"
//	@param			limit	query	int	false		"Number of items per page"
//	@param			role	query	string false	"Role to filter by"
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

	pagination = a.repository.List(*pagination)

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

	if r.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", a.apiKey) {
		if form.Role == Admin {
			a.logger.Error().Msg("Not allowed to create admin")
			http.Error(w, "Cannot create admin from the level of API!", http.StatusUnauthorized)
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
	}

	user, err1 := a.repository.GetByEmail(form.Email)
	if err1 != nil && !errors.Is(err1, gorm.ErrRecordNotFound) {
		a.logger.Error().Err(err1).Msg("Create user failed")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	} else if user != nil {
		a.logger.Error().Msg("User already exists")
		http.Error(w, "User already exists!", http.StatusConflict)
		return
	}

	newUser := form.ToModel()
	newUser.ID = GetUUID()

	_, err2 := a.repository.Create(newUser)
	if err2 != nil {
		a.logger.Error().Err(err2).Msg("Create user failed")
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

// Current godoc
//
//	@summary		Current user
//	@description	Current user
//	@tags			users
//	@accept			json
//	@produce		json
//	@success		200	{object}	UserResponse
//	@failure		400	{object}	error.Error
//	@failure		404
//	@failure		500	{object}	error.Error
//	@router			/users/current [get]
func (a *API) Current(w http.ResponseWriter, r *http.Request) {
	session, err := a.store.Get(r, "session")
	if err != nil {
		a.logger.Error().Err(err).Msg("Getting current user failed")
		e.ServerError(w, e.RespSessionAccessFailure)
		return
	}

	idString, ok := session.Values["id"].(string)
	if !ok {
		a.logger.Error().Msg("Getting current user failed")
		e.ServerError(w, e.RespSessionAccessFailure)
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		a.logger.Error().Msg(fmt.Sprintf("Error parsing UUID: %v\n", err))
		e.ServerError(w, e.RespSessionAccessFailure)
		return
	}

	user, err := a.repository.Read(id)
	if err != nil {
		a.logger.Error().Err(err).Msg("Getting current user failed")
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	response := user.ToResponse()
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logger.Error().Err(err).Msg("Getting current user failed")
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

	session.Values["id"] = user.ID.String()
	session.Values["email"] = user.Email
	session.Values["role"] = user.Role.ToString()

	if err := session.Save(r, w); err != nil {
		a.logger.Error().Err(err).Msg("Login user failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := user.ToResponse()
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logger.Error().Err(err).Msg("Update user failed")
		e.ServerError(w, e.RespJSONEncodeFailure)
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

	session.Values["id"] = nil
	session.Values["email"] = nil
	session.Values["role"] = nil
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		a.logger.Error().Err(err).Msg("Logout user failed")
	}
}
