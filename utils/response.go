package utils

import (
	"github.com/goccy/go-json"
	"net/http"
)

type ResponseEntity interface {
	Code() int
	Error() string
	Marshal() ([]byte, error)
	MarshalIdent() ([]byte, error)
}

type re struct {
	Codes   int    `json:"code"`
	Msg     string `json:"message"`
	Payload any    `json:"payload,omitempty"`
}

func (r re) Code() int {
	return r.Codes
}

func (r re) Error() string {
	return r.Msg
}

func (r re) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r re) MarshalIdent() ([]byte, error) {
	return json.MarshalIndent(r, "", "    ")
}

func NewErrResponse(code int, msg string) ResponseEntity {
	return &re{Codes: code, Msg: msg}
}

func NewResponse(payload any) ResponseEntity {
	return &re{Codes: http.StatusOK, Msg: "OK", Payload: payload}
}
