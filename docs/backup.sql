create table articles(
	id text unique,
	title text,
	url text,
	tag text,
	word_count integer,
	date_added integer,
	date_read integer
);

create table date_updated(
	date_updated integer
);

create index article_id on articles(id);
create index read_time on articles(date_added, date_read);

select id, title, url, tag, word_count, date_added, date_read from articles where 'id' =

SELECT id, title, url, tag, word_count, date_added, date_read
FROM articles
WHERE date_added >= 1554422400 and date_added <= 1554649444
or date_read >= 1554422400 and date_read <= 1554649444;
