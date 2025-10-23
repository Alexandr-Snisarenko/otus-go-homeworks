package app

import (
	"context"

	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_16_calendar/internal/domain"
)

type Logger interface {
	Info(msg string, v ...any)
	Error(msg string, v ...any)
	Debug(msg string, v ...any)
	Warn(msg string, v ...any)

	InfoCtx(ctx context.Context, msg string, v ...any)
	ErrorCtx(ctx context.Context, msg string, v ...any)
	DebugCtx(ctx context.Context, msg string, v ...any)
	WarnCtx(ctx context.Context, msg string, v ...any)
}

type Storage interface {
	EventStorage
	UserStorage
}

type EventStorage interface {
	CreateEvent(context.Context, *domain.Event) error
	UpdateEvent(context.Context, *domain.Event) error
	DeleteEvent(context.Context, int64) error
	GetEvent(context.Context, int64) (*domain.Event, error)
	GetEvents(context.Context, domain.EventFilter) ([]*domain.Event, error)
}

type UserStorage interface {
	CreateUser(context.Context, *domain.User) error
	UpdateUser(context.Context, *domain.User) error
	DeleteUser(context.Context, int64) error
	GetUser(context.Context, int64) (*domain.User, error)
}
