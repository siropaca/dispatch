-- +goose Up
-- users は identity context のアカウント。id は UUID v7 をアプリ側で採番する。
create table users (
    id                uuid        primary key,
    auth_provider_id  text        not null unique,
    handle            text        not null unique,
    display_name      text        not null,
    created_at        timestamptz not null default now(),
    updated_at        timestamptz not null default now()
);

-- +goose Down
drop table users;
