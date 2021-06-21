-- +migrate Up

set search_path to main;

ALTER TABLE chat
DROP COLUMN response_id;

-- +migrate Down

set search_path to main;

DELETE FROM chat;

ALTER TABLE chat
ADD response_id uuid not null references response(response_id) on delete cascade;
