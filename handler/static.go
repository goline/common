package handler

import (
	"io/ioutil"

	"github.com/goline/lapi"
)

type StaticHandler struct {
	File        string
	ContentType string
}

func (h *StaticHandler) Handle(c lapi.Connection) (interface{}, error) {
	data, err := ioutil.ReadFile(h.File)
	if err != nil {
		return nil, err
	}

	contentType := h.ContentType
	if contentType == "" {
		contentType = "text/plain"
	}

	c.Response().
		WithContentType("").
		WithContentBytes(data, nil).
		WithContentType(contentType)
	return nil, c.Response().Send()
}
