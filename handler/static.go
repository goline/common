package handler

import (
	"io/ioutil"

	"github.com/goline/lapi"
)

type StaticHandler struct {
	File string
}

func (h *StaticHandler) Handle(c lapi.Connection) (interface{}, error) {
	data, err := ioutil.ReadFile(h.File)
	if err != nil {
		return nil, err
	}

	c.Response().
		WithContentType("").
		WithContentBytes(data, nil).
		WithContentType("text/plain")
	return nil, c.Response().Send()
}
