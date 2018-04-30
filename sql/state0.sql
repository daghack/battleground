create extension if not exists pgcrypto;

create table if not exists players (
	id uuid primary key default gen_random_uuid(),
	username text not null,
	passkey text not null
);

create table if not exists unit_types (
	id text primary key,
	base_health int not null,
	base_movement int not null,
	base_cost int not null,
	attack_pattern bytea
);

create table if not exists active_games (
	id uuid primary key default gen_random_uuid(),
	board_size int not null,
	piece_count int not null,
	game_state bytea not null
);

create table if not exists active_players (
	game_id uuid not null references active_games(id) on delete cascade,
	player_id uuid not null references players(id) on delete cascade
);

create table if not exists available_unit_types (
	player_id uuid not null references players(id) on delete cascade,
	unit_type_id text not null references unit_types(id) on delete cascade
);

insert into unit_types (
	id, base_health, base_movement, base_cost, attack_pattern
) values ('footman', 1, 4, 1, null) on conflict do nothing;
