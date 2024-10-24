package router

import (
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

// Response
func newSuccessResponse(res domain.Message, code int, release *int64) response {
	return response{
		result:  res,
		code:    code,
		error:   nil,
		release: release,
	}
}

func newErrResponse(e *errors.Error, release int64) response {
	return response{
		result:  domain.Message{},
		code:    int(e.Code),
		error:   e,
		release: &release,
	}
}
