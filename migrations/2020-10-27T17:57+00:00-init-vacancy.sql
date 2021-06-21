-- +migrate Up
set search_path to main;
create type employment_type as enum ('full-time', 'part-time', 'remotely');

create table vacancy (
    vac_id  uuid default uuid_generate_v4() not null
        constraint vacancy_pkey primary key,
    empl_id uuid default uuid_generate_v4() not null
        references employers(empl_id) ON DELETE CASCADE,
    comp_id uuid references official_companies(comp_id) ON DELETE CASCADE,
    title varchar(128) not null,
    salary_min int default 0,
    salary_max int default 0,
    description text not null,
    requirements text null,
    duties text null,
    skills text null,
    sphere int default -1,
    gender gender_type default null,
    employment employment_type default 'full-time',
    area_search varchar(128) null,
    location varchar(512) null,
    career_level career_level_type null,
    education_level ed_level_type null,
    experience_month int null,
    empl_email  citext not null,
    empl_phone  varchar(18) not null,
    path_to_avatar varchar(256),
    date_create timestamp default current_timestamp
);

-- +migrate Down
drop table main.vacancy;
