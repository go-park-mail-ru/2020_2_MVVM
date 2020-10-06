-- +migrate Up
create schema main;

set search_path to main;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table users
(
    user_id     uuid default uuid_generate_v4() not null
        constraint users_pkey
            primary key,
    name        varchar(128)                    not null,
    surname     varchar(128),
);

create table resume
(
    resume_id   uuid default uuid_generate_v4() not null
        constraint resume_pkey
            primary key,
    user_id     uuid default uuid_generate_v4() not null
--             constraint resume_pkey
        references users (user_id),
    title       varchar(128)                    not null,
    salary      int,
    description text,
    skills      text,
    views       int
);

create table vacancy
(
    vacancy_idx         serial,
    vacancy_id          uuid default uuid_generate_v4() not null
        constraint vacancy_pkey
            primary key,
    user_id             uuid default uuid_generate_v4() not null
        references users (user_id),
    vacancy_name        varchar(128)                    not null,
    company_name        varchar(128)                    not null,
    vacancy_description text,
    work_experience     text,
    company_address     varchar(256), 15
    :
    23
    ПКЛПО
    бывший
    СС
    Фотография
    1
    новое
    сообщение
    1
    14
    :
    36

    Лиза
    Яногьян
    Ты
    чудо
-- telephone
    skills
    text,
    salary              int
);
-- +migrate Down
