create table houses(
    id serial primary key,
    public_id varchar not null,
    slug text not null,
    name varchar not null,
    country varchar not null,
    description text,
    year_founded timestamp,
    created_at timestamp,
    updated_at timestamp
);

create unique index houses_unique_public_id__idx on houses (public_id);
create unique index houses_unique_slug__idx on houses (slug);
create unique index houses_unique_name__idx on houses (name);
create index houses_country__idx on houses (country);

---- create above / drop below ----

drop table houses;
