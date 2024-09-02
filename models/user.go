package models

import (
	"errors"

	"github.com/LeonLonsdale/go-web-api/db"
	"github.com/LeonLonsdale/go-web-api/utils"
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

	hashedPassword, error := utils.HashPassword(user.Password)
	if error != nil {
		return error
	}

	result, error := statement.Exec(user.Email, hashedPassword)
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

func (user User) ValidateCredentials() error {
	// check if the user exists
	query := `SELECT password FROM users WHERE email = ?`

	row := db.DB.QueryRow(query, user.Email)

	var hashedPassword string

	error := row.Scan(&hashedPassword)
	if error != nil {
		return errors.New("email or password incorrect")
	}

	// hash the Passwordas

	passwordIsValid := utils.CompareHash(user.Password, hashedPassword)
	if !passwordIsValid {
		return errors.New("email or password incorrect")
	}
	return nil
}
