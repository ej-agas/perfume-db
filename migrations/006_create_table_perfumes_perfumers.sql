create table perfumes_perfumers(
    perfume_id serial not null,
    perfumer_id serial not null,
    constraint fk_perfume_id foreign key (perfume_id) references perfumes (id),
    constraint fk_perfumer_id foreign key (perfumer_id) references perfumers (id)
);

---- create above / drop below ----

drop table perfumes_perfumers;

-- 	ID               int
-- 	Slug             string
-- 	Name             string
-- 	Description      string
-- 	Concentration    Concentration
-- 	ImageURL         string
-- 	House            House
-- 	Perfumers        []*Perfumer
-- 	Notes            map[NoteCategory][]*Note
-- 	YearReleased     time.Time
-- 	YearDiscontinued time.Time
-- 	CreatedAt        time.Time
-- 	UpdatedAt        time.Time