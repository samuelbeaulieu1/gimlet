package validators

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/samuelbeaulieu1/gimlet/actions"
	"github.com/samuelbeaulieu1/gimlet/logger"
	"github.com/samuelbeaulieu1/gimlet/responses"
)

type Validation func(action actions.Action, value reflect.Value, field reflect.StructField) (bool, error)

type Validator struct {
	validators map[string]Validation
}

func NewValidator() *Validator {
	return &Validator{
		validators: make(map[string]Validation),
	}
}

func (validator *Validator) RegisterValidation(name string, validation Validation) {
	validator.validators[name] = validation
}

func (validator *Validator) getModelValue(model any) reflect.Value {
	var val reflect.Value
	if reflect.ValueOf(model).Kind() == reflect.Ptr {
		val = reflect.ValueOf(model).Elem()
	} else {
		val = reflect.ValueOf(model)
	}

	return val
}

func (validator *Validator) ValidateModel(action actions.Action, model any) responses.Error {
	val := validator.getModelValue(model)

	errFields := []string{}
	err := []string{}
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		validateTag := typeField.Tag.Get("validate")
		if len(validateTag) == 0 {
			continue
		}

		tags := strings.Split(typeField.Tag.Get("validate"), ",")
		valid := true
		for _, tag := range tags {
			isValid, validationErr := validator.handleValidator(action, tag, valueField, typeField)
			valid = valid && isValid
			if validationErr != nil {
				err = append(err, validationErr.Error())
			}
		}

		if !valid {
			jsonName := typeField.Tag.Get("json")
			if jsonName != "" {
				errFields = append(errFields, jsonName)
			} else {
				errFields = append(errFields, typeField.Name)
			}
		}
	}

	if len(errFields) > 0 {
		return responses.NewFieldsError(err, errFields)
	}
	return nil
}

func (validator *Validator) handleValidator(action actions.Action, validatorTag string, value reflect.Value, field reflect.StructField) (bool, error) {
	valid := false
	var err error

	switch validatorTag {
	case "required":
		valid, err = validator.ValidateRequired(action, value, field)
	case "requiredOnUpdate":
		valid, err = validator.ValidateRequiredOnUpdate(action, value, field)
	case "requiredOnCreate":
		valid, err = validator.ValidateRequiredOnCreate(action, value, field)
	default:
		if validation, ok := validator.validators[validatorTag]; ok {
			return validation(action, value, field)
		}
		logger.PrintError(fmt.Sprintf("Invalid validator %s for field %s", validatorTag, field.Name))
	}

	return valid, err
}

func GetFieldLabel(field reflect.StructField) string {
	label := field.Tag.Get("label")

	if label == "" {
		return field.Name
	}

	return label
}
