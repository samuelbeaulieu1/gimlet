package validators

import (
	"errors"
	"reflect"

	"github.com/samuelbeaulieu1/gimlet/actions"
)

func ValidateRequired(ctx *ValidationCtx) (bool, error) {
	val := ctx.Value.Interface()
	if reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface()) {
		return false, errors.New("Le champ " + GetFieldLabel(ctx.Field) + " est obligatoire")
	}

	return true, nil
}

func ValidateRequiredOnUpdate(ctx *ValidationCtx) (bool, error) {
	if ctx.Action == actions.UpdateAction {
		return ValidateRequired(ctx)
	}

	return true, nil
}

func ValidateRequiredOnCreate(ctx *ValidationCtx) (bool, error) {
	if ctx.Action == actions.CreateAction {
		return ValidateRequired(ctx)
	}

	return true, nil
}
