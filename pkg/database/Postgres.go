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
		stmt.Close()
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

func (p *PostgresClient) GetArticle(id int) (*Article, error) {
	var article Article
	stmt := `
		SELECT id, title, url, tag, word_count, date_added, date_read
		FROM articles
		WHERE id=$1;`

	row := p.db.QueryRow(stmt, id)
	err := row.Scan(
		&article.ID,
		&article.Title,
		&article.URL,
		&article.Tag,
		&article.WordCount,
		&article.DateAdded,
		&article.DateRead,
	)

	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (p *PostgresClient) UpdateArticles(articles []pocket.Article) error {
	for _, article := range articles {
		a := ConvertArticles(article)

		dbArticle, err := p.GetArticle(article.ItemID)
		if err != nil {
			// If a No rows error was returned, insert
			if err == sql.ErrNoRows {
				if iErr := p.insertArticle(a); iErr != nil {
					return iErr
				}
				continue
			}

			// Return any other error
			return err
		}

		// If no returned article was found, insert
		if dbArticle == nil {
			if iErr := p.insertArticle(a); iErr != nil {
				return iErr
			}
			continue
		}

		// Otherwise update
		if uErr := p.updateArticle(a); uErr != nil {
			return uErr
		}
	}

	return nil
}

func (p *PostgresClient) updateArticle(new Article) error {
	logrus.Debugf(
		"UPDATE articles SET date_read = %d, tag = %s WHERE id = %d",
		new.DateRead,
		new.Tag,
		new.ID,
	)

	stmt := "UPDATE articles SET date_read = $2, tag = $3 WHERE id = $1"
	_, err := p.db.Exec(stmt, new.ID, new.DateRead, new.Tag)
	return err
}

func (p *PostgresClient) insertArticle(a Article) error {
	logrus.Debugf("INSERT INTO articles %+v", a)

	stmt := `
		INSERT INTO articles (id, title, url, tag, word_count, date_added, date_read)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := p.db.Exec(
		stmt,
		a.ID, a.Title, a.URL,
		a.Tag, a.WordCount, a.DateAdded, a.DateRead,
	)

	return err
}

// GetLastUpdateDate returns the date the database is updated to.
func (p *PostgresClient) GetLastUpdateDate() (int, error) {
	var date int
	stmt := "SELECT date_updated from date_updated;"
	row := p.db.QueryRow(stmt)
	err := row.Scan(&date)
	return date, err
}

// SaveUpdateDate sets the updated date.
func (p *PostgresClient) SaveUpdateDate(date int64) error {
	stmt := "UPDATE date_updated SET date_updated = $1"
	_, err := p.db.Exec(stmt, date)
	return err
}
