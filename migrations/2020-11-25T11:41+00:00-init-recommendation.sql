-- +migrate Up
set search_path to main;

create table recommendation
(
    user_id uuid default uuid_generate_v4() not null
        references users(user_id) ON DELETE CASCADE
        constraint rec_pkey primary key,
    sphere0 numeric default 0,
    sphere1 numeric default 0,
    sphere2 numeric default 0,
    sphere3 numeric default 0,
    sphere4 numeric default 0,
    sphere5 numeric default 0,
    sphere6 numeric default 0,
    sphere7 numeric default 0,
    sphere8 numeric default 0,
    sphere9 numeric default 0,
    sphere10 numeric default 0,
    sphere11 numeric default 0,
    sphere12 numeric default 0,
    sphere13 numeric default 0,
    sphere14 numeric default 0,
    sphere15 numeric default 0,
    sphere16 numeric default 0,
    sphere17 numeric default 0,
    sphere18 numeric default 0,
    sphere19 numeric default 0,
    sphere20 numeric default 0,
    sphere21 numeric default 0,
    sphere22 numeric default 0,
    sphere23 numeric default 0,
    sphere24 numeric default 0,
    sphere25 numeric default 0,
    sphere26 numeric default 0,
    sphere27 numeric default 0,
    sphere28 numeric default 0,
    sphere29 numeric default 0
);

-- +migrate Down
drop table main.recommendation;