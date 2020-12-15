-- +migrate Up

set search_path to main;
create type sender_types as enum ('employer', 'candidate', 'technical');

create table chat
(
    chat_id uuid default uuid_generate_v4() not null
        constraint chat_pkey primary key,
    user_id_cand uuid default uuid_generate_v4() not null
        references users(user_id) ON DELETE CASCADE,
    user_id_empl uuid default uuid_generate_v4() not null
        references users(user_id) ON DELETE CASCADE,
    response_id uuid default uuid_generate_v4() not null
        references response(response_id) ON DELETE CASCADE,
    constraint chat_unique unique (user_id_cand, user_id_empl)
);

create table message
(
    message_id uuid default uuid_generate_v4() not null
        constraint message_pkey primary key,
    chat_id uuid default uuid_generate_v4() not null
        references chat(chat_id) ON DELETE CASCADE,
    sender sender_types not null,
    message text,
    is_read boolean default false,
    date_create timestamp not null
);

-- +migrate Down

drop table main.message;
drop type sender_types;
drop table main.chat;
