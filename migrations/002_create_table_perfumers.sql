create table perfumers(
    id serial primary key,
    slug text not null,
    name varchar not null,
    nationality varchar not null,
    photo_url varchar,
    birth_date timestamp,
    created_at timestamp,
    updated_at timestamp
);

create unique index perfumers_unique_slug__idx on perfumers (slug);
create unique index perfumers_unique_name__idx on perfumers (name);
create index perfumers_nationality__idx on perfumers (nationality);

---- create above / drop below ----

drop table perfumers;