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
