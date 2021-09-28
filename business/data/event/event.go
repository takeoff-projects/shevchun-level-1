package event

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return Store{
		db: db,
	}
}

func (s Store) GetEvents(ctx context.Context) ([]Event, error) {
	const q = `SELECT id, title, location, event_date FROM events`
	rows, err := s.db.NamedQueryContext(ctx, q, struct{}{})
	if err != nil {
		return []Event{}, err
	}
	var events []Event
	for rows.Next() {
		var event Event
		rows.StructScan(&event)
		events = append(events, event)
	}
	rows.Close()
	return events, nil
}

func (s Store) GetEventByID(ctx context.Context, id string) (Event, error) {
	data := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	const q = `SELECT id, title, location, event_date FROM events WHERE id = :id`
	var event Event
	rows, err := s.db.NamedQueryContext(ctx, q, data)
	if err != nil {
		return Event{}, err
	}
	if !rows.Next() {
		return Event{}, fmt.Errorf("no events")
	}
	if err := rows.StructScan(&event); err != nil {
		return Event{}, err
	}
	return event, nil
}

func (s Store) AddEvent(ctx context.Context, event Event) error {
	newID := uuid.New().String()
	event.ID = newID
	const q = `
		INSERT INTO events
			(id, title, location, event_date) 
		VALUES
			(:id, :title, :location, :event_date)
	`
	if _, err := s.db.NamedExecContext(ctx, q, event); err != nil {
		return err
	}
	return nil
}

func (s Store) UpdateEvent(ctx context.Context, updatedEvent Event) error {
	const q = `UPDATE events SET title=:title, location=:location, event_date=:event_date WHERE id=:id`
	if _, err := s.db.NamedExecContext(ctx, q, updatedEvent); err != nil {
		return err
	}
	return nil
}

func (s Store) DeleteEvent(ctx context.Context, id string) error {
	data := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}
	const q = `DELETE FROM events WHERE id = :id`
	if _, err := s.db.NamedExecContext(ctx, q, data); err != nil {
		return err
	}
	return nil
}
