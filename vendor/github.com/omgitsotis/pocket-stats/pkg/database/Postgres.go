package database

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/omgitsotis/pocket-stats/pkg/pocket"
	"github.com/pkg/errors"
)

var log *logrus.Logger

func Init(l *logrus.Logger) {
	log = l
}

type PostgresClient struct {
	db *sql.DB
}

func NewPostgresDB(db *sql.DB) *PostgresClient {
	return &PostgresClient{
		db: db,
	}
}

func (p *PostgresClient) SaveArticles(articles []pocket.Article) error {
	entries := make([]Article, 0)

	// Convert the articles to database objects
	for _, a := range articles {
		entries = append(entries, ConvertArticles(a))
	}

	// Start the transaction
	txn, err := p.db.Begin()
	if err != nil {
		return errors.WithMessage(err, "error starting transaction")
	}

	// Create the Postgres Copy function
	stmt, err := txn.Prepare(
		pq.CopyIn(
			"articles",
			"id", "title", "url", "tag", "word_count", "date_added", "date_read",
		),
	)

	if err != nil {
		return errors.WithMessage(err, "error creating Copy statement")
	}

	// Execute the queries
	for _, e := range entries {
		_, err := stmt.Exec(
			e.ID, e.Title, e.URL, e.Tag,
			e.WordCount, e.DateAdded, e.DateRead,
		)

		if err != nil {
			return errors.WithMessage(err, "error executing copy statement")
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return errors.WithMessage(err, "error executing statement")
	}

	err = stmt.Close()
	if err != nil {
		return errors.WithMessage(err, "error closing copy statement")
	}

	err = txn.Commit()
	if err != nil {
		return errors.WithMessage(err, "error committing transaction")
	}

	return nil
}
