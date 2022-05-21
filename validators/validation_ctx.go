package validators

import (
	"reflect"

	"github.com/samuelbeaulieu1/gimlet/actions"
)

type ValidationCtx struct {
	Model  any
	Action actions.Action
	Value  reflect.Value
	Field  reflect.StructField
}
