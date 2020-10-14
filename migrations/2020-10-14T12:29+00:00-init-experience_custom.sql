-- +migrate Up

set search_path to main;

create table custom_company
(
    company_id uuid default uuid_generate_v4() not null
        constraint company_id_pkey
            primary key,
    name text not null,
    location text null,
    sphere text null
);

create table experience_in_custom_company
(
    exp_custom_id uuid default uuid_generate_v4() not null
        constraint experience_in_custom_company_pkey
            primary key,
    cand_id uuid default uuid_generate_v4() not null
            references candidates(id),
    company_id uuid default uuid_generate_v4() not null
            references custom_company(company_id),
    position text,
    begin date not null,
    finish date null,
    description text null
);

-- +migrate Down

drop table main.experience_in_custom_company;
drop table main.custom_company;
