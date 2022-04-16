package gimlet

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Handler func(c *Context)

type Route struct {
	method  string
	pattern *regexp.Regexp
	handler Handler

	middlewares         []Middleware
	cancelNextExecution bool
}

func newRoute(method string, pattern string, handler Handler) *Route {
	return &Route{
		method:              method,
		pattern:             regexp.MustCompile("^" + pattern + "$"),
		handler:             handler,
		middlewares:         []Middleware{},
		cancelNextExecution: false,
	}
}

func (r *Route) matches(request *http.Request) (bool, [][]string) {
	paramsRegex := regexp.MustCompile("({[a-zA-Z]+})")
	regexParamsRegex := regexp.MustCompile(`({[a-zA-Z]+:\s?.+})`)
	route := paramsRegex.ReplaceAllStringFunc(r.pattern.String(), convertParamToRegex)
	routeRegex := regexp.MustCompile(regexParamsRegex.ReplaceAllStringFunc(route, convertRegexParam))

	path := "/" + strings.Trim(request.URL.Path, "/")
	matches := routeRegex.FindStringSubmatch(path)
	if len(matches) > 0 {
		matches := matches[1:]
		params := routeRegex.SubexpNames()[1:]
		return true, [][]string{params, matches}
	}

	return false, nil
}

func convertRegexParam(param string) string {
	strParam := strings.Trim(strings.Replace(strings.Replace(param, "{", "", 1), "}", "", 1), " ")
	splitRegex := regexp.MustCompile(`:\s?`)
	keys := splitRegex.Split(strParam, -1)
	regex := fmt.Sprintf("(?P<%v>%v)", strings.Trim(keys[0], " "), strings.TrimLeft(keys[1], " "))
	return regex
}

func convertParamToRegex(param string) string {
	strParam := strings.Replace(strings.Replace(param, "{", "", 1), "}", "", 1)
	regex := fmt.Sprintf("(?P<%v>[a-zA-Z0-9_\\-\\.\\!\\~\\*\\\\'\\(\\)\\:\\@\\&\\=\\$\\+,%%{}\"']+)", strParam)
	return regex
}

func (r *Route) CancelExecution() {
	r.cancelNextExecution = true
}

func (r *Route) Handle(ctx *Context) {
	if r.shouldCancelExecution() {
		return
	}
	r.executeMiddlewares(ctx)
	if r.shouldCancelExecution() {
		return
	}

	r.handler(ctx)
}

func (r *Route) shouldCancelExecution() bool {
	if r.cancelNextExecution {
		r.cancelNextExecution = false
		return true
	}

	return false
}

func (r *Route) Use(middleware ...Middleware) *Route {
	r.middlewares = append(r.middlewares, middleware...)
	return r
}

func (r *Route) executeMiddlewares(ctx *Context) {
	for _, middleware := range r.middlewares {
		middleware(r, ctx)
	}
}
