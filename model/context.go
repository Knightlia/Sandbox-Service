package model

import (
	"net/http"

	"github.com/go-chi/render"
)

type Context struct {
	writer  http.ResponseWriter
	request *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) Context {
	return Context{w, r}
}

func (c Context) Response() http.ResponseWriter {
	return c.writer
}

func (c Context) Request() *http.Request {
	return c.request
}

func (c Context) Header(k string) string {
	return c.request.Header.Get(k)
}

func (c Context) PlainString(status int, v string) {
	render.Status(c.request, status)
	render.PlainText(c.writer, c.request, v)
}

func (c Context) JSON(status int, v interface{}) {
	render.Status(c.request, status)
	render.JSON(c.writer, c.request, v)
}
