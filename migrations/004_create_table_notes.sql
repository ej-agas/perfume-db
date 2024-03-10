create table notes(
    id serial primary key,
    slug text not null,
    name varchar not null,
    description text,
    photo_url varchar,
    note_group_id serial not null,
    created_at timestamp,
    updated_at timestamp,
    constraint fk_note_group_id foreign key (note_group_id) references note_groups (id)
);

create unique index notes_unique_slug__idx on notes (slug);
create unique index notes_unique_name__idx on notes (name);

---- create above / drop below ----

drop table notes;

-- ID          int
-- Name        string
-- Slug        string
-- Description string
-- ImageURL    string
-- CreatedAt   time.Time
-- UpdatedAt   time.Time