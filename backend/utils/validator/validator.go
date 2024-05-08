package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrResponse struct {
	Errors []string `json:"errors"`
}

func New() *validator.Validate {
	validate := validator.New()
	validate.SetTagName("form")

	// Using the names which have been specified for JSON representations of structs, rather than normal Go field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	_ = validate.RegisterValidation("alpha_space", isAlphaSpace)
	_ = validate.RegisterValidation("password", isPassword)
	_ = validate.RegisterValidation("role", isRole)
	_ = validate.RegisterValidation("page", greaterOrEqual0)
	_ = validate.RegisterValidation("limit", greaterOrEqual0)

	return validate
}

func ToErrResponse(err error) *ErrResponse {
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		resp := ErrResponse{
			Errors: make([]string, len(fieldErrors)),
		}

		for i, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				resp.Errors[i] = fmt.Sprintf("%s is a required field", err.Field())
			case "max":
				resp.Errors[i] = fmt.Sprintf("%s must be a maximum of %s in length", err.Field(), err.Param())
			case "alpha_space":
				resp.Errors[i] = fmt.Sprintf("%s can only contain alphabetic and space characters", err.Field())
			case "email":
				resp.Errors[i] = fmt.Sprintf("%s must be a valid email address", err.Field())
			case "password":
				resp.Errors[i] = fmt.Sprintf("%s must contain at least one uppercase letter, one lowercase letter, one digit, and one special character", err.Field())
			case "role":
				resp.Errors[i] = fmt.Sprintf("%s must be one of the following: patient, doctor, admin", err.Field())
			case "required_without":
				resp.Errors[i] = fmt.Sprintf("%s is required when another field is absent", err.Field())
			case "page":
				resp.Errors[i] = fmt.Sprintf("%s must be greater than 0", err.Field())
			case "limit":
				resp.Errors[i] = fmt.Sprintf("%s must be greater than 0", err.Field())
			default:
				resp.Errors[i] = fmt.Sprintf("something wrong on %s; %s", err.Field(), err.Tag())
			}
		}

		return &resp
	}

	return nil
}

func isAlphaSpace(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile("^[a-zA-Z ]+$")
	return reg.MatchString(fl.Field().String())
}

func isPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check length
	if len(password) < 8 {
		return false
	}

	// Check for at least one uppercase letter, one lowercase letter, one digit, and one special character
	hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowerCase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	return hasUpperCase && hasLowerCase && hasDigit && hasSpecial
}

func isRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	roles := []string{"patient", "doctor", "admin"}

	return slices.Contains(roles, role)
}

func greaterOrEqual0(fl validator.FieldLevel) bool {
	return fl.Field().Int() >= 0
}
