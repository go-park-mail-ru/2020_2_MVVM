-- +migrate Up
set search_path to main;
create type employment_type as enum ('full-time', 'part-time', 'remotely');

create table vacancy (
    vac_id  uuid default uuid_generate_v4() not null
        constraint vacancy_pkey primary key,
    empl_id uuid default uuid_generate_v4() not null
        references employers(empl_id),
    comp_id uuid references official_companies(comp_id),
    title varchar(128) not null,
    salary_min int null,
    salary_max int null,
    description text not null,
    requirements text null,
    duties text null,
    skills text null,
    spheres text null,
    employment employment_type default 'full-time',
    week_work_hours int null,
    location varchar(512) null,
    career_level career_level_type null,
    education_level ed_level_type null,
    experience_month int null,
    date_create date default current_date
);

-- +migrate Down
drop table main.vacancy;
