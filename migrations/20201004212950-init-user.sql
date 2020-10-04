
-- +migrate Up
create schema main;

set search_path to main;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table users
(
    user_id uuid default uuid_generate_v4() not null
        constraint users_pkey
            primary key,
    name varchar(128) not null,
    surname varchar(128)
);
-- +migrate Down
