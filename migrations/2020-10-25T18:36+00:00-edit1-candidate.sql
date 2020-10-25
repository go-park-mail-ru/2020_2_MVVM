-- +migrate Up
alter table main.candidates drop column area_search;

-- +migrate Down
alter table main.candidates add column area_search text;
