
-- +migrate Up
set search_path to main;


-- +migrate StatementBegin
create or replace function update_vac_cnt()
    returns trigger
    language plpgsql as
$BODY$
begin
    if TG_OP = 'INSERT' then
        if EXISTS(select 1 from main.official_companies where comp_id = new.comp_id) then
            update main.official_companies set count_vacancy = count_vacancy + 1 where comp_id = new.comp_id;
        end if;
    elseif TG_OP = 'DELETE' then
        if EXISTS(select 1 from main.official_companies where comp_id = old.comp_id) then
            update main.official_companies set count_vacancy = count_vacancy - 1 where comp_id = old.comp_id;
        end if;
    end if;
return new;
end
$BODY$;
-- +migrate StatementEnd

-- +migrate StatementBegin
create trigger update_vac_cnt_in_comp
    after insert or delete
on main.vacancy
    for each row
execute procedure update_vac_cnt();
-- +migrate StatementEnd



-- +migrate Down
drop trigger if exists update_vac_cnt_in_comp on main.vacancy;
drop function if exists main.update_vac_cnt();

