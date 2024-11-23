package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"strings"
)

type ValidationError struct {
	Message string `json:"message"`
	Field   string `json:"field"`
}

type ValidationErrors struct {
	Message          string            `json:"message"`
	ValidationErrors []ValidationError `json:"errors"`
}

type Validator struct {
	validator  *validator.Validate
	translator ut.Translator
}

type ValidatorInterface interface {
	Validate(data interface{}) ValidationErrors
}

func NewValidator() ValidatorInterface {
	validator := validator.New(validator.WithRequiredStructEnabled())
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	translator := registerTranslator(validator)

	return &Validator{
		validator:  validator,
		translator: translator,
	}
}

func (v *Validator) Validate(data interface{}) ValidationErrors {
	errs := v.validator.Struct(data)
	if errs == nil {
		return ValidationErrors{}
	}

	validationErrors := ValidationErrors{
		Message: "Validation error",
	}
	for _, err := range errs.(validator.ValidationErrors) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, ValidationError{
			Message: err.Translate(v.translator),
			Field:   err.Field(),
		})
	}

	return validationErrors
}

func registerTranslator(validator *validator.Validate) ut.Translator {
	en := en.New()
	uni := ut.New(en, en)
	translator, _ := uni.GetTranslator("en")

	_ = enTranslations.RegisterDefaultTranslations(validator, translator)

	return translator
}
