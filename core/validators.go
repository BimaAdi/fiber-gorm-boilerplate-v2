package core

import (
	"github.com/BimaAdi/fiberGormBoilerplate/schemas"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

// ValidateSchemas
// Return value (is_valid, validation_error)
func ValidateSchemas(s interface{}) (bool, schemas.UnprocessableEntityResponse) {
	var errors []map[string]string
	err := Validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, map[string]string{err.Field(): err.Error()})
		}
		return false, schemas.UnprocessableEntityResponse{
			Message: errors,
		}
	}
	return true, schemas.UnprocessableEntityResponse{}
}
