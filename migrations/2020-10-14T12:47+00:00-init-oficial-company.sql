-- +migrate Up

set search_path to main;

create table official_companies
(
    comp_id uuid default uuid_generate_v4() not null
        constraint comp_id_pkey
            primary key,
    name text not null,
    sphere text[] null,
    location text not null,
    link text null,
    avatar text null,
    phone text null
);

-- +migrate Down

drop table main.official_companies;
