package model

import (
	"net/http"

	"github.com/go-chi/render"
)

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) Context {
	return Context{w, r}
}

func (c Context) Response() http.ResponseWriter {
	return c.w
}

func (c Context) Request() *http.Request {
	return c.r
}

func (c Context) PlainString(status int, v string) {
	render.Status(c.r, status)
	render.PlainText(c.w, c.r, v)
}

func (c Context) JSON(status int, v interface{}) {
	render.Status(c.r, status)
	render.JSON(c.w, c.r, v)
}
