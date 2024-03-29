create table perfumers(
    id serial primary key,
    public_id varchar not null,
    slug text not null,
    name varchar not null,
    nationality varchar not null,
    image_url varchar,
    birth_date timestamp,
    created_at timestamp,
    updated_at timestamp
);

create unique index perfumers_unique_public_id__idx on perfumers (public_id);
create unique index perfumers_unique_slug__idx on perfumers (slug);
create unique index perfumers_unique_name__idx on perfumers (name);
create index perfumers_nationality__idx on perfumers (nationality);

---- create above / drop below ----

drop table perfumers;