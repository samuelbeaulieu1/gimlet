package gimlet

import (
	"bytes"
	"encoding/json"
)

type DTO interface {
	GetNewInstance() DTO
}

func ParseModelToDTO(model any, dto any) {
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(model)
	json.NewDecoder(buffer).Decode(&dto)
}

func ParseModelsToDTO[T any](src *[]T, dto DTO) *[]DTO {
	modelsDTO := []DTO{}
	for _, val := range *src {
		modelDTO := dto.GetNewInstance()
		ParseModelToDTO(&val, &modelDTO)
		modelsDTO = append(modelsDTO, modelDTO)
	}

	return &modelsDTO
}
