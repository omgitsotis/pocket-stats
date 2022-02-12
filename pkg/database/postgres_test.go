package database_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/omgitsotis/pocket-stats/pkg/database"
	migrate "github.com/rubenv/sql-migrate"
)

const postgresURL = "postgresql://otis:password@localhost:5432/test_database?sslmode=disable"

func dropMigrations(pool *pgxpool.Pool) error {
	migrationSource := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: database.Migrations,
		Root:       "migrations",
	}

	db := stdlib.OpenDB(*pool.Config().ConnConfig)

	_, err := migrate.Exec(db, "postgres", migrationSource, migrate.Down)
	if err != nil {
		return err
	}

	return nil
}

func TestBatch(t *testing.T) {
	ctx := context.Background()

	db := database.MustConnect(ctx, postgresURL)
	defer db.Close()

	s := database.NewStore(db)
	defer dropMigrations(db)

	s.BeginBatch()
	// Insert new article
	s.SaveArticle(&database.Article{
		ID:        "1",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "tag",
		WordCount: 123,
		DateAdded: 1640189477,
		DateRead:  0,
	})

	// Update article
	s.SaveArticle(&database.Article{
		ID:        "1",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "tag",
		WordCount: 123,
		DateAdded: 1640189477,
		DateRead:  1640189522,
	})

	// Insert another article
	s.SaveArticle(&database.Article{
		ID:        "2",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "tag",
		WordCount: 123,
		DateAdded: 1640189477,
		DateRead:  0,
	})

	// Delete the article
	s.DeleteArticle("2")

	if err := s.CommitBatch(ctx); err != nil {
		t.Fatal(err)
	}

	want := &database.Article{
		ID:        "1",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "tag",
		WordCount: 123,
		DateAdded: 1640189477,
		DateRead:  1640189522,
	}

	got, err := s.GetArticle(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetArticle() mismatch: (-want +got):\n%s", diff)
	}
}

func TestGetByDateRange(t *testing.T) {
	ctx := context.Background()

	db := database.MustConnect(ctx, postgresURL)
	defer db.Close()

	s := database.NewStore(db)
	defer dropMigrations(db)

	var (
		startDate int64 = 1640131200
		endDate   int64 = 1640217600
	)

	s.BeginBatch()

	// Article read in our date range
	s.SaveArticle(&database.Article{
		ID:        "1",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "tag",
		WordCount: 123,
		DateAdded: 1640044800,
		DateRead:  1640174400,
	})

	// Article added in our date range
	s.SaveArticle(&database.Article{
		ID:        "2",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "tag",
		WordCount: 123,
		DateAdded: 1640174400,
		DateRead:  0,
	})

	// Article read in our date range (differnt tag)
	s.SaveArticle(&database.Article{
		ID:        "3",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "tag-2",
		WordCount: 123,
		DateAdded: 1640044800,
		DateRead:  1640174400,
	})

	// Article added in our date range
	s.SaveArticle(&database.Article{
		ID:        "4",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "tag-2",
		WordCount: 123,
		DateAdded: 1640174400,
		DateRead:  0,
	})

	// Article added out of range (before start)
	s.SaveArticle(&database.Article{
		ID:        "5",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "tag",
		WordCount: 123,
		DateAdded: 1640088000,
		DateRead:  0,
	})

	// Article read out of range (before start)
	s.SaveArticle(&database.Article{
		ID:        "6",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "tag",
		WordCount: 123,
		DateAdded: 1640044800,
		DateRead:  1640088000,
	})

	// Article added out of range (after end)
	s.SaveArticle(&database.Article{
		ID:        "7",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "tag",
		WordCount: 123,
		DateAdded: 1640304000,
		DateRead:  0,
	})

	// Article read out of range (after end)
	s.SaveArticle(&database.Article{
		ID:        "8",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "tag",
		WordCount: 123,
		DateAdded: 1640304000,
		DateRead:  1640347200,
	})

	if err := s.CommitBatch(ctx); err != nil {
		t.Fatal(err)
	}

	want := []*database.Article{
		{
			ID:        "1",
			URL:       "test-url",
			Title:     "test-title",
			Tag:       "tag",
			WordCount: 123,
			DateAdded: 1640044800,
			DateRead:  1640174400,
		},
		{
			ID:        "2",
			URL:       "test-url",
			Title:     "test-title",
			Tag:       "tag",
			WordCount: 123,
			DateAdded: 1640174400,
			DateRead:  0,
		}, {
			ID:        "3",
			URL:       "test-url",
			Title:     "test-title",
			Tag:       "tag-2",
			WordCount: 123,
			DateAdded: 1640044800,
			DateRead:  1640174400,
		},
		{
			ID:        "4",
			URL:       "test-url",
			Title:     "test-title",
			Tag:       "tag-2",
			WordCount: 123,
			DateAdded: 1640174400,
			DateRead:  0,
		},
	}

	got, err := s.GetArticlesByDate(ctx, startDate, endDate)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetArticleByDate() mismatch: (-want +got):\n%s", diff)
	}

	want = []*database.Article{
		{
			ID:        "1",
			URL:       "test-url",
			Title:     "test-title",
			Tag:       "tag",
			WordCount: 123,
			DateAdded: 1640044800,
			DateRead:  1640174400,
		},
		{
			ID:        "2",
			URL:       "test-url",
			Title:     "test-title",
			Tag:       "tag",
			WordCount: 123,
			DateAdded: 1640174400,
			DateRead:  0,
		},
	}

	got, err = s.GetArticlesByDateAndTag(ctx, startDate, endDate, "tag")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetArticleByDateAndTag() mismatch: (-want +got):\n%s", diff)
	}
}

func TestGetArticles(t *testing.T) {
	ctx := context.Background()

	db := database.MustConnect(ctx, postgresURL)
	defer db.Close()

	s := database.NewStore(db)
	defer dropMigrations(db)

	s.BeginBatch()

	// Article read in our date range
	s.SaveArticle(&database.Article{
		ID:        "1",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "tag",
		WordCount: 123,
		DateAdded: 1640044800,
		DateRead:  1640174400,
	})

	// Article added in our date range
	s.SaveArticle(&database.Article{
		ID:        "2",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "tag",
		WordCount: 123,
		DateAdded: 1640174400,
		DateRead:  0,
	})

	// Article read in our date range (differnt tag)
	s.SaveArticle(&database.Article{
		ID:        "3",
		URL:       "test-url",
		Title:     "test-title",
		Tag:       "",
		WordCount: 123,
		DateAdded: 1640044800,
		DateRead:  1640174400,
	})

	if err := s.CommitBatch(ctx); err != nil {
		t.Fatal(err)
	}

	want := []*database.Article{
		{
			ID:        "1",
			URL:       "test-url",
			Title:     "test-title",
			Tag:       "tag",
			WordCount: 123,
			DateAdded: 1640044800,
			DateRead:  1640174400,
		},
		{
			ID:        "3",
			URL:       "test-url",
			Title:     "test-title",
			Tag:       "",
			WordCount: 123,
			DateAdded: 1640044800,
			DateRead:  1640174400,
		},
		{
			ID:        "2",
			URL:       "test-url",
			Title:     "test-title",
			Tag:       "tag",
			WordCount: 123,
			DateAdded: 1640174400,
			DateRead:  0,
		},
	}

	got, err := s.GetArticles(ctx, 10, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetArticles() without empty tags mismatch: (-want +got):\n%s", diff)
	}

	want = []*database.Article{
		{
			ID:        "3",
			URL:       "test-url",
			Title:     "test-title",
			Tag:       "",
			WordCount: 123,
			DateAdded: 1640044800,
			DateRead:  1640174400,
		},
	}

	got, err = s.GetArticles(ctx, 10, 0, true)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetArticles() with empty tags  mismatch: (-want +got):\n%s", diff)
	}
}
