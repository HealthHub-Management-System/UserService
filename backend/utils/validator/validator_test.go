package validator_test

import (
	"testing"

	"backend/utils/validator"
)

type testCase struct {
	name     string
	input    interface{}
	expected string
}

var tests = []*testCase{
	{
		name: `required`,
		input: struct {
			Name string `json:"name" form:"required"`
		}{},
		expected: "name is a required field",
	},
	{
		name: `max`,
		input: struct {
			Name string `json:"name" form:"max=7"`
		}{Name: "12345678"},
		expected: "name must be a maximum of 7 in length",
	},
	{
		name: `alpha_space`,
		input: struct {
			Name string `json:"name" form:"alpha_space"`
		}{Name: "Some Name 2"},
		expected: "name can only contain alphabetic and space characters",
	},
	{
		name: `password`,
		input: struct {
			Password string `json:"password" form:"password"`
		}{Password: "password"},
		expected: "password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character",
	},
}

func TestToErrResponse(t *testing.T) {
	vr := validator.New()

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := vr.Struct(tc.input)
			if errResp := validator.ToErrResponse(err); errResp == nil || len(errResp.Errors) != 1 {
				t.Fatalf(`Expected:"{[%v]}", Got:"%v"`, tc.expected, errResp)
			} else if errResp.Errors[0] != tc.expected {
				t.Fatalf(`Expected:"%v", Got:"%v"`, tc.expected, errResp.Errors[0])
			}
		})
	}
}
