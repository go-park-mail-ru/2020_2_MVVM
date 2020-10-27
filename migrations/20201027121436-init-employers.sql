set search_path to main;
create type donate_class as enum ('minimum', 'standard', 'premium', 'private');

create table employers (
    employer_id uuid default uuid_generate_v4() not null
        constraint  primary key,
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