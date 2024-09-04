package models

import (
	"errors"
	"fmt"

	"github.com/LeonLonsdale/go-web-api/db"
	"github.com/LeonLonsdale/go-web-api/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
	Salt     string
}

func (user *User) Save() error {
	query := `INSERT INTO users(email, password, salt) VALUES (?, ?, ?)`

	statement, error := db.DB.Prepare(query)
	if error != nil {
		return error
	}
	defer statement.Close()

	hashedPassword, salt, error := utils.HashPasswordArgon2id(user.Password)

	if error != nil {
		return error
	}

	result, error := statement.Exec(user.Email, hashedPassword, salt)
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

func (user *User) ValidateCredentials() error {
	// check if the user exists
	query := `SELECT id, password, salt FROM users WHERE email = ?`

	row := db.DB.QueryRow(query, user.Email)

	var hashedPassword, salt string

	fmt.Print("pre scan")

	error := row.Scan(&user.ID, &hashedPassword, &salt)
	if error != nil {
		return errors.New("email or password incorrect")
	}
	fmt.Println("got hashed and salt")
	// hash the Passwordas

	passwordIsValid := utils.CompareHashArgon2id(user.Password, salt, hashedPassword)
	if !passwordIsValid {
		return errors.New("email or password incorrect")
	}
	return nil
}
