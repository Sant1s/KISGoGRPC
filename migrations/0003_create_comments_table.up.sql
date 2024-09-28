create table if not exists comments
(
    id          bigint not null
        primary key,
    author_id   uuid   not null references users,
    parent_id   bigint default 0,
    post_id     bigint not null
        references posts,
    data        text   not null,
    created_at  timestamp default CURRENT_TIMESTAMP,
    likes_count bigint    default 0
);

alter table comments
    owner to postgres;

alter table comments
alter column parent_id set default 0;