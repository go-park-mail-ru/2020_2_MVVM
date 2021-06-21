-- +migrate Up

set search_path to main;

create type sender_types_new AS ENUM ('employer', 'candidate');
delete from message where sender = 'technical';
alter table message
    alter column sender type sender_types_new
        using (sender::text::sender_types_new);

drop type sender_types;
alter type sender_types_new rename to sender_types;

create table tech_message
(
    message_id uuid default uuid_generate_v4() not null
        constraint tech_message_pkey primary key,
    chat_id uuid not null
        references chat(chat_id) on delete cascade,
    response_id uuid not null references response(response_id) on delete cascade,
    read_by_cand boolean default false,
    read_by_empl boolean default false,
    date_create timestamp not null
);

-- +migrate Down

set search_path to main;

create type sender_types_new AS ENUM ('employer', 'candidate', 'technical');
alter table message
    alter column sender type sender_types_new
        using (sender::text::sender_types_new);

drop type sender_types;
alter type sender_types_new rename to sender_types;

drop table tech_message;
