begin;

create schema if not exists demo;
create extension if not exists "pgcrypto" with schema demo;

set schema 'demo';

-- Functions and Triggers
create
    or replace function nownano(curtime timestamptz)
    returns bigint as
$$
begin
    return EXTRACT(EPOCH FROM curtime) * 1000;
end;
$$
    language 'plpgsql' IMMUTABLE;

create
    or replace function trigger_updated_at() returns trigger
    language plpgsql
as
$$
begin
    new.updated_at
        = demo.nownano(CURRENT_TIMESTAMP);
    return new;
end;
$$;

create
    or replace function label_and_desc_tsearch(label varchar, description varchar)
    returns tsvector as
$$
declare
begin
    return setweight(to_tsvector(coalesce(label, '')), 'A') || setweight(to_tsvector(coalesce(description, '')), 'B');
end;
$$
    language 'plpgsql' IMMUTABLE;

-- experimental
create or replace function uri_alias() returns varchar as
$$
declare
    dict    varchar = 'abcdefghijklmnpqrstuvwxyzABCDEFGHIJLMNOPQRSTUVWXYZ0123456789'; -- We will remove two characters o and K
    randstr varchar;
begin
    select into randstr string_agg(substr(dict, ceil(random() * 60)::integer, 1), '') FROM generate_series(1, 22);
    return randstr;
end;
$$
    language 'plpgsql' immutable;

-- Insert trigger for org
create or replace function trigger_org_uri() returns trigger as
$$
declare
    found bool = false;
begin
    while found
        loop
            new.path_uri =  demo.uri_alias();
            select into found exists(select 1 from demo.org where path_uri = new.path_uri);
        END LOOP;
    return new;
end;
$$
    language 'plpgsql';

create or replace function trigger_repo_uri() returns trigger as
$$
declare
    found bool = false;
begin
    while found
        loop
            new.path_uri =  demo.uri_alias();
            select into found exists(select 1 from demo.repository where path_uri = new.path_uri);
        END LOOP;
    return new;
end;
$$
    language 'plpgsql';

create or replace function trigger_folder_uri() returns trigger as
$$
declare
    found bool = false;
begin
    while found
        loop
            new.path_uri =  demo.uri_alias();
            select into found exists(select 1 from demo.folder where path_uri = new.path_uri);
        END LOOP;
    return new;
end;
$$
    language 'plpgsql';

create or replace function trigger_project_uri() returns trigger as
$$
declare
    found bool = false;
begin
    while found
        loop
            new.path_uri =  demo.uri_alias();
            select into found exists(select 1 from demo.project where path_uri = new.path_uri);
        END LOOP;
    return new;
end;
$$
    language 'plpgsql';

--  Table Definitions
CREATE TABLE base
(
    id           uuid    not null primary key default gen_random_uuid(),
    path_uri     varchar not null,
    previous_uri varchar,
    type         varchar not null,
    hierarchy    jsonb,
    label        varchar,
    alt_label    varchar,
    description  varchar,
    tsearch      tsvector generated always as ( label_and_desc_tsearch(label, description) ) stored,
    owner        varchar not null,
    updated_by   varchar,
    created_at   bigint default nownano(CURRENT_TIMESTAMP),
    updated_at   bigint default nownano(CURRENT_TIMESTAMP),
    deleted      bool default false,
    constraint pathing unique (path_uri, label, type)
);

CREATE TABLE org
(
    LIKE base including all
) inherits (base);

CREATE TABLE repository
(
    LIKE base including all
) inherits (base);

CREATE TABLE folder
(
    LIKE base including all
) INHERITS (base);

CREATE TABLE project
(
    LIKE base including all,
    colour varchar
) INHERITS (base);



-- trigger assignments
create trigger org_updated_at
    before update
    on org
    for each row
execute procedure trigger_updated_at();
create trigger repository_updated_at
    before update
    on repository
    for each row
execute procedure trigger_updated_at();
create trigger project_updated_at
    before update
    on project
    for each row
execute procedure trigger_updated_at();
create trigger folder_updated_at
    before update
    on folder
    for each row
execute procedure trigger_updated_at();

--triggers insert resources
create trigger org_insert before insert on demo.org for each row execute procedure demo.trigger_org_uri();
create trigger folder_insert before insert on demo.folder for each row execute procedure demo.trigger_folder_uri();
create trigger project_insert before insert on demo.project for each row execute procedure demo.trigger_project_uri();
create trigger repo_insert before insert on demo.repository for each row execute procedure demo.trigger_repo_uri();


set schema 'public';

end;




