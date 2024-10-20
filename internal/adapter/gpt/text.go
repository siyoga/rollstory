package gpt

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

func (a *adapter) CreateThread(
	ctx context.Context,
) (openai.Thread, error) {
	resp, err := a.client.CreateThread(
		ctx,
		openai.ThreadRequest{},
	)
	if err != nil {
		return openai.Thread{}, a.log.AdapterError(err, errors.ErrAdapterRequestFailed, "create_thread")
	}

	return resp, nil
}

// TODO: fix similar error details
func (a *adapter) Request(ctx context.Context, threadId string, msg string, respLimit int, respOrder domain.ReturnOrder) (openai.MessagesList, error) {
	if _, err := a.client.CreateMessage(
		ctx,
		threadId,
		openai.MessageRequest{
			Role:    "user",
			Content: msg,
		},
	); err != nil {
		return openai.MessagesList{}, a.log.AdapterError(err, errors.ErrAdapterRequestFailed, "create_message")
	}

	runResp, err := a.client.CreateRun(ctx, threadId, openai.RunRequest{
		AssistantID: a.cfg.Assistants[0],
	})
	if err != nil {
		return openai.MessagesList{}, a.log.AdapterError(err, errors.ErrAdapterRequestFailed, "create_run")
	}

	for {
		select {
		case <-ctx.Done():
			return openai.MessagesList{}, a.log.AdapterError(
				fmt.Errorf("context deadline exceeded"),
				errors.ErrAdapterRequestFailed,
				"retrieve_run",
			)
		default:
			resp, err := a.client.RetrieveRun(ctx, threadId, runResp.ID)
			if err != nil {
				return openai.MessagesList{}, a.log.AdapterError(err, errors.ErrAdapterRequestFailed, "retrieve_run")
			}

			if slices.Contains(domain.FailedStatuses, resp.Status) {
				return openai.MessagesList{}, a.log.AdapterError(
					fmt.Errorf("request failed with status %s", resp.Status),
					errors.ErrAdapterRequestFailed,
					"retrieve_run",
				)
			}

			if resp.Status == openai.RunStatusCompleted {
				order := string(respOrder)
				messages, err := a.client.ListMessage(ctx, threadId, &respLimit, &order, nil, nil, &runResp.ID)
				if err != nil {
					return openai.MessagesList{}, a.log.AdapterError(err, errors.ErrAdapterRequestFailed, "retrieve_run")
				}

				return messages, nil
			}

			time.Sleep(1 * time.Second)
		}
	}
}
