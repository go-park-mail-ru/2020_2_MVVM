
-- +migrate Up
-- create schema main;

set search_path to main;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table resume
(
    resume_id uuid default uuid_generate_v4() not null
        constraint resume_pkey
            primary key,
    user_id uuid default uuid_generate_v4() not null
--             constraint resume_pkey
            references users(user_id),
    title varchar(128) not null,
    salary int,
    description text,
    skills text,
    views int
);
-- +migrate Down
