package sqlite

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/omgitsotis/pocket-stats/pocket/model"
	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger("sqlite")

// SQLiteDAO is the object that connects to the sqlite database
type SQLiteDAO struct {
	db *sql.DB
}

// AddUser adds a user to the database
func (dao *SQLiteDAO) AddUser(name string) (int64, error) {
	stmt, err := dao.db.Prepare("INSERT INTO users(username) VALUES (?)")
	if err != nil {
		logger.Errorf("Error preparing database: [%s]", err.Error())
		return 0, err
	}

	logger.Debugf(
		"Running statement [INSERT INTO users(username) VALUES (%s)]",
		name,
	)

	res, err := stmt.Exec(name)
	if err != nil {
		logger.Errorf("Error executing database [%s]", err.Error())
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		logger.Errorf("Error executing database: [%s]", err.Error())
		return 0, err
	}

	logger.Infof("Created user [%d]", id)

	return id, nil
}

// AddArticle adds a new row to the article table
func (dao *SQLiteDAO) AddArticle(r model.Article) error {
	stmt, err := dao.db.Prepare("INSERT INTO articles VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		logger.Errorf("Error preparing database: [%s]", err.Error())
		return err
	}

	logger.Debugf("Inserting into articles values (%#v)", r)

	res, err := stmt.Exec(
		r.ID,
		r.DateAdded,
		r.DateRead,
		r.WordCount,
		r.Status,
		r.UserID,
		r.Tag,
	)

	if err != nil {
		logger.Errorf("Error executing database [%s]", err.Error())
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		logger.Errorf("Error executing database [%s]", err.Error())
		return err
	}

	logger.Debugf("Row(s) added: [%d]", n)
	return nil
}

// IsUser checks to see if the user is in the database
func (dao *SQLiteDAO) IsUser(id int64) (bool, error) {
	var username string
	logger.Debugf("SELECT username FROM users WHERE id=%d", id)
	err := dao.db.QueryRow("SELECT username FROM users WHERE id=?", id).Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		logger.Errorf("No user with id [%d]\n", id)
		return false, nil
	case err != nil:
		logger.Errorf("Error reading username: [%s]", err.Error())
		return false, err
	default:
		logger.Infof("Found user [%s]", username)
		return true, nil
	}
}

// GetArticles returns a list of articles between two dates
func (dao *SQLiteDAO) GetArticles(start, end int64) ([]model.Article, error) {
	query := "SELECT id, date_added, date_read, word_count, status, tag FROM articles " +
		"WHERE date_added >= ? AND date_added <= ? " +
		"OR date_read >= ? AND date_read <= ? " +
		"ORDER BY date_added"

	articles := make([]model.Article, 0)

	logger.Infof("selecting articles between %d and %d", start, end)

	rows, err := dao.db.Query(query, start, end, start, end)
	if err != nil {
		logger.Errorf("Error executing query: %s", err.Error())
		return nil, err
	}

	for rows.Next() {
		var a model.Article
		if err := rows.Scan(&a.ID, &a.DateAdded, &a.DateRead, &a.WordCount, &a.Status, &a.Tag); err != nil {
			logger.Errorf("Error reading data: %s", err.Error())
			return nil, err
		}
		articles = append(articles, a)
	}

	if err := rows.Err(); err != nil {
		logger.Errorf("Error looping results: %s", err.Error())
		return nil, err
	}

	return articles, nil
}

// GetArticle returns an article for a given id
func (dao *SQLiteDAO) GetArticle(id int64) (*model.Article, error) {
	logger.Debugf("Getting article [%d]", id)

	var r model.Article
	err := dao.db.QueryRow("SELECT * FROM articles WHERE id=?", id).Scan(
		&r.ID,
		&r.DateAdded,
		&r.DateRead,
		&r.WordCount,
		&r.Status,
		&r.UserID,
		&r.Tag,
	)

	switch {
	case err == sql.ErrNoRows:
		logger.Warningf("No article with id [%d]", id)
		return nil, nil
	case err != nil:
		logger.Errorf("Error reading article: [%s]", err.Error())
		return nil, err
	default:
		logger.Debugf("Found article [%d]", id)
		return &r, nil
	}
}

// UpdateArticle updates a article in the database
func (dao *SQLiteDAO) UpdateArticle(r *model.Article) error {
	stmt, err := dao.db.Prepare("UPDATE articles SET date_read = ?, status = ?, tag = ? WHERE id = ?")
	if err != nil {
		logger.Errorf("Error creating statment: %s", err.Error())
		return err
	}

	logger.Debugf(
		"Updating article [%d] date read [%d] status [%s] tag [%s]",
		r.ID, r.DateRead, r.Status, r.Tag,
	)

	res, err := stmt.Exec(r.DateRead, r.Status, r.Tag, r.ID)

	if err != nil {
		logger.Errorf("Error executing database [%s]", err.Error())
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		logger.Errorf("Error executing database [%s]", err.Error())
		return err
	}

	logger.Debugf("Row(s) updated: [%d]", n)
	return nil
}

// GetLastAdded gets the latest article and returns the most recent date
func (dao *SQLiteDAO) GetLastAdded() (int64, error) {
	logger.Debug("Getting last added")
	var id, dateAdded, dateRead int64
	err := dao.db.QueryRow("SELECT MAX(id), date_added, date_read FROM articles ORDER BY ID DESC").Scan(
		&id,
		&dateAdded,
		&dateRead,
	)

	switch {
	case err == sql.ErrNoRows:
		logger.Errorf("No articles found")
		return 0, nil
	case err != nil:
		logger.Errorf("Error reading article: [%s]", err.Error())
		return 0, err
	default:
		logger.Debugf("Found article [%d]", id)
		if dateAdded > dateRead {
			return dateAdded, nil
		}

		return dateRead, nil
	}
}

// CloseDB closes the connection of the database
func (dao *SQLiteDAO) CloseDB() {
	dao.db.Close()
}

// NewSQLiteDAO creates a new SQLite Client
func NewSQLiteDAO(file string) (*SQLiteDAO, error) {
	setupLogging()
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		logger.Errorf("Error opening database: [%s]", err.Error())
		return nil, err
	}

	dao := SQLiteDAO{db}
	return &dao, nil
}

func setupLogging() {
	var format = logging.MustStringFormatter(
		`%{color}[%{time:Mon 02 Jan 2006 15:04:05.000}] %{level:.5s} %{shortfile} %{color:reset} %{message}`,
	)

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	formatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(formatter)
}
