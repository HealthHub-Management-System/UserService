package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	e "backend/api/resource/common/error"
	validatorUtil "backend/utils/validator"
)

var GetUUID = uuid.New

type API struct {
	repository *Repository
	validator  *validator.Validate
	logger     *zerolog.Logger
}

func New(l *zerolog.Logger, db *gorm.DB, v *validator.Validate) *API {
	return &API{
		repository: NewRepository(db),
		validator:  v,
		logger:     l,
	}
}

// List godoc
//
//	@summary		List users
//	@description	List users
//	@tags			users
//	@accept			json
//	@produce		json
//	@success		200	{array}		ListResponse
//	@failure		500	{object}	error.Error
//	@router			/users [get]
func (a *API) List(w http.ResponseWriter, _ *http.Request) {
	users, err := a.repository.List()
	if err != nil {
		a.logger.Error().Err(err).Msg("List users failed")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	if len(users) == 0 {
		_, _ = fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(users.ToResponse()); err != nil {
		a.logger.Error().Err(err).Msg("List users failed")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
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
//	@success		201
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
//	@success		200
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
