create extension if not exists pgcrypto;

create table if not exists teams (
	id uuid primary key default gen_random_uuid(),
	name text not null unique,
	country_code text null,
	created_at timestamp default now(),
	updated_at timestamp default now()
);

create table if not exists team_aliases (
	id uuid primary key default gen_random_uuid(),
	team_id uuid not null references teams(id),
	alias text not null unique,
	created_at timestamp default now()
);

create table if not exists historical_matches (
	id uuid primary key default gen_random_uuid(),
	match_date date not null,
	home_team_id uuid not null references teams(id),
	away_team_id uuid not null references teams(id),
	home_score int not null,
	away_score int not null,
	tournament text,
	city text,
	country text,
	neutral boolean default false,
	source text not null default 'international-results',
	created_at timestamp default now(),
	updated_at timestamp default now()
);

create unique index if not exists historical_matches_unique_idx
	on historical_matches (
		match_date,
		home_team_id,
		away_team_id,
		coalesce(tournament, '')
	);

create table if not exists historical_goalscorers (
	id uuid primary key default gen_random_uuid(),
	match_id uuid references historical_matches(id),
	match_date date not null,
	team_id uuid references teams(id),
	scorer text,
	minute int null,
	own_goal boolean default false,
	penalty boolean default false,
	source text not null default 'international-results',
	created_at timestamp default now()
);

create table if not exists fifa_rankings (
	id uuid primary key default gen_random_uuid(),
	team_id uuid not null references teams(id),
	ranking_date date not null,
	rank int not null,
	total_points numeric null,
	previous_points numeric null,
	rank_change int null,
	confederation text null,
	source text not null default 'fifa-ranking-historical',
	created_at timestamp default now()
);

create unique index if not exists fifa_rankings_unique_idx
	on fifa_rankings (team_id, ranking_date);

create table if not exists data_import_logs (
	id uuid primary key default gen_random_uuid(),
	import_type text not null,
	file_path text,
	status text not null,
	processed_count int default 0,
	inserted_count int default 0,
	skipped_count int default 0,
	error_count int default 0,
	error_message text null,
	started_at timestamp default now(),
	finished_at timestamp null
);

create table if not exists external_api_snapshots (
	id uuid primary key default gen_random_uuid(),
	provider text not null,
	endpoint text not null,
	payload_json jsonb not null,
	fetched_at timestamp default now(),
	expires_at timestamp not null
);

create index if not exists external_api_snapshots_lookup_idx
	on external_api_snapshots (provider, endpoint, expires_at);

