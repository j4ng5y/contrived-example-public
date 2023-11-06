package storage

import (
	"context"
	"fmt"
)

type (
	RepoType int
	Repo     interface {
		Connect(ctx context.Context, uri string) error
		Insert(message, err_message string) error
	}
)

const (
	RepoType_Memory RepoType = iota
	RepoType_Sqlite
)

func NewRepo(rt RepoType) (Repo, error) {
	switch {
	case rt == RepoType_Memory:
		return &Memory{}, nil
	case rt == RepoType_Sqlite:
		return &Sqlite{}, nil
	default:
		return nil, fmt.Errorf("invalid repo type")
	}
}
