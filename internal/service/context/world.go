package context

import (
	"context"
	"fmt"

	"github.com/siyoga/rollstory/internal/domain"
	"github.com/siyoga/rollstory/internal/errors"
)

func (s *service) CreateWorld(ctx context.Context, userId int64, worldDesc string) (string, *errors.Error) {
	threadId, err := s.threadRepository.GetThreadByUser(ctx, userId)
	if err != nil {
		return "", errors.DatabaseError(err)
	}

	worldDesc = fmt.Sprintf("Описание мира:\n%s", worldDesc)

	resp, err := s.gptAdapter.Request(ctx, threadId, worldDesc, 1, domain.Asc)
	if err != nil {
		return "", errors.AdapterError(err)
	}

	return resp.Messages[0].Content[0].Text.Value, nil
}
