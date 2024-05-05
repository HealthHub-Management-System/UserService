package mock

import (
	e "backend/api/resource/common/error"
	"backend/api/resource/users"
	validatorUtil "backend/utils/validator"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"net/http"
)

type mockAPI struct {
	repository *users.Repository
	validator  *validator.Validate
	logger     *zerolog.Logger
}

func NewMockAPI(l *zerolog.Logger, db *gorm.DB, v *validator.Validate) *mockAPI {
	return &mockAPI{
		repository: users.NewRepository(db),
		validator:  v,
		logger:     l,
	}
}

func (a *mockAPI) Create(w http.ResponseWriter, r *http.Request) {
	form := &users.Form{}
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
	newUser.ID = uuid.New()

	_, err := a.repository.Create(newUser)
	if err != nil {
		a.logger.Error().Err(err).Msg("Create user failed")
		e.ServerError(w, e.RespDBDataInsertFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(newUser.Password)
}
