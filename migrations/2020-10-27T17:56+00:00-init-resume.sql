-- +migrate Up

set search_path to main;
create type gender_type as enum ('male', 'female');
create type career_level_type as enum ('junior', 'middle', 'senior');
create type ed_level_type as enum ('middle', 'specialized_secondary', 'incomplete_higher',
                                    'higher', 'bachelor', 'master', 'phD', 'doctoral');

create table resume
(
    resume_id uuid default uuid_generate_v4() not null
        constraint resume_pkey
            primary key,
    cand_id uuid default uuid_generate_v4() not null
            references candidates(cand_id) ON DELETE CASCADE,
    title varchar(128) not null,
    description text not null,
    salary_min int null,
    salary_max int null,
    gender gender_type null,
    career_level career_level_type null,
    education_level ed_level_type null,
    experience_month int null,
    skills text null,
    place text null,
    area_search text,
    path_to_avatar varchar(256),
    date_create timestamp default current_timestamp
);

-- +migrate Down

drop table main.resume;
