package db

import (
  "database/sql"
  /* 
    The underscore tells Go to keep this import in the file even if it appears that we are not using it. 
    The built in package will user features of this package under the hood.
  */
  _ "github.com/mattn/go-sqlite3" 
)

var DB *sql.DB

func InitDB() {
  DB, err := sql.Open("sqlite3", "api.db") // creates the file if it doesn't exist

  if err != nil {
    panic("unable to connect to the database")
  }

  DB.SetMaxOpenConns(10) // Maximum simultaneous connections.
  DB.SetMaxIdleConns(5) // How many connections to keep open if nobody is using the available connections
}
