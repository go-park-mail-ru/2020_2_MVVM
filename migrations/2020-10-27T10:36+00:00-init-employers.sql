-- +migrate Up
set search_path to main;

create table employers
(
    empl_id uuid default uuid_generate_v4() not null
        constraint employers_pkey primary key,
    comp_id uuid default uuid_generate_v4() not null
            references official_companies(comp_id),
    name varchar(128) not null,
    surname varchar(128),
    email citext not null unique,
    nickname varchar(128) not null unique,
    password_hash bytea not null,
    phone varchar(32)
);

-- +migrate Down
drop table main.employers;