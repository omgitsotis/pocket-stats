package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func (s *Store) BeginBatch() {
	s.batch = &pgx.Batch{}
}

func (s *Store) CommitBatch(ctx context.Context) error {
	res := s.db.SendBatch(ctx, s.batch)
	s.batch = nil
	return res.Close()
}

func (s *Store) SaveArticle(a *Article) {
	q := `
		INSERT INTO articles(
					id,
					title,
					url,
					tag,
					word_count,
					date_added,
					date_read)
		VALUES 		($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT	(id)
		DO UPDATE
		SET			title 		= $2,
					url 		= $3,
					tag 		= $4,
					word_count 	= $5,
					date_added	= $6,
					date_read	= $7;
	`

	s.batch.Queue(q, a.ID, a.Title, a.URL, a.Tag, a.WordCount, a.DateAdded, a.DateRead)
}

func (s *Store) DeleteArticle(id string) {
	q := `DELETE FROM articles WHERE id = $1`
	s.batch.Queue(q, id)
}

func (s *Store) GetArticle(ctx context.Context, id string) (*Article, error) {
	q := `
		SELECT 	id, 
				title, 
				url, 
				tag, 
				word_count, 
				date_added, 
				date_read
		FROM 	articles
		WHERE 	id = $1;
	`

	var a Article

	err := s.db.QueryRow(ctx, q, id).Scan(
		&a.ID,
		&a.Title,
		&a.URL,
		&a.Tag,
		&a.WordCount,
		&a.DateAdded,
		&a.DateRead,
	)

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *Store) GetArticlesByDate(ctx context.Context, start, end int64) ([]*Article, error) {
	q := `
		SELECT 	id, 
				title, 
				url, 
				tag, 
				word_count, 
				date_added, 
				date_read
		FROM 	articles
		WHERE 	date_added BETWEEN $1 AND $2
		OR 		date_read BETWEEN $1 AND $2;
	`

	articles := make([]*Article, 0)

	rows, err := s.db.Query(ctx, q, start, end)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var article Article
		err = rows.Scan(
			&article.ID,
			&article.Title,
			&article.URL,
			&article.Tag,
			&article.WordCount,
			&article.DateAdded,
			&article.DateRead,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		articles = append(articles, &article)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error reading results: %w", err)
	}

	return articles, nil
}

func (s *Store) GetArticlesByDateAndTag(ctx context.Context, start, end int64, tag string) ([]*Article, error) {
	q := `
		SELECT 	id, 
				title, 
				url, 
				tag, 
				word_count, 
				date_added, 
				date_read
		FROM 	articles
		WHERE 	(date_added BETWEEN $1 AND $2
		OR 		date_read BETWEEN $1 AND $2)
		AND		tag = $3

	`

	articles := make([]*Article, 0)

	rows, err := s.db.Query(ctx, q, start, end, tag)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var article Article
		err = rows.Scan(
			&article.ID,
			&article.Title,
			&article.URL,
			&article.Tag,
			&article.WordCount,
			&article.DateAdded,
			&article.DateRead,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		articles = append(articles, &article)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error reading results: %w", err)
	}

	return articles, nil
}

func (s *Store) GetLastUpdateDate(ctx context.Context) (int, error) {
	q := `SELECT date_updated FROM date_updated`
	var date int
	if err := s.db.QueryRow(ctx, q).Scan(&date); err != nil {
		return 0, err
	}

	return date, nil
}

func (s *Store) SaveLastUpdateDate(ctx context.Context, date int) error {
	q := `UPDATE date_updated SET date_updated = $1`
	_, err := s.db.Exec(ctx, q, date)
	return err
}
