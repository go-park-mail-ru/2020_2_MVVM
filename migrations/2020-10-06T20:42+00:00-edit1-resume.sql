-- +migrate Up

set search_path to main;

alter table resume
    drop column salary,
    drop column title,
    drop column skills,
    drop column views;

ALTER TABLE resume
    ADD salary_min int null;

alter table resume
    add salary_max int null;

create type gender_type as enum ('male', 'female');
alter table resume
    add gender gender_type null;

create type level_type as enum ('junior', 'middle', 'senior');
alter table resume
    add level level_type null;

create type education_type as enum ('master', 'bachelor');
alter table resume
    add education education_type null;

alter table resume
    add experience_month int null;

-- +migrate Down

set search_path to main;

alter table resume
    drop column salary_max,
    drop column salary_min,
    drop column gender,
    drop column level,
    drop column education,
    drop column experience_month;

drop type gender_type, level_type, education_type;

ALTER TABLE resume
    ADD salary int null,
    ADD title int null,
    ADD views int null,
    ADD skills int null;
