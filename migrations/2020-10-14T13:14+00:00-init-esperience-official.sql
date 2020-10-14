-- +migrate Up

set search_path to main;

create table experience_in_official_company
(
    exp_custom_id uuid default uuid_generate_v4() not null
        constraint experience_in_official_company_pkey
            primary key,
    cand_id uuid default uuid_generate_v4() not null
            references candidates(id),
    company_id uuid default uuid_generate_v4() not null
            references official_companies(comp_id),
    position text,
    begin date not null,
    finish date null,
    description text null
);

-- +migrate Down

drop table main.experience_in_official_company;
