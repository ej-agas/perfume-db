create table perfumes_notes(
    perfume_id uuid not null,
    note_id uuid not null,
    category varchar,
    constraint fk_perfume_id foreign key (perfume_id) references perfumes (id),
    constraint fk_note_id foreign key (note_id) references notes (id),
    constraint unique_perfume_id_note_id unique (perfume_id, note_id)
);

create index perfumes_notes_category__idx on perfumes_notes (category);

---- create above / drop below ----

drop table perfumes_notes;
