-- +migrate Up
set search_path to main;
drop table main.response;

create table response
(
    response_id uuid default uuid_generate_v4() not null
        constraint response_pkey primary key,
    vacancy_id uuid default uuid_generate_v4() not null
        references vacancy(vac_id),
    resume_id uuid default uuid_generate_v4() not null
        references resume(resume_id),
    initial users_types not null,
    status status_response default 'sent',
    date_create date not null,
    constraint response_unique unique (vacancy_id, resume_id)

);

-- +migrate Down
drop table main.response;
