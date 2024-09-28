create table users (
    id uuid not null primary key,
    nickname varchar(25) not null,
    password_hash varchar(255) not null,
    created_at timestamp default CURRENT_TIMESTAMP
);

alter table users owner to postgres;