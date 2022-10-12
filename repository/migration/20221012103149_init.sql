-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- users is the exception to the rule of singular table names, since it's 
-- a reserved name in postgres.
CREATE TABLE IF NOT EXISTS users (
	id uuid DEFAULT uuid_generate_v4() not null, 
	email varchar not null,
	password_hash varchar not null,
	primary key (id),
	UNIQUE (email)
);

-- +goose Down
DROP TABLE IF EXISTS users;

