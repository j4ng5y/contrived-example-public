package storage

import (
	"context"
)

type (
	repo struct {
		message       string
		error_message string
	}
	Memory struct {
		repo map[uint]repo
	}
)

func (r *Memory) Connect(ctx context.Context, uri string) error {
	return nil
}

func (r *Memory) Insert(message, err_message string) error {
	if len(r.repo) == 0 {
		r.repo[0] = repo{
			message,
			err_message,
		}
		return nil
	} else {
		r.repo[uint(len(r.repo))] = repo{
			message,
			err_message,
		}
		return nil
	}
}
