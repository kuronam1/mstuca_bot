package botErrors

import (
	"errors"
	"log/slog"
)

var (
	ErrUserNotPresents = errors.New("no user in cache")
	ErrNoGroupsFound   = errors.New("группы не найдены, попробуйте еще раз")
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
