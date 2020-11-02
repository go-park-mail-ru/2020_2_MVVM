-- +migrate Up
set search_path to main;
    create type status_response as enum ('sent', 'accepted', 'refusal');

create table response
(
    response_id uuid default uuid_generate_v4() not null
        constraint response_pkey primary key,
    vacancy_id uuid default uuid_generate_v4() not null
        references vacancy(vac_id),
    resume_id uuid default uuid_generate_v4() not null
        references resume(resume_id),
    initial users_types not null,
    isApply status_response default 'sent',
    date_create date not null
);
-- +migrate Down
drop table main.response;
