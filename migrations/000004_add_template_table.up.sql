create table templates
(
    id          text not null
        primary key,
    prompt         text,
    deleted_at  bigint,
    created_at  bigint,
    updated_at  bigint
);

-- alter table articles add column valid boolean default false not null;

