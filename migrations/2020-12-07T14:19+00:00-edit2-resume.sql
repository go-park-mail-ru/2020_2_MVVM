-- +migrate Up

set search_path to main;

ALTER TABLE resume
    ADD COLUMN sphere int default -1;

-- +migrate Down

set search_path to main;

ALTER TABLE resume
    DROP COLUMN sphere RESTRICT;
