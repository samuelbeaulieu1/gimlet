package gimlet

import (
	"github.com/samuelbeaulieu1/gimlet/actions"
	"github.com/samuelbeaulieu1/gimlet/responses"
	"github.com/samuelbeaulieu1/gimlet/validators"
)

type ServiceHandler[M Model] interface {
	GetEntity() Entity[M]
	RegisterValidators(actions.Action, *M, *validators.Validator)
}

type OnUpdate[M Model] interface {
	BeforeUpdate(string, *M) responses.Error
}

type AfterUpdate[M Model] interface {
	AfterUpdate(*M) responses.Error
}

type OnCreate[M Model] interface {
	BeforeCreate(*M) responses.Error
}

type AfterCreate[M Model] interface {
	AfterCreate(*M) responses.Error
}

type OnDelete interface {
	BeforeDelete(string) responses.Error
}

type AfterDelete interface {
	AfterDelete(string) responses.Error
}

type ServiceInterface[M Model] interface {
	Exists(string) responses.Error
	GetAll() (*[]M, responses.Error)
	Get(string) (*M, responses.Error)
	Update(string, *M) responses.Error
	Create(*M) (*M, responses.Error)
	Delete(string) responses.Error
	NewValidator(actions.Action, *M) *validators.Validator
	ServiceHandler[M]
}

type Service[M Model] struct {
	ServiceHandler[M]
	ResponseHandler responses.ResponseHandler
}

func NewService[M Model]() *Service[M] {
	return &Service[M]{
		ResponseHandler: &responses.ResponseMapper{},
	}
}

func (service *Service[M]) Exists(id string) responses.Error {
	if !service.GetEntity().Exists(id) {
		return service.ResponseHandler.Error(responses.ERR_RECORD_NOT_FOUND)
	}

	return nil
}

func (service *Service[M]) GetAll() (*[]M, responses.Error) {
	records, err := service.GetEntity().GetAll()
	if err != nil {
		return nil, service.ResponseHandler.Error(responses.ERR_GET_RECORDS)
	}

	return records, nil
}

func (service *Service[M]) Get(id string) (*M, responses.Error) {
	if err := service.Exists(id); err != nil {
		return nil, err
	}

	record, err := service.GetEntity().Get(id)
	if err != nil {
		return nil, service.ResponseHandler.Error(responses.ERR_GET_RECORD)
	}

	return record, nil
}

func (service *Service[M]) Update(id string, request *M) responses.Error {
	if err := service.NewValidator(actions.UpdateAction, request).ValidateModel(actions.UpdateAction, request); err != nil {
		return err
	}
	if err := service.Exists(id); err != nil {
		return err
	}

	if onUpdateHandler, ok := service.ServiceHandler.(OnUpdate[M]); ok {
		if err := onUpdateHandler.BeforeUpdate(id, request); err != nil {
			return err
		}
	}
	err := service.GetEntity().Update(id, request)
	if err != nil {
		err = service.ResponseHandler.Error(responses.ERR_UPDATE_RECORD)
	} else {
		if afterUpdateHandler, ok := service.ServiceHandler.(AfterUpdate[M]); ok {
			err = afterUpdateHandler.AfterUpdate(request)
		}
	}

	return err
}

func (service *Service[M]) Create(request *M) (*M, responses.Error) {
	if err := service.NewValidator(actions.CreateAction, request).ValidateModel(actions.CreateAction, request); err != nil {
		return nil, err
	}

	if onCreateHandler, ok := service.ServiceHandler.(OnCreate[M]); ok {
		if err := onCreateHandler.BeforeCreate(request); err != nil {
			return nil, err
		}
	}
	record, err := service.GetEntity().Create(request)
	if err != nil {
		err = service.ResponseHandler.Error(responses.ERR_CREATE_RECORD)
		record = nil
	} else {
		if afterCreateHandler, ok := service.ServiceHandler.(AfterCreate[M]); ok {
			if err = afterCreateHandler.AfterCreate(record); err != nil {
				record = nil
			}
		}
	}

	return record, err
}

func (service *Service[M]) Delete(id string) responses.Error {
	if err := service.Exists(id); err != nil {
		return err
	}

	if onDeleteHandler, ok := service.ServiceHandler.(OnDelete); ok {
		if err := onDeleteHandler.BeforeDelete(id); err != nil {
			return err
		}
	}
	err := service.GetEntity().Delete(id)
	if err != nil {
		err = service.ResponseHandler.Error(responses.ERR_DELETE_RECORD)
	} else {
		if afterDeleteHandler, ok := service.ServiceHandler.(AfterDelete); ok {
			err = afterDeleteHandler.AfterDelete(id)
		}
	}

	return err
}

func (service *Service[M]) NewValidator(action actions.Action, request *M) *validators.Validator {
	validator := validators.NewValidator()
	service.RegisterValidators(action, request, validator)

	return validator
}
