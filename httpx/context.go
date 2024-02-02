package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/charliego3/pallas/encoding"
	"github.com/charliego3/pallas/encoding/json"
	"github.com/charliego3/pallas/encoding/xml"
	"github.com/charliego3/pallas/utility"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

const (
	contentTypeHeader = "Content-Type"
)

var defaultCodecType = json.Type

func SetDefaultCodeType(ctype string) {
	defaultCodecType = ctype
}

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

func (c *Context) bind(contentType string, v any) error {
	codec := encoding.CodecWithType(contentType)
	if coder, ok := codec.(encoding.Coder); ok {
		return coder.Decoder(c.Body).Decode(v)
	}

	b, err := io.ReadAll(c.Body)
	if err != nil {
		return err
	}
	return codec.Unmarshal(b, v)
}

func (c *Context) BindJSON(v any) error {
	return c.bind(json.Type, v)
}

func (c *Context) BindXML(v any) error {
	return c.bind(xml.Type, v)
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

	contentType := c.Header.Get(contentTypeHeader)
	contentType = SubContentType(contentType)
	switch contentType {
	case "x-www-form-urlencoded":
		return c.BindForm(v)
	default:
		return c.bind(contentType, v)
	}
}

func (c *Context) write(contentType string, v any, code []int) error {
	codec := encoding.CodecWithType(contentType)
	c.Writer.WriteHeader(utility.First(http.StatusOK, code))
	c.Writer.Header().Set("Content-Type", "application/"+codec.Type())
	if coder, ok := codec.(encoding.Coder); ok {
		return coder.Encoder(c.Writer).Encode(v)
	}

	b, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(c.Writer, b)
	return err
}

func (c *Context) Write(v any, code ...int) error {
	contentType := SubContentType(c.Header.Get("Accept"))
	return c.write(contentType, v, code)
}

func (c *Context) JSON(v any, code ...int) error {
	return c.write(json.Type, v, code)
}

func (c *Context) XML(v any, code ...int) error {
	return c.write(xml.Type, v, code)
}

func (c *Context) Aborted(msg ...string) error {
	return AbortedErr{msg}
}

func SubContentType(contentType string) string {
	if len(contentType) == 0 {
		return defaultCodecType
	}
	first := strings.Index(contentType, "/")
	if first == -1 {
		return defaultCodecType
	}
	last := strings.Index(contentType, ";")
	if last == -1 {
		return contentType[first+1:]
	}
	if last <= first {
		return defaultCodecType
	}
	return contentType[first+1 : last]
}
