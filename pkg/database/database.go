package database

import (
	"context"
	"embed"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/omgitsotis/pocket-stats/pkg/pocket"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

//go:embed migrations/*.sql
var Migrations embed.FS

// Article is the object that will be saved in the database
type Article struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	Tag       string `json:"tag"`
	WordCount int64  `json:"word_count"`
	DateAdded int64  `json:"date_added"`
	DateRead  int64  `json:"date_read"`
}

// DBCLient is the interface to the database
type DBCLient interface {
	SaveArticles(articles []pocket.Article) error
	GetArticle(id int) (*Article, error)
	GetLastUpdateDate() (int, error)
	UpsertArticles(articles []pocket.Article) error
	SaveUpdateDate(int64) error
	GetArticlesByDate(start, end int64) ([]Article, error)
	GetArticlesByTag(start, end int64, tag string) ([]Article, error)
	DeleteArticle(int) error
}

type Store struct {
	db    *pgxpool.Pool
	batch *pgx.Batch
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
	}
}

func MustConnect(ctx context.Context, dsn string) *pgxpool.Pool {
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		logrus.WithError(err).Panic("error connecting to DB")
	}

	if err = runMigration(pool); err != nil {
		logrus.WithError(err).Panic("error running migration")
	}

	return pool
}

func runMigration(pool *pgxpool.Pool) error {
	migrationSource := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: Migrations,
		Root:       "migrations",
	}

	db := stdlib.OpenDB(*pool.Config().ConnConfig)

	_, err := migrate.Exec(db, "postgres", migrationSource, migrate.Up)
	if err != nil {
		return err
	}

	return nil
}
