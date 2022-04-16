package validators

import (
	"errors"
	"reflect"

	"github.com/samuelbeaulieu1/gimlet/actions"
)

func (validator *Validator) ValidateRequired(action actions.Action, value reflect.Value, field reflect.StructField) (bool, error) {
	val := value.Interface()
	if reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface()) {
		return false, errors.New("Le champ " + GetFieldLabel(field) + " est obligatoire")
	}

	return true, nil
}

func (validator *Validator) ValidateRequiredOnUpdate(action actions.Action, value reflect.Value, field reflect.StructField) (bool, error) {
	if action == actions.UpdateAction {
		return validator.ValidateRequired(action, value, field)
	}

	return true, nil
}

func (validator *Validator) ValidateRequiredOnCreate(action actions.Action, value reflect.Value, field reflect.StructField) (bool, error) {
	if action == actions.CreateAction {
		return validator.ValidateRequired(action, value, field)
	}

	return true, nil
}
