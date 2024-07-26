create table artworks
(
    id           serial,
    name         text not null,
    width        int,
    height       int,
    bg_color     char(6),
    text_color_1 char(6),
    text_color_2 char(6),
    text_color_3 char(6),
    text_color_4 char(6),
    primary key (id)
);

create type visibility as enum ('public', 'private');

create type venue_type as enum ('unknown', 'club', 'bar', 'restaurant', 'space', 'concert_hall', 'open_air');

create table venues
(
    id             serial,
    name           varchar(64) not null,
    type           venue_type           default 'unknown',
    description    text        not null default '',
    stages         text[]      not null default '{"main"}',
    artwork_id     int,

    address        text        not null default '',
    metro_stations text[]      not null default '{}',

    visibility     visibility  not null default 'private',
    created_at     timestamptz not null default now(),
    updated_at     timestamptz not null default now(),

    primary key (id),
    foreign key (artwork_id) references artworks (id),
    unique (name)
);

create type event_status as enum ('draft', 'review_required', 'approved');

create table events
(
    id              serial,
    title           varchar(64)  not null,
    description     text         not null default '',
    artwork_id      int,

    -- Period when event is happening
    -- TODO: choose one
    date            date         not null,
    time            tstzrange    not null,
    start_time      time         not null,
    end_time        time         null,

    age_restriction int          not null default 18,

    promoter        varchar(128) not null default '',
    venue_id        int,
    tickets_url     text         not null default '',

    status          event_status not null default 'draft',
    visibility      visibility   not null default 'private',

    updated_at      timestamptz  not null default now(),
    created_at      timestamptz  not null default now(),
    primary key (id),
    foreign key (artwork_id) references artworks (id),
    foreign key (venue_id) references venues (id),
    unique (title, date)
);

create index events_date_idx on events (date);
create index events_is_public_idx on events (id) where visibility = 'public';

create table artists
(
    id          serial,
    name        text not null,
    social_link text not null default '',
    primary key (id),
    unique (name)
);

create table lineups
(
    id        serial,
    event_id  int,
    stage     text not null default 'main',
    with_time boolean       default false,
    primary key (id),
    foreign key (event_id) references events (id) on delete cascade,
    unique (event_id, stage)
);

create index events_lineup_idx on lineups (event_id);

create table lineups_artists
(
    lineup_id int     not null,
    artist_id int     not null,
    live      boolean not null default false,
    start_at  time,
    b2b       int,
    primary key (lineup_id, artist_id),
    foreign key (lineup_id) references lineups (id) on delete cascade,
    foreign key (event_id) references events (id),
    foreign key (artist_id) references artists (id),
    foreign key (b2b) references artists (id)
);

-- create table events_lineups
-- (
--     event_id  int,
--     lineup_id int,
--     foreign key (event_id) references events (id),
--     foreign key (lineup_id) references lineups (id)
-- );

-- create table events_lineups
-- (
--     event_id  int     not null,
--     stage     text    not null default 'main',
--     artist_id int     not null,
--     live      boolean not null default false,
--     start_at  time,
--     b2b       int,
--     primary key (event_id, stage, artist_id),
--     foreign key (event_id) references events (id) on delete cascade,
--     foreign key (artist_id) references artists (id),
--     foreign key (b2b) references artists (id)
-- );

create table genres
(
    id          serial,
    name        varchar(32) not null,
    description varchar(128) default '',
    primary key (id),
    unique (name)
);

create table events_genres
(
    event_id int not null,
    genre_id int not null,
    primary key (event_id, genre_id)
);

create table events_tickets
(
    id          serial,
    event_id    int         not null,
    title       varchar(64) not null,
    price       int         not null,
    description text        not null default '',
    primary key (id),
    unique (event_id, title),
    foreign key (event_id) references events (id) on delete cascade
);