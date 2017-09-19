package handler

import (
	"github.com/goline/lapi"
	"io/ioutil"
)

type StaticHandler struct {
	File string
}

func (h *StaticHandler) Handle(c lapi.Connection) (interface{}, error) {
	data, err := ioutil.ReadFile(h.File)
	if err != nil {
		return nil, err
	}
	c.Response().WithContentType("").WithContentBytes(data, nil)
	c.Response().Send()
	return nil, nil
}
