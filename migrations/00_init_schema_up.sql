create table
    artworks (
        id serial primary key,
        image_url text unique not null,
        width int,
        height int,
        bg_color char(6),
        text_color_1 char(6),
        text_color_2 char(6),
        text_color_3 char(6),
        text_color_4 char(6)
    );

create type venue_type as enum('club', 'bar', 'cafe', 'space', 'concert_hall');

create table
    venues (
        id serial primary key,
        name varchar(64) unique not null,
        type venue_type,
        description text default '',
        artwork_id int references artworks (id),
        stages text[],
        address text default '',
        metro_stations text[],
        public boolean default false,
        created_at timestamptz default now(),
        updated_at timestamptz default now()
    );

create table
    events (
        id serial primary key,
        title varchar(64) not null,
        description text default '',
        date date not null,
        start_time time not null,
        end_time time null,
        age_restriction int default 18,
        promoter varchar(128),
        tickets_url text,
        artwork_id int references artworks (id),
        venue_id int references venues (id),
        public boolean default false,
        updated_at timestamp default now(),
        created_at timestamp default now(),
        unique (title, date)
    );

create index events_date_idx on events (date);
create index events_is_public_idx on events (id) where public is true;

create table
    events_lineup (
        event_id int references events(id),
        stage varchar(64),
        artist_name varchar(64),
        live boolean,
        start_at time,
        primary key (event_id, stage_name, artist_name)
    );

create table
    genres (
        id serial primary key,
        name varchar(32) unique not null,
        description varchar(128) default ''
    );

create table
    events_genres (
        event_id int,
        genre_id int,
        primary key (event_id, genre_id)
    );

create table
    events_tickets (
        id int primary key,
        event_id int references events (id),
        title varchar(64) not null,
        price int not null,
        description text default '',
        unique (event_id, title)
    );