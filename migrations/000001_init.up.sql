create table articles
(
    id          text not null
        primary key,
    url         text,
    tool        text,
    style       text,
    keyword     text,
    author_id   text,
    author_name text,
    deleted_at  bigint,
    created_at  bigint,
    updated_at  bigint
);

