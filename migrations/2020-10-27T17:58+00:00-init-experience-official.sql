-- +migrate Up

set search_path to main;

create table experience_in_official_company
(
    exp_custom_id uuid default uuid_generate_v4() not null
        constraint experience_in_official_company_pkey
            primary key,
    cand_id uuid default uuid_generate_v4() not null
            references candidates(cand_id) ON DELETE CASCADE,
    resume_id uuid default uuid_generate_v4() not null
            references resume(resume_id) ON DELETE CASCADE,
    company_id uuid default uuid_generate_v4() not null
            references official_companies(comp_id) ON DELETE CASCADE,
    position text,
    begin date not null,
    finish date null,
    duties text null,
    continue_to_today boolean
);

-- +migrate Down

drop table main.experience_in_official_company;
