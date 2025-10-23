package memory

import (
	"context"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_16_calendar/internal/domain"
)

type Storage struct {
	events map[int64]*domain.Event
	users  map[int64]*domain.User
	mu     sync.RWMutex
}

func (s *Storage) CreateEvent(_ context.Context, event *domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	evID := NewID()
	event.ID = evID
	s.events[evID] = event
	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, event *domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[event.ID] = event
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, eventID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, eventID)
	return nil
}

func (s *Storage) GetEvent(_ context.Context, eventID int64) (*domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, ok := s.events[eventID]
	if !ok {
		return nil, domain.ErrEventNotFound
	}
	return event, nil
}

// Отбираем все события по фильтру.
func (s *Storage) GetEvents(_ context.Context, filter domain.EventFilter) ([]*domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]*domain.Event, 0, len(s.events))

	for _, event := range s.events {
		// Фильтр по времени
		if filter.DateFrom != nil && event.StartTime.Before(*filter.DateFrom) {
			continue
		}
		if filter.DateTo != nil && event.StartTime.After(*filter.DateTo) {
			continue
		}

		// Фильтр по пользователю
		if filter.UserID != nil && event.UserID != *filter.UserID {
			continue
		}

		events = append(events, event)
	}

	// Сортируем по дате старта
	sort.Slice(events, func(i, j int) bool {
		return events[i].StartTime.Before(events[j].StartTime)
	})

	return events, nil
}

// --- Методы для работы с User ---
func (s *Storage) CreateUser(_ context.Context, user *domain.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	userID := NewID()
	user.ID = userID
	if s.users == nil {
		s.users = make(map[int64]*domain.User)
	}
	s.users[userID] = user
	return nil
}

func (s *Storage) UpdateUser(_ context.Context, user *domain.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.users == nil {
		s.users = make(map[int64]*domain.User)
	}
	s.users[user.ID] = user
	return nil
}

func (s *Storage) DeleteUser(_ context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.users == nil {
		return nil
	}
	delete(s.users, userID)
	return nil
}

func (s *Storage) GetUser(_ context.Context, userID int64) (*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.users == nil {
		return nil, domain.ErrEventNotFound
	}
	user, ok := s.users[userID]
	if !ok {
		return nil, domain.ErrEventNotFound
	}
	return user, nil
}

func (s *Storage) Close() error {
	return nil
}

func NewID() int64 {
	now := time.Now().UnixNano()       // наносекунды (int64)
	randPart := rand.Int63n(1_000_000) //nolint:gosec // not used for security
	return now*1_000_000 + randPart
}

func New() *Storage {
	return &Storage{events: make(map[int64]*domain.Event)}
}
