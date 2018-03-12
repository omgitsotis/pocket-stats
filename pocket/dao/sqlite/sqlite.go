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

func (dao *SQLiteDAO) GetCountForDates(start, end int64) (*model.Stats, error) {
	query := "SELECT date_added, (word_count) FROM articles " +
		"WHERE date_added >= ? AND date_added <= ? " +
		"GROUP BY date_added " +
		"ORDER BY date_added DESC"

	res := make(map[int64]*model.WordCount)

	log.Printf("selecting date_added between %d and %d", start, end)

	rows, err := dao.db.Query(query, start, end)
	if err != nil {
		log.Printf("Error executing query: %s", err.Error())
		return nil, err
	}

	for rows.Next() {
		var date, count int64
		if err := rows.Scan(&date, &count); err != nil {
			log.Printf("Error reading data: %s", err.Error())
			return nil, err
		}

		res[date] = &model.WordCount{Added: count}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error looping results: %s", err.Error())
		return nil, err
	}

	query = "SELECT date_read, (word_count) FROM articles " +
		"WHERE date_read >= ? AND date_added <= ? " +
		"GROUP BY date_read " +
		"ORDER BY date_read DESC"

	log.Printf("selecting date_read between %d and %d", start, end)

	rows, err = dao.db.Query(query, start, end)
	if err != nil {
		log.Printf("Error executing query: %s", err.Error())
		return nil, err
	}

	for rows.Next() {
		var date, count int64
		if err := rows.Scan(&date, &count); err != nil {
			log.Printf("Error reading data: %s", err.Error())
			return nil, err
		}

		v, ok := res[date]
		if !ok {
			res[date] = &model.WordCount{Read: count}
		} else {
			v.Read = count
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error looping results: %s", err.Error())
		return nil, err
	}

	stats := model.Stats{Value: res}
	return &stats, nil
}

func (dao *SQLiteDAO) GetArticle(id int64) (*model.Row, error) {
	var r model.Row
	log.Printf("Getting article [%d]", id)
	err := dao.db.QueryRow("SELECT * FROM articles WHERE id=?", id).Scan(
		&r.ID,
		&r.DateAdded,
		&r.DateRead,
		&r.WordCount,
		&r.Status,
		&r.UserID,
	)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("No article with id [%d]", id)
		return nil, nil
	case err != nil:
		log.Printf("Error reading article: [%s]", err.Error())
		return nil, err
	default:
		log.Printf("Found article [%d]", id)
		return &r, nil
	}
}

func (dao *SQLiteDAO) UpdateArticle(r *model.Row) error {
	stmt, err := dao.db.Prepare("UPDATE articles SET date_read = ?, status = ? WHERE id = ?")
	if err != nil {
		log.Printf("Error creating statment: %s", err.Error())
		return err
	}

	log.Printf(
		"Updating article [%d] date read [%d] status [%s]",
		r.ID,
		r.DateRead,
		r.Status,
	)

	res, err := stmt.Exec(
		r.DateRead,
		r.Status,
		r.ID,
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

	log.Printf("Row(s) updated: [%d]", n)
	return nil
}

func (dao *SQLiteDAO) GetLastAdded() (int64, error) {
	log.Printf("Getting last added")
	var id, dateAdded, dateRead int64
	err := dao.db.QueryRow("SELECT MAX(id), date_added, date_read FROM articles ORDER BY ID DESC").Scan(
		&id,
		&dateAdded,
		&dateRead,
	)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("No articles found")
		return 0, nil
	case err != nil:
		log.Printf("Error reading article: [%s]", err.Error())
		return 0, err
	default:
		log.Printf("Found article [%d]", id)
		if dateAdded > dateRead {
			return dateAdded, nil
		}

		return dateRead, nil
	}
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
