package database

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/omgitsotis/pocket-stats/pkg/pocket"
	"github.com/pkg/errors"
)

var log *logrus.Logger

// Init sets the logger for this package
func Init(l *logrus.Logger) {
	log = l
}

// PostgresClient is the client used to connect to the Postgres database.
type PostgresClient struct {
	db *sql.DB
}

// NewPostgresDB creates a new PostgresClient from a Postgres connection
func NewPostgresDB(db *sql.DB) *PostgresClient {
	return &PostgresClient{
		db: db,
	}
}

// SaveArticles saves a list of articles to the database.
// This function is deprecated, as it is only used for the initial population of
// the database.
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

// GetArticle gets an article for a given ID
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

// UpsertArticles will loop through a given list of articles and either insert
// or update the article.
func (p *PostgresClient) UpsertArticles(articles []pocket.Article) error {
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

// updateArticle updates an article with a new read date and tag
func (p *PostgresClient) updateArticle(new Article) error {
	if new.Tag == "" && new.DateRead == 0 {
		// If the article has no tag or read date, skip it, as it has not been
		// read.
		return nil
	}

	log.Debugf(
		"UPDATE articles SET date_read = %d, tag = %s WHERE id = %d",
		new.DateRead,
		new.Tag,
		new.ID,
	)

	stmt := "UPDATE articles SET date_read = $2, tag = $3 WHERE id = $1"
	_, err := p.db.Exec(stmt, new.ID, new.DateRead, new.Tag)
	return err
}

// insertArticle inserts a new article
func (p *PostgresClient) insertArticle(a Article) error {
	log.Debugf("INSERT INTO articles %+v", a)

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

func (p *PostgresClient) GetArticlesByDate(start, end int64) ([]Article, error) {
	log.Debugf("Getting articles from [%d] to [%d]", start, end)
	stmt := `
		SELECT id, title, url, tag, word_count, date_added, date_read
		FROM articles
		WHERE date_added >= $1 and date_added <= $2
		OR date_read >= $1 and date_read <= $2;
	`
	articles := make([]Article, 0)

	rows, err := p.db.Query(stmt, start, end)
	if err != nil {
		return nil, errors.Wrap(err, "error executing query")
	}

	defer rows.Close()
	for rows.Next() {
		var article Article
		sErr := rows.Scan(
			&article.ID,
			&article.Title,
			&article.URL,
			&article.Tag,
			&article.WordCount,
			&article.DateAdded,
			&article.DateRead,
		)

		if sErr != nil {
			return nil, errors.Wrap(sErr, "Error scanning row")
		}

		articles = append(articles, article)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "error iterating results")
	}

	return articles, nil
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
