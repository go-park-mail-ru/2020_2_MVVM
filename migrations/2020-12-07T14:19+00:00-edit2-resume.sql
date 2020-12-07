-- +migrate Up

set search_path to main;

ALTER TABLE resume
    ADD COLUMN sphere int default -1;

-- +migrate Down

ALTER TABLE resume
    DROP COLUMN sphere RESTRICT;
