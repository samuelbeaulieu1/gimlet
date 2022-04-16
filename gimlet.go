package gimlet

import (
	"net"
	"net/http"

	"github.com/samuelbeaulieu1/gimlet/logger"
)

type Engine struct {
	router

	Config
	middlewares []Middleware
}

func NewEngine() *Engine {
	engine := &Engine{
		Config: NewConfig(),
	}

	return engine
}

func (engine *Engine) LoadConfig(filePath string) bool {
	return engine.init(filePath)
}

func (engine *Engine) Run() (err error) {
	listener, err := net.Listen("tcp", ":"+engine.Port)
	if err == nil {
		logger.PrintInfo("Listening and serving HTTP on port %v.", engine.Port)
		http.Serve(listener, engine)
	} else {
		logger.PrintError("Error serving HTTP: %v", err)
	}
	return err
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route, ctxParams := engine.resolveRequest(w, req)

	if route == nil {
		http.NotFound(w, req)
	} else if ctxParams != nil {
		ctx := NewContext(w, req, engine, ctxParams)
		engine.executeMiddlewares(route, ctx)
		route.Handle(ctx)
	}
}

func (engine *Engine) Use(middleware ...Middleware) *Engine {
	engine.middlewares = append(engine.middlewares, middleware...)
	return engine
}

func (engine *Engine) executeMiddlewares(route *Route, ctx *Context) {
	for _, middleware := range engine.middlewares {
		middleware(route, ctx)
	}
}
