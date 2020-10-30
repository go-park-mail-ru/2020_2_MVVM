-- +migrate Up

set search_path to main;
create type ed_level_item_type as enum ('school', 'courses', 'middle', 'bachelor', 'master', 'doctoral');

create table education
(
    ed_id uuid default uuid_generate_v4() not null
        constraint education_pkey
            primary key,
    cand_id uuid default uuid_generate_v4() not null
            references users(user_id),
    resume_id uuid default uuid_generate_v4() not null
            references resume(resume_id),
    level ed_level_item_type null,
    begin date null,
    finish date not null,
    university text not null,
    department text null,
    description text null
);

-- +migrate Down

drop table main.education;
