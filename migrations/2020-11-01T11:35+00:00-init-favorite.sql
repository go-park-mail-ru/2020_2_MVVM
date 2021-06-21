-- +migrate Up
set search_path to main;

create table favorite_for_cand
(
    favorite_id uuid default uuid_generate_v4() not null
        constraint favorite_for_cand_pkey primary key,
    cand_id uuid default uuid_generate_v4() not null
        references candidates(cand_id) ON DELETE CASCADE,
    vacancy_id uuid default uuid_generate_v4() not null
        references vacancy(vac_id) ON DELETE CASCADE,
    constraint like_unique_cand unique (cand_id, vacancy_id)
);

create table favorite_for_empl
(
    favorite_id uuid default uuid_generate_v4() not null
        constraint favorite_for_empl_pkey primary key,
    empl_id uuid default uuid_generate_v4() not null
        references employers(empl_id) ON DELETE CASCADE,
    resume_id uuid default uuid_generate_v4() not null
        references resume(resume_id) ON DELETE CASCADE,
    constraint like_unique_empl unique (empl_id, resume_id)
);

-- +migrate Down
drop table main.favorite_for_cand, main.favorite_for_empl;
