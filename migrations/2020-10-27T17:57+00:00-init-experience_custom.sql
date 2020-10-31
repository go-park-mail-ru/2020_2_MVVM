-- +migrate Up

set search_path to main;

create table experience_in_custom_company
(
    exp_custom_id uuid default uuid_generate_v4() not null
        constraint experience_in_custom_company_pkey
            primary key,
    cand_id uuid default uuid_generate_v4() not null
            references candidates(cand_id),
    resume_id uuid default uuid_generate_v4() not null
            references resume(resume_id),
    name_job text not null,
    position text null,
    duties text null,
    begin date not null,
    finish date null,
    continue_to_today boolean default false
);

-- +migrate Down

drop table main.experience_in_custom_company;
