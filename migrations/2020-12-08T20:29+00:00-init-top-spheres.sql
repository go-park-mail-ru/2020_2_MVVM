-- +migrate Up
set search_path to main;

create table if not exists sphere
(
    sphere_idx int not null,
    sphere_cnt int default 0
);

-- +migrate StatementBegin
create or replace function update_sphere()
    returns trigger
    language plpgsql as
$BODY$
begin
    if TG_OP = 'INSERT' then
        if EXISTS(select 1 from main.sphere where sphere_idx = new.sphere) then
            update main.sphere set sphere_cnt = sphere_cnt + 1 where sphere_idx = new.sphere;
        else
            insert into main.sphere(sphere_idx, sphere_cnt) values (new.sphere, 1);
        end if;
    elseif TG_OP = 'DELETE' then
        if EXISTS(select 1 from main.sphere where sphere_idx = old.sphere) then
            update main.sphere set sphere_cnt = sphere_cnt - 1 where sphere_idx = old.sphere;
        end if;
    end if;
    return new;
end
$BODY$;
-- +migrate StatementEnd

-- +migrate StatementBegin
create trigger update_top_spheres
    after insert or delete
    on main.vacancy
    for each row
execute procedure update_sphere();
-- +migrate StatementEnd

-- +migrate Up notransaction
create index if not exists date_idx on main.vacancy (date_create);

-- +migrate Down
drop trigger if exists update_top_spheres on main.vacancy;
drop function if exists main.update_sphere();
drop table if exists main.sphere;
drop index main.date_idx;
