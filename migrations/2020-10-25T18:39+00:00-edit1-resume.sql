-- +migrate Up
alter table main.resume add column area_search text;

-- +migrate Down
alter table main.resume drop column area_search;
