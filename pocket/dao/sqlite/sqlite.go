package sqlite

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/omgitsotis/pocket-stats/pocket/model"
)

var logger = log.New(os.Stdout, "SQLite:", log.Ldate|log.Ltime|log.Lshortfile)

type SQLiteDAO struct {
	db *sql.DB
}

func (dao *SQLiteDAO) AddUser(name string) (int64, error) {
	stmt, err := dao.db.Prepare("INSERT INTO users(username) VALUES (?)")
	if err != nil {
		log.Printf("Error preparing database: [%s]", err.Error())
		return 0, err
	}

	logger.Printf(
		"Running statement [INSERT INTO users(username) VALUES (%s)]",
		name,
	)

	res, err := stmt.Exec(name)
	if err != nil {
		log.Printf("Error executing database [%s]", err.Error())
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error executing database: [%s]", err.Error())
		return 0, err
	}

	log.Printf("Created user [%d]", id)

	return id, nil
}

func (dao *SQLiteDAO) AddArticle(r model.Row) error {
	stmt, err := dao.db.Prepare("INSERT INTO articles VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing database: [%s]", err.Error())
		return err
	}

	logger.Printf("Inserting into articles values (%#v)", r)

	res, err := stmt.Exec(
		r.ID,
		r.DateAdded,
		r.DateRead,
		r.WordCount,
		r.Status,
		r.UserID,
	)

	if err != nil {
		log.Printf("Error executing database [%s]", err.Error())
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error executing database [%s]", err.Error())
		return err
	}

	log.Printf("Row(s) added: [%d]", n)
	return nil
}

func (dao *SQLiteDAO) IsUser(id int64) (bool, error) {
	var username string
	log.Printf("SELECT username FROM users WHERE id=%d", id)
	err := dao.db.QueryRow("SELECT username FROM users WHERE id=?", id).Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with id [%d]\n", id)
		return false, nil
	case err != nil:
		log.Printf("Error reading username: [%s]", err.Error())
		return false, err
	default:
		log.Printf("Found user [%s]", username)
		return true, nil
	}
}

func (dao *SQLiteDAO) GetStatsForDate(int64) error {
	//TODO This
	return nil
}

func (dao *SQLiteDAO) CloseDB() {
	dao.db.Close()
}

func NewSQLiteDAO(file string) (*SQLiteDAO, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		logger.Printf("Error opening database: [%s]", err.Error())
		return nil, err
	}

	dao := SQLiteDAO{db}
	return &dao, nil
}
