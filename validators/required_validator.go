package validators

import (
	"errors"
	"reflect"

	"github.com/samuelbeaulieu1/gimlet/actions"
)

func (validator *Validator) ValidateRequired(ctx *ValidationCtx) (bool, error) {
	val := ctx.Value.Interface()
	if reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface()) {
		return false, errors.New("Le champ " + GetFieldLabel(ctx.Field) + " est obligatoire")
	}

	return true, nil
}

func (validator *Validator) ValidateRequiredOnUpdate(ctx *ValidationCtx) (bool, error) {
	if ctx.Action == actions.UpdateAction {
		return validator.ValidateRequired(ctx)
	}

	return true, nil
}

func (validator *Validator) ValidateRequiredOnCreate(ctx *ValidationCtx) (bool, error) {
	if ctx.Action == actions.CreateAction {
		return validator.ValidateRequired(ctx)
	}

	return true, nil
}
