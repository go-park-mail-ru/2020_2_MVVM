-- +migrate Up
set search_path to main;

create type users_types as enum ('employer', 'candidate');

create table users
(
    user_id uuid default uuid_generate_v4() not null
        constraint users_pkey primary key,
    user_type users_types default 'candidate',
    email citext not null unique,
    password_hash bytea not null,
    name varchar(128) not null,
    surname varchar(128),
    phone varchar(18),
    social_network text
);

create table candidates
(
    cand_id uuid default uuid_generate_v4() not null
        constraint candidates_pkey primary key,
    user_id uuid default uuid_generate_v4() not null
            references users(user_id)
);

create table employers
(
    empl_id uuid default uuid_generate_v4() not null
        constraint employers_pkey primary key,
    user_id uuid default uuid_generate_v4() not null
        references users(user_id),
    comp_id uuid references official_companies(comp_id)
);

-- +migrate Down

drop table main.employers, main.candidates, main.users;
