create extension if not exists pgcrypto;

create table if not exists team_metrics (
	id uuid primary key default gen_random_uuid(),
	team_id uuid not null references teams(id),
	metric_date date not null,
	elo_score numeric null,
	attack_score numeric null,
	defense_score numeric null,
	recent_form_score numeric null,
	world_cup_history_score numeric null,
	knockout_score numeric null,
	group_stage_score numeric null,
	avg_goals_scored numeric null,
	avg_goals_conceded numeric null,
	win_rate numeric null,
	draw_rate numeric null,
	loss_rate numeric null,
	matches_played int default 0,
	source text not null default 'metrics-engine-v1',
	created_at timestamp default now(),
	updated_at timestamp default now()
);

create unique index if not exists team_metrics_team_date_idx
	on team_metrics (team_id, metric_date);

create table if not exists team_metric_snapshots (
	id uuid primary key default gen_random_uuid(),
	team_id uuid not null references teams(id),
	snapshot_type text not null,
	payload_json jsonb not null,
	calculated_at timestamp default now()
);

create index if not exists team_metric_snapshots_team_type_calculated_idx
	on team_metric_snapshots (team_id, snapshot_type, calculated_at desc);

create table if not exists match_features (
	id uuid primary key default gen_random_uuid(),
	match_id uuid null,
	match_date date not null,
	home_team_id uuid not null references teams(id),
	away_team_id uuid not null references teams(id),
	tournament text null,
	stage text null,
	home_elo_score numeric null,
	away_elo_score numeric null,
	elo_diff numeric null,
	home_attack_score numeric null,
	away_attack_score numeric null,
	home_defense_score numeric null,
	away_defense_score numeric null,
	home_recent_form_score numeric null,
	away_recent_form_score numeric null,
	home_fifa_rank int null,
	away_fifa_rank int null,
	fifa_rank_diff int null,
	home_avg_goals_scored numeric null,
	away_avg_goals_scored numeric null,
	home_avg_goals_conceded numeric null,
	away_avg_goals_conceded numeric null,
	home_world_cup_history_score numeric null,
	away_world_cup_history_score numeric null,
	neutral boolean default false,
	created_at timestamp default now(),
	updated_at timestamp default now()
);

create unique index if not exists match_features_match_unique_idx
	on match_features (
		match_date,
		home_team_id,
		away_team_id,
		tournament
	);

create index if not exists match_features_match_date_idx
	on match_features (match_date);
