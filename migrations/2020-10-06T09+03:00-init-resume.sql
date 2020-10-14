
-- +migrate Up

set search_path to main;
create type gender_type as enum ('male', 'female');
create type career_level_type as enum ('junior', 'middle', 'senior');
create type ed_level_type as enum ('middle', 'bachelor', 'master', 'doctoral');

create table resume
(
    resume_id uuid default uuid_generate_v4() not null
        constraint resume_pkey
            primary key,
    cand_id uuid default uuid_generate_v4() not null
--             constraint resume_pkey
            references candidates(cand_id),
    title varchar(128) not null,
    about_me text not null,
    salary_min int null,
    salary_max int null,
    gender gender_type null,
    career_level career_level_type null,
    education_level ed_level_type null,
    experience_month int null,
    skills text[] null,
    date_create date not null

);

-- +migrate Down

drop table main.resume;
