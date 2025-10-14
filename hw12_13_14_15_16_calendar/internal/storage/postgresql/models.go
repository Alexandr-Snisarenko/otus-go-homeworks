package postgresql

import (
	"time"

	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_16_calendar/internal/domain"
)

type eventRow struct {
	ID           int64         `db:"id"`
	Title        string        `db:"title"`
	Description  string        `db:"description"`
	StartTime    time.Time     `db:"start_time"`
	EndTime      time.Time     `db:"end_time"`
	UserID       int64         `db:"user_id"`
	NotifyPeriod time.Duration `db:"notify_period"`
}

func (r eventRow) toDomain() *domain.Event {
	return &domain.Event{
		ID:           r.ID,
		Title:        r.Title,
		Description:  r.Description,
		StartTime:    r.StartTime.UTC(),
		EndTime:      r.EndTime.UTC(),
		UserID:       r.UserID,
		NotifyPeriod: r.NotifyPeriod,
	}
}

func fromDomainEvent(e *domain.Event) eventRow {
	return eventRow{
		ID:           e.ID,
		Title:        e.Title,
		Description:  e.Description,
		StartTime:    e.StartTime, // предполагаем, что уже в UTC на слое app
		EndTime:      e.EndTime,
		UserID:       e.UserID,
		NotifyPeriod: e.NotifyPeriod,
	}
}

func toDomainEvents(rows []eventRow) []*domain.Event {
	events := make([]*domain.Event, 0, len(rows))
	for _, r := range rows {
		events = append(events, r.toDomain())
	}
	return events
}

type userRow struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func (r userRow) toDomain() *domain.User {
	return &domain.User{
		ID:    r.ID,
		Name:  r.Name,
		Email: r.Email,
	}
}

func fromDomainUser(u *domain.User) userRow {
	return userRow{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
