package render

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"text/template"

	"github.com/x-mod/errors"
)

type ResponseOpt func(http.ResponseWriter)

func StatusCode(code int) ResponseOpt {
	return func(wr http.ResponseWriter) {
		wr.WriteHeader(code)
	}
}

func Cookie(cookie *http.Cookie) ResponseOpt {
	return func(wr http.ResponseWriter) {
		http.SetCookie(wr, cookie)
	}
}

func HeaderAdd(key string, value string) ResponseOpt {
	return func(wr http.ResponseWriter) {
		wr.Header().Add(key, value)
	}
}
func HeaderSet(key string, value string) ResponseOpt {
	return func(wr http.ResponseWriter) {
		wr.Header().Set(key, value)
	}
}
func HeaderDel(key string) ResponseOpt {
	return func(wr http.ResponseWriter) {
		wr.Header().Del(key)
	}
}

//Responsor interface
type Responsor interface {
	Response(http.ResponseWriter, ...ResponseOpt) error
}

type JSONRender struct {
	data interface{}
}

func JSON(data interface{}) *JSONRender {
	return &JSONRender{data: data}
}

func (r *JSONRender) Response(wr http.ResponseWriter, opts ...ResponseOpt) error {
	wr.Header().Set("Content-Type", "application/json; charset=utf-8")
	for _, opt := range opts {
		opt(wr)
	}
	return json.NewEncoder(wr).Encode(r.data)
}

type ErrorRender struct {
	err error
}

func Error(err error) *ErrorRender {
	return &ErrorRender{err: err}
}

func (r *ErrorRender) Response(wr http.ResponseWriter, opts ...ResponseOpt) error {
	return JSON(map[string]interface{}{
		"code":    errors.ValueFrom(r.err),
		"message": r.err.Error(),
	}).Response(wr, opts...)
}

type TextRender struct {
	text string
}

func Text(txt string) *TextRender {
	return &TextRender{text: txt}
}

func (r *TextRender) Response(wr http.ResponseWriter, opts ...ResponseOpt) error {
	wr.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, opt := range opts {
		opt(wr)
	}
	_, err := wr.Write([]byte(r.text))
	return err
}

type XMLRender struct {
	data interface{}
}

func XML(data interface{}) *XMLRender {
	return &XMLRender{
		data: data,
	}
}

func (r *XMLRender) Response(wr http.ResponseWriter, opts ...ResponseOpt) error {
	wr.Header().Set("Content-Type", "application/xml; charset=utf-8")
	for _, opt := range opts {
		opt(wr)
	}
	return xml.NewEncoder(wr).Encode(r.data)
}

type TemplateRender struct {
	tpl  *template.Template
	data interface{}
}

func Template(tpl *template.Template, data interface{}) *TemplateRender {
	return &TemplateRender{
		tpl:  tpl,
		data: data,
	}
}

func (r *TemplateRender) Response(wr http.ResponseWriter, opts ...ResponseOpt) error {
	wr.Header().Set("Content-Type", "application/html; charset=utf-8")
	for _, opt := range opts {
		opt(wr)
	}
	return r.tpl.Execute(wr, r.data)
}
