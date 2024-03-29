create table perfumes_notes(
    perfume_id varchar not null,
    note_id varchar not null,
    category varchar,
    constraint fk_perfume_id foreign key (perfume_id) references perfumes (public_id),
    constraint fk_note_id foreign key (note_id) references notes (public_id)
);

create index perfumes_notes_category__idx on perfumes_notes (category);

---- create above / drop below ----

drop table perfumes_notes;
