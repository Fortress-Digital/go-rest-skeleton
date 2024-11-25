package validation

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

type ValidationStruct struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Nothing string `json:"-"`
}

func TestValidator_Validate(t *testing.T) {
	d := ValidationStruct{
		Name:  "Test",
		Email: "test@example.com",
	}

	sut := NewValidator()
	errs := sut.Validate(d)

	assert.Equal(t, errs, ValidationErrors{})
}

func TestValidator_Validate_Invalid(t *testing.T) {
	d := ValidationStruct{
		Name:  "",
		Email: "test",
	}

	sut := NewValidator()
	errs := sut.Validate(d)

	assert.Equal(t, errs.Message, "Validation error")
	assert.Equal(t, len(errs.ValidationErrors), 2)
	assert.Equal(t, errs.ValidationErrors[0].Field, "name")
	assert.Equal(t, errs.ValidationErrors[0].Message, "name is a required field")
	assert.Equal(t, errs.ValidationErrors[1].Field, "email")
	assert.Equal(t, errs.ValidationErrors[1].Message, "email must be a valid email address")
}
