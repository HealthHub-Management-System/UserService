package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"

	e "backend/api/resource/common/error"
	validatorUtil "backend/utils/validator"
)

type API struct {
	repository *Repository
	validator  *validator.Validate
}

func New(db *gorm.DB, v *validator.Validate) *API {
	return &API{
		repository: NewRepository(db),
		validator:  v,
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
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	users, err := a.repository.List()
	if err != nil {
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	if len(users) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(users.ToResponse()); err != nil {
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
		e.ServerError(w, e.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			e.ServerError(w, e.RespJSONEncodeFailure)
			return
		}

		e.ValidationErrors(w, respBody)
		return
	}

	newUser := form.ToModel()
	newUser.ID = uuid.New()

	_, err := a.repository.Create(newUser)
	if err != nil {
		e.ServerError(w, e.RespDBDataInsertFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
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
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	user, err := a.repository.Read(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	response := user.ToResponse()
	if err := json.NewEncoder(w).Encode(response); err != nil {
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
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		e.ServerError(w, e.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
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
		e.ServerError(w, e.RespDBDataUpdateFailure)
		return
	}
	if rows == 0 {
		e.NotFound(w)
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
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	rows, err := a.repository.Delete(id)
	if err != nil {
		e.BadRequest(w, e.RespDBDataRemoveFailure)
		return
	}
	if rows == 0 {
		e.NotFound(w)
		return
	}
}
