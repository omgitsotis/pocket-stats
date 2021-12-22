-- +migrate Up
CREATE TABLE IF NOT EXISTS articles (
    id          TEXT UNIQUE,
    title       TEXT,
    url         TEXT,
    tag         TEXT,
    word_count  INTEGER,
    date_added  INTEGER,
    date_read   INTEGER
);

CREATE TABLE IF NOT EXISTS date_updated (
    date_updated INTEGER
);

CREATE INDEX IF NOT EXISTS article_id ON articles(id);
CREATE INDEX IF NOT EXISTS read_time ON articles(date_added, date_read);

ALTER TABLE articles ADD PRIMARY KEY (id);

-- +migrate Down
DROP INDEX IF EXISTS article_id;
DROP INDEX IF EXISTS read_time;

DROP TABLE IF EXISTS date_updated;
DROP TABLE IF EXISTS articles;