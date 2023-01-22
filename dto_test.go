package gimlet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestDTO struct {
	TestValue  string
	TestStruct *TestDTO
}

func (d *TestDTO) GetNewInstance() DTO {
	return &TestDTO{}
}

type TestDTOModel struct {
	TestValue  string
	TestStruct *TestDTOModel
}

func (t TestDTOModel) ToDTO() DTO {
	return nil
}

func (t TestDTOModel) TableName() string {
	return ""
}

func TestParseModelToDTO(t *testing.T) {
	model := &TestDTOModel{
		TestValue: "value1",
		TestStruct: &TestDTOModel{
			TestValue: "value2",
		},
	}

	dto := &TestDTO{}
	ParseModelToDTO(&model, &dto)

	assert.Equal(t, model.TestValue, dto.TestValue)
	assert.NotNil(t, dto.TestStruct)
	assert.Equal(t, model.TestStruct.TestValue, dto.TestStruct.TestValue)
}

func TestParseModelsToDTO(t *testing.T) {
	models := []*TestDTOModel{}
	for i := 0; i < 5; i++ {
		models = append(models, &TestDTOModel{
			TestValue: fmt.Sprintf("value%d", i),
			TestStruct: &TestDTOModel{
				TestValue: "value2",
			},
		})
	}

	dtos := ParseModelsToDTO(&models, &TestDTO{})

	assert.NotNil(t, dtos)
	assert.Equal(t, len(*dtos), len(models))

	for i, dto := range *dtos {
		if dtoInstance, ok := dto.(*TestDTO); ok {
			assert.Equal(t, models[i].TestValue, dtoInstance.TestValue)
			assert.Equal(t, models[i].TestStruct.TestValue, dtoInstance.TestStruct.TestValue)
		} else {
			t.Error("Parse of models to DTO didn't return right DTO struct type")
		}
	}
}
