package gimlet

type Middleware func(route *Route, ctx *Context)
