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
