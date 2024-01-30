package httpx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/charliego3/pallas/utility"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

const (
	contentTypeHeader = "Content-Type"
)

type AbortedErr struct {
	msg []string
}

func (e AbortedErr) Error() string {
	msg := "Aborted"
	if len(e.msg) > 0 {
		msg += strings.Join(e.msg, ". ")
	}
	return msg
}

type Context struct {
	context.Context
	*http.Request
	Writer  http.ResponseWriter
	decoder *schema.Decoder
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := new(Context)
	ctx.Context = r.Context()
	ctx.Request = r
	ctx.Writer = w
	ctx.decoder = schema.NewDecoder()
	return ctx
}

func (c *Context) BindJSON(v any) error {
	return nil
}

func (c *Context) BindXML(v any) error {
	return nil
}

func (c *Context) BindQuery(v any) error {
	values := c.URL.Query()
	if len(values) == 0 {
		return nil
	}

	return c.decoder.Decode(v, values)
}

func (c *Context) BindVars(v any) error {
	vars := mux.Vars(c.Request)
	if len(vars) == 0 {
		return nil
	}

	values := make(url.Values)
	for k, v := range vars {
		values.Set(k, v)
	}
	return c.decoder.Decode(v, values)
}

func (c *Context) BindForm(v any) error {
	if len(c.PostForm) == 0 {
		return nil
	}

	return c.decoder.Decode(v, c.PostForm)
}

func (c *Context) Bind(v any) error {
	if err := c.BindVars(v); err != nil {
		return err
	}

	if c.Method == http.MethodGet {
		return c.BindQuery(v)
	}

	contentType := c.Request.Header.Get(contentTypeHeader)
	switch contentType {

	}
	return nil
}

func (c *Context) JSON(v any, code ...int) error {
	c.Writer.Header().Set(contentTypeHeader, "application/json")
	c.Writer.WriteHeader(utility.First(http.StatusOK, code))
	return json.NewEncoder(c.Writer).Encode(v)
}

func (c *Context) Aborted(msg ...string) error {
	return AbortedErr{msg}
}
