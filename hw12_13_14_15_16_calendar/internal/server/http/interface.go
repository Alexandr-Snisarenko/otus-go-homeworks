package http

import (
	"context"

	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_16_calendar/internal/domain"
)

// Application описывает контракт прикладного уровня, используемый HTTP-сервером.
// Интерфейс зеркалирует методы, реализованные в internal/app.App.
type Application interface {
	// Events
	CreateEvent(ctx context.Context, e domain.Event) (int64, error)
	UpdateEvent(ctx context.Context, e domain.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetEvent(ctx context.Context, id int64) (*domain.Event, error)
	GetEvents(ctx context.Context, f domain.EventFilter) ([]*domain.Event, error)

	// Users
	CreateUser(ctx context.Context, u domain.User) (int64, error)
	UpdateUser(ctx context.Context, u domain.User) error
	DeleteUser(ctx context.Context, id int64) error
	GetUser(ctx context.Context, id int64) (*domain.User, error)
}
