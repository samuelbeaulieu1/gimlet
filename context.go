package gimlet

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
	"github.com/samuelbeaulieu1/gimlet/logger"
	"github.com/samuelbeaulieu1/gimlet/responses"
)

type Context struct {
	Request *http.Request
	Params  *ContextParams
	Writer  http.ResponseWriter

	Engine *Engine

	body    url.Values
	decoder *schema.Decoder

	Authentication *AuthTokenPayload
}

func NewContext(w http.ResponseWriter, r *http.Request, engine *Engine, cp *ContextParams) *Context {
	ctx := &Context{
		Writer:  w,
		Request: r,
		Params:  cp,
		Engine:  engine,
		decoder: schema.NewDecoder(),
	}
	ctx.decoder.SetAliasTag("json")
	ctx.parsePostForm()

	return ctx
}

func (ctx *Context) parsePostForm() {
	if err := ctx.Request.ParseMultipartForm(ctx.Engine.MaxPostFormSizeMB); err != nil {
		if err != http.ErrNotMultipart {
			logger.PrintDebug("Error parsing post form: %v", err)
		}
	}
	ctx.body = ctx.Request.PostForm
}

func (ctx *Context) ParseBody(s any) bool {
	var err error
	if ctx.Request.Header.Get("Content-Type") == "application/json" {
		err = json.NewDecoder(ctx.Request.Body).Decode(s)
	} else {
		err = ctx.decoder.Decode(s, ctx.body)
	}

	return err != nil
}

func (ctx *Context) GetBodyParam(key string) string {
	if val := ctx.GetBodyParams(key); len(val) > 0 {
		return val[0]
	}

	return ""
}

func (ctx *Context) GetBodyParams(key string) []string {
	if val, ok := ctx.body[key]; ok {
		return val
	}

	return []string{}
}

func (ctx *Context) GetParam(key string) string {
	return ctx.Params.Get(key)
}

func (ctx *Context) Status(code int) {
	ctx.Writer.WriteHeader(code)
}

func (ctx *Context) WriteMessage(res string) {
	ctx.Writer.Write([]byte(res))
}

func (ctx *Context) SetJSONContent() {
	ctx.Writer.Header().Set("Content-Type", "application/json")
}

func (ctx *Context) WriteError(code int, message ...string) {
	ctx.SetJSONContent()

	ctx.Status(code)
	if len(message) > 0 {
		ctx.WriteMessage(message[0])
	}
}

func (ctx *Context) WriteJSON(data any) {
	ctx.SetJSONContent()

	json.NewEncoder(ctx.Writer).Encode(data)
}

func (ctx *Context) WriteJSONResponse(data any) {
	ctx.SetJSONContent()

	json.NewEncoder(ctx.Writer).Encode(&responses.RequestResponse{
		Status: "success",
		Data:   data,
	})
}

func (ctx *Context) WriteJSONError(code int, err ...any) {
	ctx.SetJSONContent()

	ctx.Status(code)
	if len(err) > 0 {
		ctx.WriteJSON(err[0])
	}
}
