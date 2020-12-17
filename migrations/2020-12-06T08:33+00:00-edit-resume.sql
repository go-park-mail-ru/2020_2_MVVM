-- +migrate Up

set search_path to main;
ALTER TYPE gender_type ADD value '';

ALTER TABLE resume
    ADD COLUMN cand_name varchar(128),
    ADD COLUMN cand_surname varchar(128),
    ADD COLUMN cand_email citext;

-- +migrate Down

set search_path to main;
ALTER TABLE resume
    DROP COLUMN cand_name RESTRICT,
    DROP COLUMN cand_surname RESTRICT,
    DROP COLUMN cand_email RESTRICT;
