create table tasks
(
    id          text not null
        primary key,
    prompt         text,
    article_id        text,
    status int default 0 not null,
    deleted_at  bigint,
    created_at  bigint,
    updated_at  bigint
);

-- alter table articles add column valid boolean default false not null;

