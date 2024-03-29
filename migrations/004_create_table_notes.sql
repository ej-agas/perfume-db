create table notes(
    id serial primary key,
    public_id varchar not null,
    slug text not null,
    name varchar not null,
    description text,
    image_url varchar,
    note_group_id varchar not null,
    created_at timestamp,
    updated_at timestamp,
    constraint fk_note_group_id foreign key (note_group_id) references note_groups (public_id)
);

create unique index notes_unique_slug__idx on notes (slug);
create unique index notes_unique_name__idx on notes (name);
create unique index notes_unique_public_id__idx on notes (public_id);

---- create above / drop below ----

drop table notes;

-- ID          int
-- Name        string
-- Slug        string
-- Description string
-- ImageURL    string
-- CreatedAt   time.Time
-- UpdatedAt   time.Time