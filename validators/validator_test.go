package validators

import (
	"errors"
	"reflect"
	"testing"

	"github.com/samuelbeaulieu1/gimlet/actions"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	TestValue        string `label:"testLabel"`
	TestValueNoLabel string `validate:"required"`
}

type TestModelRequired struct {
	TestValue string `validate:"required" label:"testLabel"`
}

type TestModelRequiredCreate struct {
	TestValue string `validate:"requiredOnCreate" label:"testLabel"`
}

type TestModelRequiredUpdate struct {
	TestValue string `validate:"requiredOnUpdate" label:"testLabel"`
}

type TestModelCustomValidator struct {
	TestValue string `validate:"custom"`
}

func ValidateFuncTest(ctx *ValidationCtx) (bool, error) {
	return true, nil
}

func ValidateFuncTestInvalid(ctx *ValidationCtx) (bool, error) {
	return false, errors.New("test")
}

func TestNewValidator(t *testing.T) {
	validator := NewValidator()

	assert.NotNil(t, validator)
}

func TestRegisterValidation(t *testing.T) {
	validator := NewValidator()
	validator.validators = make(map[string]Validation)

	assert.NotNil(t, validator)
	assert.Equal(t, 0, len(validator.validators))

	validator.RegisterValidation("test", ValidateFuncTest)
	assert.Equal(t, 1, len(validator.validators))

	_, ok := validator.validators["test"]
	assert.True(t, ok)
}

func TestGetModelValue(t *testing.T) {
	validator := NewValidator()

	assert.NotNil(t, validator)

	model := TestModel{}
	reflectValue := validator.getModelValue(model)
	assert.NotNil(t, reflectValue)
}

func TestGetModelPtrValue(t *testing.T) {
	validator := NewValidator()

	assert.NotNil(t, validator)

	ptrModel := &TestModel{}
	reflectPtrValue := validator.getModelValue(ptrModel)
	assert.NotNil(t, reflectPtrValue)
}

func TestGetFieldLabel(t *testing.T) {
	validator := NewValidator()

	assert.NotNil(t, validator)

	model := &TestModel{
		TestValue: "val",
	}
	value := reflect.ValueOf(model).Elem()
	field, _ := value.Type().FieldByName("TestValue")
	label := GetFieldLabel(field)
	assert.Equal(t, "testLabel", label)
}

func TestGetFieldLabelEmpty(t *testing.T) {
	validator := NewValidator()

	assert.NotNil(t, validator)

	model := &TestModel{
		TestValue: "val",
	}
	value := reflect.ValueOf(model).Elem()

	fieldNoLabel, _ := value.Type().FieldByName("TestValueNoLabel")
	noLabel := GetFieldLabel(fieldNoLabel)
	assert.Equal(t, "TestValueNoLabel", noLabel)
}

func TestInitDefaultValidators(t *testing.T) {
	validator := NewValidator()

	assert.NotNil(t, validator)
	assert.Equal(t, len(defaultValidators), len(validator.validators))

	for name, _ := range defaultValidators {
		_, ok := validator.validators[name]
		assert.True(t, ok)
	}
}

func TestHandleValidator(t *testing.T) {
	validator := NewValidator()
	validator.validators = make(map[string]Validation)
	validator.validators["test"] = ValidateFuncTest

	valid, err := validator.handleValidator("test", &ValidationCtx{})
	assert.Nil(t, err)
	assert.True(t, valid)
}

func TestHandleUnknownValidator(t *testing.T) {
	validator := NewValidator()
	validator.validators = make(map[string]Validation)

	valid, err := validator.handleValidator("test", &ValidationCtx{})
	assert.Nil(t, err)
	assert.True(t, valid)
}

func TestValidateModel(t *testing.T) {
	validator := NewValidator()
	validator.validators = make(map[string]Validation)
	validator.validators["required"] = ValidateFuncTest
	model := &TestModel{}

	errors := validator.ValidateModel(actions.CreateAction, model)
	assert.Nil(t, errors)
}

func TestValidateInvalidModel(t *testing.T) {
	validator := NewValidator()
	validator.validators = make(map[string]Validation)
	validator.validators["required"] = ValidateFuncTestInvalid
	model := &TestModel{}

	errors := validator.ValidateModel(actions.CreateAction, model)
	assert.NotNil(t, errors)
	assert.Equal(t, "test", errors.Error())
}

func TestValidateModelCustomValidate(t *testing.T) {
	validator := NewValidator()
	validator.RegisterValidation("custom", ValidateFuncTest)
	model := &TestModelCustomValidator{}

	errors := validator.ValidateModel(actions.CreateAction, model)
	assert.Nil(t, errors)
}

func TestValidateRequiredModel(t *testing.T) {
	validator := NewValidator()
	model := &TestModelRequired{
		TestValue: "test",
	}

	errors := validator.ValidateModel(actions.CreateAction, model)
	assert.Nil(t, errors)

	invalidModel := &TestModelRequired{}

	errors = validator.ValidateModel(actions.CreateAction, invalidModel)
	assert.NotNil(t, errors)
}

func TestValidateRequiredOnCreateModel(t *testing.T) {
	validator := NewValidator()
	model := &TestModelRequiredCreate{
		TestValue: "test",
	}

	errors := validator.ValidateModel(actions.CreateAction, model)
	assert.Nil(t, errors)

	invalidModel := &TestModelRequiredCreate{}

	errors = validator.ValidateModel(actions.CreateAction, invalidModel)
	assert.NotNil(t, errors)

	errors = validator.ValidateModel(actions.UpdateAction, invalidModel)
	assert.Nil(t, errors)
}

func TestValidateRequiredOnUpdateModel(t *testing.T) {
	validator := NewValidator()
	model := &TestModelRequiredUpdate{
		TestValue: "test",
	}

	errors := validator.ValidateModel(actions.UpdateAction, model)
	assert.Nil(t, errors)

	invalidModel := &TestModelRequiredUpdate{}

	errors = validator.ValidateModel(actions.UpdateAction, invalidModel)
	assert.NotNil(t, errors)

	errors = validator.ValidateModel(actions.CreateAction, invalidModel)
	assert.Nil(t, errors)
}
