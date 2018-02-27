package pocket

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func AddUser(name string) (int64, error) {
	db, err := sql.Open("sqlite3", "database/pocket.db")
	if err != nil {
		fmt.Printf("Error opening database: %s\n", err.Error())
		return 0, err
	}

	stmt, err := db.Prepare("INSERT INTO users(username) VALUES (?)")
	if err != nil {
		fmt.Printf("Error preparing database: %s\n", err.Error())
		return 0, err
	}

	res, err := stmt.Exec(name)
	if err != nil {
		fmt.Printf("Error executing database %s\n", err.Error())
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Error executing database %s\n", err.Error())
		return 0, err
	}

	fmt.Printf("Created user %d\n", id)
	db.Close()

	return id, nil
}

func AddRow(r Row) error {
	db, err := sql.Open("sqlite3", "database/pocket.db")
	if err != nil {
		fmt.Printf("Error opening database: %s\n", err.Error())
		return err
	}

	stmt, err := db.Prepare("INSERT INTO articles VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Printf("Error preparing database: %s\n", err.Error())
		return err
	}

	res, err := stmt.Exec(
		r.ID,
		r.DateAdded,
		r.DateRead,
		r.WordCount,
		r.Status,
		r.UserID,
	)

	if err != nil {
		fmt.Printf("Error executing database %s\n", err.Error())
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Error executing database %s\n", err.Error())
		return err
	}
	fmt.Printf("Row(s) affected %d\n", n)
	db.Close()

	return nil
}

func HasID(id int64) (bool, error) {
	db, err := sql.Open("sqlite3", "database/pocket.db")
	if err != nil {
		fmt.Printf("Error opening database: %s\n", err.Error())
		return false, err
	}

	var username string

	err = db.QueryRow("SELECT username FROM users WHERE id=?", id).Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("No user with id %s\n", id)
		return false, nil
	case err != nil:
		fmt.Printf("Error reading username: %s", err.Error())
		return false, err
	default:
		return true, nil
	}
}
