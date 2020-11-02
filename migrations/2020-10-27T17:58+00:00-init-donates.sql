-- +migrate Up
set search_path to main;
create type donate_class_type as enum ('minimum', 'standard', 'premium', 'private');

create table donates
(
    donate_id uuid default uuid_generate_v4() not null
        constraint donate_pkey primary key,
    vac_id uuid default uuid_generate_v4() not null
            references vacancy(vac_id),
    activation_date date default current_date,
    isActive bool,
    class donate_class_type default 'private'
);
-- +migrate Down
drop table main.donates;
