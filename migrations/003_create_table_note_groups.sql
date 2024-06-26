create table note_groups(
    id serial primary key,
    public_id varchar not null,
    slug text not null,
    name varchar not null,
    description text,
    image_url varchar,
    created_at timestamp,
    updated_at timestamp
);

create unique index note_groups_unique_public_id__idx on note_groups (public_id);
create unique index note_groups_unique_slug__idx on note_groups (slug);
create unique index note_groups_unique_name__idx on note_groups (name);

---- create above / drop below ----

drop table note_groups;
