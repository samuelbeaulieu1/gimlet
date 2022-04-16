package gimlet

import (
	"github.com/samuelbeaulieu1/gimlet/actions"
	"github.com/samuelbeaulieu1/gimlet/responses"
	"github.com/samuelbeaulieu1/gimlet/validators"
)

type ServiceHandler[M Model] interface {
	GetEntity() Entity[M]
	RegisterValidators(*validators.Validator)
}

type OnUpdate[M Model] interface {
	BeforeUpdate(string, *M)
}

type AfterUpdate[M Model] interface {
	AfterUpdate(*M)
}

type OnCreate[M Model] interface {
	BeforeCreate(*M)
}

type AfterCreate[M Model] interface {
	AfterCreate(*M)
}

type OnDelete interface {
	BeforeDelete(string)
}

type AfterDelete interface {
	AfterDelete(string)
}

type ServiceInterface[M Model] interface {
	Exists(string) responses.Error
	GetAll() (*[]M, responses.Error)
	Get(string) (*M, responses.Error)
	Update(string, *M) responses.Error
	Create(*M) (*M, responses.Error)
	Delete(string) responses.Error
	NewValidator() *validators.Validator
	ServiceHandler[M]
}

type Service[M Model] struct {
	ServiceHandler[M]
}

func (service *Service[M]) Exists(id string) responses.Error {
	if !service.GetEntity().Exists(id) {
		return responses.NewError("Record inexistant")
	}

	return nil
}

func (service *Service[M]) GetAll() (*[]M, responses.Error) {
	records, err := service.GetEntity().GetAll()
	if err != nil {
		return nil, responses.NewError("Impossible de récupérer la liste")
	}

	return records, nil
}

func (service *Service[M]) Get(id string) (*M, responses.Error) {
	if err := service.Exists(id); err != nil {
		return nil, err
	}

	record, err := service.GetEntity().Get(id)
	if err != nil {
		return nil, responses.NewError("Impossible de récupérer le record")
	}

	return record, nil
}

func (service *Service[M]) Update(id string, request *M) responses.Error {
	if err := service.NewValidator().ValidateModel(actions.UpdateAction, request); err != nil {
		return err
	}
	if err := service.Exists(id); err != nil {
		return err
	}

	if onUpdateHandler, ok := service.ServiceHandler.(OnUpdate[M]); ok {
		onUpdateHandler.BeforeUpdate(id, request)
	}
	err := service.GetEntity().Update(id, request)
	if err != nil {
		return responses.NewError("Impossible de modifier le record")
	}
	if afterUpdateHandler, ok := service.ServiceHandler.(AfterUpdate[M]); ok {
		afterUpdateHandler.AfterUpdate(request)
	}

	return nil
}

func (service *Service[M]) Create(request *M) (*M, responses.Error) {
	if err := service.NewValidator().ValidateModel(actions.CreateAction, request); err != nil {
		return nil, err
	}

	if onCreateHandler, ok := service.ServiceHandler.(OnCreate[M]); ok {
		onCreateHandler.BeforeCreate(request)
	}
	record, err := service.GetEntity().Create(request)
	if err != nil {
		return nil, responses.NewError("Impossible de créer le record")
	}
	if afterCreateHandler, ok := service.ServiceHandler.(AfterCreate[M]); ok {
		afterCreateHandler.AfterCreate(record)
	}

	return record, nil
}

func (service *Service[M]) Delete(id string) responses.Error {
	if err := service.Exists(id); err != nil {
		return err
	}

	if onDeleteHandler, ok := service.ServiceHandler.(OnDelete); ok {
		onDeleteHandler.BeforeDelete(id)
	}
	if err := service.GetEntity().Delete(id); err != nil {
		return responses.NewError("Impossible de supprimer le record")
	}
	if afterDeleteHandler, ok := service.ServiceHandler.(AfterDelete); ok {
		afterDeleteHandler.AfterDelete(id)
	}

	return nil
}

func (service *Service[M]) NewValidator() *validators.Validator {
	validator := validators.NewValidator()
	service.RegisterValidators(validator)

	return validator
}
