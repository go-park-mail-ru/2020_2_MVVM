-- +migrate Up
create schema if not exists main;

set search_path to main;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table vacancy
(
    vacancy_idx         serial,
    vacancy_id          uuid default uuid_generate_v4() not null
        constraint vacancy_pkey
            primary key,
--     user_id             uuid default uuid_generate_v4() not null
--         references candidates (cand_id),
    vacancy_name        varchar(128)                    not null,
    company_name        varchar(128)                    not null,
    vacancy_description text,
    work_experience     text,
    company_address     varchar(256),
-- telephone
    skills
    text,
    salary              int
);

-- +migrate Down
drop table main.vacancy;
