package gimlet

import (
	"net/http"
	"strconv"

	"github.com/samuelbeaulieu1/gimlet/responses"
)

const RouteIdentifier = "id"

type ControllerHandler[M Model] interface {
	RegisterRoutes(router Router)
	GetService() ServiceInterface[M]
}

type ControllerInterface[M Model] interface {
	GetAll(*Context)
	Get(*Context)
	Update(*Context)
	Delete(*Context)
	Create(*Context)
	ControllerHandler[M]
}

type Controller[M Model] struct {
	ControllerHandler[M]
	responseHandler responses.ResponseHandler
}

func NewController[M Model]() *Controller[M] {
	return &Controller[M]{
		responseHandler: &responses.ResponseMapper{},
	}
}

func ParseModelUintID(id string) (uint, error) {
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsedId), nil
}

func (controller *Controller[M]) ParseRouteUintIdentifier(key string, ctx *Context) (uint, bool) {
	if ctx.GetParam(key) == "" {
		ctx.WriteJSONError(http.StatusBadRequest, controller.responseHandler.Error(responses.ERR_EMPTY_ID))
		return 0, false
	}
	id, err := ParseModelUintID(ctx.GetParam(key))
	if err != nil {
		ctx.WriteJSONError(http.StatusBadRequest, controller.responseHandler.Error(responses.ERR_INVALID_ID))
		return 0, false
	}

	return id, true
}

func (controller *Controller[M]) ParseRouteStrIdentifier(key string, ctx *Context) (string, bool) {
	if ctx.GetParam(key) == "" {
		ctx.WriteJSONError(http.StatusBadRequest, controller.responseHandler.Error(responses.ERR_EMPTY_ID))
		return "", false
	}

	return ctx.GetParam(key), true
}

func (controller *Controller[M]) GetAll(ctx *Context) {
	records, err := controller.GetService().GetAll()
	if err != nil {
		ctx.WriteJSONError(http.StatusBadRequest, err)
		return
	}

	if len(*records) > 0 {
		ctx.WriteJSONResponse(ParseModelsToDTO(records, (*records)[0].ToDTO()))
	} else {
		ctx.WriteJSONResponse([]M{})
	}
}

func (controller *Controller[M]) Get(ctx *Context) {
	id, ok := controller.ParseRouteStrIdentifier(RouteIdentifier, ctx)
	if !ok {
		return
	}
	record, err := controller.GetService().Get(id)
	if err != nil {
		ctx.WriteJSONError(http.StatusBadRequest, err)
		return
	}

	ctx.WriteJSONResponse((*record).ToDTO())
}

func (controller *Controller[M]) Update(ctx *Context) {
	id, ok := controller.ParseRouteStrIdentifier(RouteIdentifier, ctx)
	if !ok {
		return
	}
	var request M
	ctx.ParseBody(&request)

	err := controller.GetService().Update(id, &request)
	if err != nil {
		ctx.WriteJSONError(http.StatusBadRequest, err)
	} else {
		ctx.WriteJSONResponse(&responses.RequestResponseMessage{
			Message: controller.responseHandler.String(responses.SUCC_UPDATE_RECORD),
		})
	}
}

func (controller *Controller[M]) Delete(ctx *Context) {
	id, ok := controller.ParseRouteStrIdentifier(RouteIdentifier, ctx)
	if !ok {
		return
	}

	err := controller.GetService().Delete(id)
	if err != nil {
		ctx.WriteJSONError(http.StatusBadRequest, err)
	} else {
		ctx.WriteJSONResponse(&responses.RequestResponseMessage{
			Message: controller.responseHandler.String(responses.SUCC_DELETE_RECORD),
		})
	}
}

func (controller *Controller[M]) Create(ctx *Context) {
	var request M
	ctx.ParseBody(&request)

	record, err := controller.GetService().Create(&request)
	if err != nil {
		ctx.WriteJSONError(http.StatusBadRequest, err)
	} else {
		ctx.WriteJSONResponse((*record).ToDTO())
	}
}
