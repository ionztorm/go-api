package models

import (
	"github.com/LeonLonsdale/go-web-api/db"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (user User) Save() error {
	query := `INSERT INTO users(email, password) VALUES (?, ?)`

	statement, error := db.DB.Prepare(query)

	if error != nil {
		return error
	}

	defer statement.Close()

	result, error := statement.Exec(user.Email, user.Password)

	if error != nil {
		return error
	}

	userId, error := result.LastInsertId()

	if error != nil {
		return error
	}

	user.ID = userId

	return nil
}
