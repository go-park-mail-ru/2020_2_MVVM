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
    count_vacancy int default 0
);

-- +migrate Down

drop table main.official_companies;
