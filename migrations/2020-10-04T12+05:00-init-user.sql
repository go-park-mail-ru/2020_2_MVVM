
-- +migrate Up
create schema if not exists main;

set search_path to main;

CREATE EXTENSION IF NOT EXISTS "citext";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table users
(
    user_id uuid default uuid_generate_v4() not null
        constraint users_pkey primary key,
    email citext not null unique,
    avatar_path varchar(128),
    nickname varchar(128) not null unique,
    password_hash bytea not null,
    name varchar(128) not null,
    surname varchar(128)
);
-- +migrate Down

drop table main.users;