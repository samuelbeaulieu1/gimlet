package gimlet

import (
	"net/http"
	"strings"
)

type GroupHandler func(r Router)

type router struct {
	routes []*Route

	engine *Engine

	rootPattern string
	groups      []*router
}

type Router interface {
	GET(pattern string, handler Handler) *Route
	POST(pattern string, handler Handler) *Route
	DELETE(pattern string, handler Handler) *Route
	PUT(pattern string, handler Handler) *Route

	Method(method string, pattern string, handler Handler) *Route

	addRoute(method string, pattern string, handler Handler) *Route

	Group(pattern string, h GroupHandler)

	resolveRequest(w http.ResponseWriter, req *http.Request) (*Route, *ContextParams)
	resolveGroups(w http.ResponseWriter, req *http.Request) (*Route, *ContextParams)
}

func NewRouter(engine *Engine) *router {
	return &router{
		engine:      engine,
		rootPattern: "",
	}
}

func (r *router) GET(pattern string, handler Handler) *Route {
	return r.addRoute(http.MethodGet, pattern, handler)
}

func (r *router) POST(pattern string, handler Handler) *Route {
	return r.addRoute(http.MethodPost, pattern, handler)
}

func (r *router) DELETE(pattern string, handler Handler) *Route {
	return r.addRoute(http.MethodDelete, pattern, handler)
}

func (r *router) PUT(pattern string, handler Handler) *Route {
	return r.addRoute(http.MethodPut, pattern, handler)
}

func (r *router) Method(method string, pattern string, handler Handler) *Route {
	return r.addRoute(method, pattern, handler)
}

func (r *router) addRoute(method string, pattern string, handler Handler) *Route {
	pattern = "/" + strings.Trim(r.rootPattern+"/"+strings.Trim(pattern, "/"), "/")
	newRoute := newRoute(method, pattern, handler)
	r.routes = append(r.routes, newRoute)

	return newRoute
}

func (r *router) Group(pattern string, handler GroupHandler) {
	newRouter := NewRouter(r.engine)
	newRouter.rootPattern = strings.Trim(r.rootPattern+pattern, "/")

	r.groups = append(r.groups, newRouter)
	handler(newRouter)
}

func (r *router) resolveRequest(w http.ResponseWriter, req *http.Request) (*Route, *ContextParams) {
	var route *Route
	for _, currentRoute := range r.routes {
		matches, params := currentRoute.matches(req)
		if matches {
			if currentRoute.method != req.Method && req.Method != http.MethodOptions {
				route = currentRoute
				continue
			}
			return currentRoute, NewContextParams(params[0], params[1])
		}
	}
	if route != nil {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return route, nil
	}
	return r.resolveGroups(w, req)
}

func (r *router) resolveGroups(w http.ResponseWriter, req *http.Request) (*Route, *ContextParams) {
	for _, group := range r.groups {
		if route, ctxParams := group.resolveRequest(w, req); route != nil {
			return route, ctxParams
		}
	}
	return nil, nil
}
