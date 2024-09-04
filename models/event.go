package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/LeonLonsdale/go-web-api/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

var events = []Event{}

func (e *Event) Save() error {
	// Prepare the SQL statement for execution
	// Preparing statements can improve performance, especially if the statement is executed multiple times.
	query := `
    INSERT INTO events (name, description, location, dateTime, user_id)
    VALUES (?,?,?,?,?)
  `

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	// a statement uses up resources as it sits in a ready state waiting to be executed by Exec() or Query()
	// it's therefore important to close it when it's no longer needed, to free up those resources
	defer stmt.Close()

	// Execute the prepared statement with the event's fields
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}

	// Get the ID of the last inserted row in the database and assign the generated ID back to the event object
	id, err := result.LastInsertId()
	e.ID = id

	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event

	// Next() returns a bool as long as there are rows left to iterate
	for rows.Next() {
		var event Event
		// pointers used - event populated this way
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	var event Event
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id)
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (e Event) Update() error {
	query := `
  UPDATE events
  SET name = ?, description = ?, location = ?, dateTime = ?
  WHERE id = ?
  `

	statement, error := db.DB.Prepare(query)
	if error != nil {
		return error
	}

	defer statement.Close()

	_, error = statement.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)
	return error
}

func (e Event) Delete(id int64) error {
	query := `DELETE FROM events WHERE id = ?`

	statement, error := db.DB.Prepare(query)
	if error != nil {
		return error
	}

	defer statement.Close()
	_, error = statement.Exec(e.ID)

	return error
}

func (e Event) Register(userId int64) error {
	query := "INSERT INTO registrations(event_id, user_id) VALUES (?, ?)"
	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(e.ID, userId)

	return err
}

func CheckIfAlreadyRegistered(userId int64, eventId int64) (bool, error) {
	query := "SELECT 1 FROM registrations WHERE user_id = ? AND event_id = ?"
	row := db.DB.QueryRow(query, userId, eventId)

	var exists int
	err := row.Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// no registrations found
			return false, nil
		}
		return false, fmt.Errorf("check registrations for user %d and event %d: %w", userId, eventId, err)
	}
	// registration found
	return true, nil
}

func (e Event) CancelRegistration(userId int64) error {
	query := "DELETE FROM registrations WHERE user_id = ? AND event_id = ?"

	statement, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("Preparing query statement: %w", err)
	}

	_, err = statement.Exec(userId, &e.ID)
	if err != nil {
		return fmt.Errorf("Deleting registration: %w", err)
	}

	return nil
}
