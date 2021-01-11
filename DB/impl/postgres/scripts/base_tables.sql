begin;

create schema if not exists demo;
create extension if not exists "pgcrypto" with schema demo;

set schema 'demo';

create table user_create_history
(
    created_at bigint  not null default extract(epoch from current_timestamp) * 1000,
    created_by varchar not null
);

create table user_update_history
(
    updated_at bigint  not null default extract(epoch from current_timestamp) * 1000,
    updated_by varchar not null
);

create table org
(
    id      uuid primary key default gen_random_uuid(),
    label   varchar          default '',
    deleted bool             default false
) inherits (user_create_history, user_update_history);

create table repo
(
    id      uuid primary key default gen_random_uuid(),
    label   varchar          default '',
    path    varchar not null,
    deleted bool             default false
) inherits (user_create_history, user_update_history);

create table project
(
    id      uuid primary key default gen_random_uuid(),
    label   varchar          default '',
    path    varchar not null,
    deleted bool             default false
) inherits (user_create_history, user_update_history);

create table folder
(
    id      uuid primary key default gen_random_uuid(),
    label   varchar          default '',
    path    varchar not null,
    deleted bool             default false
) inherits (user_create_history, user_update_history);

create table respath
(
    path_url     varchar primary key,
    resource_id  uuid    not null,
    type         varchar not null, -- this should point to the resource table
    previous_url varchar not null, -- self-reference
    hierarchy    jsonb           -- the list contains all the resources that form the path keep in the hierarchy order
);

commit;