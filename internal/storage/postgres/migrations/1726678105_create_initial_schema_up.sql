create table if not exists
    artworks (
        id serial primary key,
        name text not null,
        path text not null,
        width int,
        height int,
        bg_color char(6),
        text_color_1 char(6),
        text_color_2 char(6),
        text_color_3 char(6),
        text_color_4 char(6)
    );

insert into
    artworks (id, name, url)
values
    (-1, '', '');

create type venue_type as enum(
    'unknown',
    'club',
    'bar',
    'restaurant',
    'space',
    'concert_hall',
    'open_air'
);

create table if not exists
    venues (
        id serial primary key,
        name varchar(64) unique not null,
        type venue_type default 'unknown',
        description text not null default '',
        stages text[] not null default '{"main"}',
        artwork_id int not null default -1 references artworks (id) on delete set default,
        address text not null default '',
        metro_stations text[] not null default '{}',
        created_at timestamptz not null default now(),
        updated_at timestamptz not null default now(),
        is_draft boolean not null default true
    );

create table if not exists
    events (
        id serial primary key,
        title varchar(64) not null,
        description text not null default '',
        artwork_id int not null default -1 references artworks (id) on delete set default,
        -- Period when event is happening
        -- TODO: choose one
        date date not null,
        start_time time not null,
        end_time time null,
        age_restriction int not null default 18,
        promoter varchar(128) not null default '',
        venue_id int references venues (id) on delete set null,
        tickets_url text not null default '',
        updated_at timestamptz not null default now(),
        created_at timestamptz not null default now(),
        is_draft boolean not null default true,
        unique (title, date)
    );

create index if not exists events_date_idx on events (date);

create index if not exists events_is_public_idx on events (id)
where
    is_draft = false;

create table if not exists
    artists (
        id serial primary key,
        name text unique not null,
        social_profile_url text not null default ''
    );

-- This table represents the schedule or lineup of artists at the event, 
-- specifying when and where they are performing.
create table if not exists
    lineups (
        id serial primary key,
        event_id int not null references events (id) on delete cascade,
        stage text not null default 'main',
        unique (event_id, stage)
    );

create index if not exists lineups_event_id_idx on lineups (event_id);

create table if not exists
    lineups_artists (
        lineup_id int not null references lineups (id) on delete cascade,
        artist_id int not null references artists (id) on delete cascade,
        b2b_artist_id int,
        start_time time,
        live boolean not null default false,
        foreign key (b2b_artist_id, lineup_id) references lineups_artists (artist_id, lineup_id) on delete set null,
        primary key (lineup_id, artist_id)
    );

create table if not exists
    genres (
        id serial primary key,
        name varchar(32) unique not null,
        description varchar(128) default ''
    );

create table if not exists
    events_genres (
        event_id int not null references events (id) on delete cascade,
        genre_id int not null references genres (id) on delete cascade,
        primary key (event_id, genre_id)
    );

create table if not exists
    tickets (
        id serial primary key,
        event_id int not null references events (id) on delete cascade,
        title varchar(64) not null,
        price int not null,
        description text not null default '',
        unique (event_id, title)
    );