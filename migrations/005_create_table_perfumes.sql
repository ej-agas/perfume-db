create table perfumes(
    id serial primary key,
    public_id varchar not null,
    slug text not null,
    name varchar not null,
    description text,
    concentration smallint not null,
    image_url varchar,
    house_id varchar not null,
    year_released timestamp not null,
    year_discontinued timestamp,
    created_at timestamp,
    updated_at timestamp,
    constraint fk_house_id foreign key (house_id) references houses (public_id)
);

create unique index perfumes_unique_public_id__idx on perfumes (public_id);
create unique index perfumes_unique_slug__idx on perfumes (slug);
create unique index perfumes_unique_name__idx on perfumes (name);
create index perfumes_concentration__idx on perfumes (concentration);

---- create above / drop below ----

drop table perfumes;

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