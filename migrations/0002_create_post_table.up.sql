create table posts (
    id bigint not null primary key,
    author_id uuid not null references users,
    data text not null,
    created_at timestamp default CURRENT_TIMESTAMP,
    comments_count bigint default 0,
    likes_count bigint default 0
);

alter table posts owner to postgres;