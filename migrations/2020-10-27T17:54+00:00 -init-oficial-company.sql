-- +migrate Up
create schema if not exists main;
set search_path to main;

CREATE EXTENSION IF NOT EXISTS "citext";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table official_companies
(
    comp_id uuid default uuid_generate_v4() not null
        constraint comp_id_pkey
            primary key,
    name text not null,
    spheres int[] null,
    description text null,
    location text not null,
    link text null,
    count_vacancy int default 0
);

-- +migrate Down

drop table main.official_companies;
