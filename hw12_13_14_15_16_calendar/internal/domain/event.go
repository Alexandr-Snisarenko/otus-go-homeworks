package domain

import "time"

// Event - структура для объекта "Событие"
// ID - уникальный идентификатор события (BIGINT);
// Заголовок - короткий текст;
// Дата и время события;
// Длительность события (или дата и время окончания);
// Описание события - длинный текст, опционально;
// ID пользователя, владельца события;
// За сколько времени высылать уведомление, опционально.
type Event struct {
	ID           int64
	Title        string
	Description  string
	StartTime    time.Time
	EndTime      time.Time
	UserID       int64
	NotifyPeriod time.Duration
}

type EventFilter struct {
	UserID   *int64
	DateFrom *time.Time
	DateTo   *time.Time
}
