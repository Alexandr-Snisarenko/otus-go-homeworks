package app

import (
	"context"

	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_16_calendar/internal/domain"
)

type App struct {
	storage Storage
	logger  Logger
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

// Events
func (a *App) CreateEvent(ctx context.Context, e domain.Event) (int64, error) {
	if err := a.storage.CreateEvent(ctx, &e); err != nil {
		return 0, err
	}
	return e.ID, nil
}

func (a *App) UpdateEvent(ctx context.Context, e domain.Event) error {
	return a.storage.UpdateEvent(ctx, &e)
}

func (a *App) DeleteEvent(ctx context.Context, id int64) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetEvent(ctx context.Context, id int64) (*domain.Event, error) {
	return a.storage.GetEvent(ctx, id)
}

func (a *App) GetEvents(ctx context.Context, f domain.EventFilter) ([]*domain.Event, error) {
	return a.storage.GetEvents(ctx, f)
}

// Users
func (a *App) CreateUser(ctx context.Context, usr domain.User) (int64, error) {
	if err := a.storage.CreateUser(ctx, &usr); err != nil {
		return 0, err
	}
	return usr.ID, nil
}

func (a *App) UpdateUser(ctx context.Context, usr domain.User) error {
	return a.storage.UpdateUser(ctx, &usr)
}

func (a *App) DeleteUser(ctx context.Context, id int64) error {
	return a.storage.DeleteUser(ctx, id)
}

func (a *App) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	return a.storage.GetUser(ctx, id)
}
