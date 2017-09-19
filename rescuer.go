package utils

import (
	"net/http"
	"github.com/goline/lapi"
)

func NewRescuer() lapi.Rescuer {
	return &FactoryRescuer{}
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type FactoryRescuer struct{}

func (h *FactoryRescuer) Rescue(connection lapi.Connection, err error) error {
	if connection == nil {
		return err
	}
	switch e := err.(type) {
	case lapi.SystemError:
		h.handleSystemError(connection, e)
	case lapi.StackError:
		h.handleStackError(connection, e)
	default:
		h.handleUnknownError(connection, e)
	}

	return nil
}

func (h *FactoryRescuer) handleSystemError(c lapi.Connection, err lapi.SystemError) {
	switch err.Code() {
	case lapi.ERROR_HTTP_NOT_FOUND:
		c.Response().WithStatus(http.StatusNotFound).
			WithContent(&ErrorResponse{"ERROR_HTTP_NOT_FOUND", http.StatusText(http.StatusNotFound)})
	case lapi.ERROR_HTTP_BAD_REQUEST:
		c.Response().WithStatus(http.StatusBadRequest).
			WithContent(&ErrorResponse{"ERROR_HTTP_BAD_REQUEST", http.StatusText(http.StatusBadRequest)})
	default:
		c.Response().WithStatus(http.StatusInternalServerError).
			WithContent(&ErrorResponse{"ERROR_INTERNAL_SERVER_ERROR", err.Error()})
	}
}

func (h *FactoryRescuer) handleStackError(c lapi.Connection, err lapi.StackError) {
	c.Response().WithStatus(err.Status()).WithContent(&ErrorResponse{"", err.Error()})
}

func (h *FactoryRescuer) handleUnknownError(c lapi.Connection, err error) {
	if e, ok := err.(lapi.ErrorStatus); ok == true {
		c.Response().WithStatus(e.Status())
	} else {
		c.Response().WithStatus(http.StatusInternalServerError)
	}
	code := "ERROR_UNKNOWN_ERROR"
	if e, ok := err.(lapi.ErrorCoder); ok == true {
		code = e.Code()
	}
	c.Response().WithContent(&ErrorResponse{code, err.Error()})
}