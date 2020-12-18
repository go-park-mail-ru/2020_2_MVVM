-- +migrate Up

set search_path to main;

ALTER TABLE tech_message
    ADD COLUMN response_status status_response default 'sent';

-- +migrate Down

set search_path to main;

ALTER TABLE tech_message
    DROP COLUMN response_status RESTRICT;
