
-- +migrate Up
create schema if not exists main;

set search_path to main;

CREATE EXTENSION IF NOT EXISTS "citext";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table candidates
(
    cand_id uuid default uuid_generate_v4() not null
        constraint candidates_pkey primary key,
    email citext not null unique,
    nickname varchar(128) not null unique,
    password_hash bytea not null,
    name varchar(128) not null,
    surname varchar(128),
    phone varchar(32),
    area_search text,
    social_network text,
    avatar text
);
-- +migrate Down

drop table main.candidates;