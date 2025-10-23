package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_16_calendar/internal/config"
	"github.com/Alexandr-Snisarenko/otus-go-homeworks/hw12_13_14_15_16_calendar/internal/domain"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v4/stdlib" // register pgx driver
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

// Создает нового пользователя и обновляет ID созданного пользователя в структуре user.
func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	const query = `
		INSERT INTO users (name, email)
		VALUES (:name, :email)
		RETURNING id`
	stmt, args, _ := sqlx.Named(query, fromDomainUser(user))
	stmt = s.db.Rebind(stmt)
	return s.db.QueryRowContext(ctx, stmt, args...).Scan(&user.ID)
}

func (s *Storage) UpdateUser(ctx context.Context, user *domain.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	const query = `
		UPDATE users SET
			name = :name,
			email = :email
		WHERE id = :id
	`
	res, err := s.db.NamedExecContext(ctx, query, fromDomainUser(user))
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, userID int64) error {
	const query = `DELETE FROM users WHERE id = $1`
	res, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (s *Storage) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	const query = `
		SELECT id, name, email
		FROM users
		WHERE id = $1
	`
	var user userRow
	if err := s.db.GetContext(ctx, &user, query, id); err != nil {
		return nil, err
	}
	return user.toDomain(), nil
}

// Создает новое событие и обновляет ID созданного события в структуре event.
func (s *Storage) CreateEvent(ctx context.Context, event *domain.Event) error {
	if event == nil {
		return domain.ErrEventIsEmpty
	}

	const query = `
    INSERT INTO events (title, description, start_time, end_time, user_id, notify_period)
    VALUES (:title, :description, :start_time, :end_time, :user_id, :notify_period)
    RETURNING id`

	stmt, args, _ := sqlx.Named(query, fromDomainEvent(event))
	stmt = s.db.Rebind(stmt)

	return s.db.QueryRowContext(ctx, stmt, args...).Scan(&event.ID)
}

func (s *Storage) UpdateEvent(ctx context.Context, event *domain.Event) error {
	if event == nil {
		return errors.New("event is nil")
	}

	const query = `
        UPDATE events SET
            title = :title,
            description = :description,
            start_time = :start_time,
            end_time = :end_time,
            user_id = :user_id,
            notify_period = :notify_period
        WHERE id = :id
    `

	res, err := s.db.NamedExecContext(ctx, query, fromDomainEvent(event))
	if err != nil {
		return err
	}

	// проверяем был ли UPDATE
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID int64) error {
	const query = `DELETE FROM events WHERE id = $1`

	res, err := s.db.ExecContext(ctx, query, eventID)
	if err != nil {
		return err
	}

	// проверяем был ли DELETE
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrEventNotFound
	}

	return err
}

func (s *Storage) GetEvent(ctx context.Context, id int64) (*domain.Event, error) {
	const query = `
        SELECT id, title, description, start_time, end_time, user_id, notify_period
        FROM events
        WHERE id = $1
    `

	var event eventRow
	if err := s.db.GetContext(ctx, &event, query, id); err != nil {
		return nil, err
	}

	return event.toDomain(), nil
}

func (s *Storage) GetEvents(ctx context.Context, filter domain.EventFilter) ([]*domain.Event, error) {
	// Собираем запрос через squirrel. Заполняем фильтрацию, если параметры фильтра не nil
	qb := sq.
		Select("id", "title", "description", "start_time", "end_time", "user_id", "notify_period").
		From("events").
		PlaceholderFormat(sq.Dollar) // $1, $2...

	if filter.DateFrom != nil {
		qb = qb.Where(sq.GtOrEq{"start_time": *filter.DateFrom})
	}
	if filter.DateTo != nil {
		qb = qb.Where(sq.LtOrEq{"start_time": *filter.DateTo})
	}
	if filter.UserID != nil {
		qb = qb.Where(sq.Eq{"user_id": *filter.UserID})
	}

	// Соритруем по дате старта (по умолчанию)
	qb = qb.OrderBy("start_time")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var events []eventRow
	if err := s.db.SelectContext(ctx, &events, query, args...); err != nil {
		return nil, err
	}
	return toDomainEvents(events), nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func New(cfg config.Database) (*Storage, error) {
	db, err := OpenDB(cfg)
	if err != nil {
		return nil, err
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(cfg.Postgresql.Pool.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Postgresql.Pool.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Postgresql.Pool.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.Postgresql.Pool.ConnMaxIdleTime)

	return &Storage{db: db}, nil
}

func OpenDB(cfg config.Database) (*sqlx.DB, error) {
	// Работает только для режима postgresql
	if cfg.Workmode != "postgresql" {
		return nil, errors.New("workmode must be 'postgresql'")
	}

	// Если заполнен параметр конфиге "dns", то используем его. В этом случае параметры User и т.д. - игнорируются
	dsn := cfg.Postgresql.Dsn
	if dsn == "" {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.Postgresql.User, cfg.Postgresql.Password,
			cfg.Postgresql.Host, cfg.Postgresql.Port, cfg.Postgresql.Name)
		if dsn == "" {
			return nil, errors.New("empty DSN")
		}
	}

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
