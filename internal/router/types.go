package router

import (
	"context"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

type (
	Route struct {
		name    string
		handler cmdHandler
		buttons []button
	}

	button struct {
		text string
	}

	response struct {
		result domain.Message
		code   int
		error  *errors.Error

		// pass the user ID of the person whose execution of the command has finished
		release *int64
	}

	cmdHandler func(ctx context.Context, userId int64, msg domain.Message) response
)
