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

create type venue_type as enum(
    'DEFAULT',
    'CLUB',
    'BAR',
    'CAFE',
    'CONCERT_HALL',
    'SPACE'
);

create table if not exists
    venues (
        id serial primary key,
        name varchar(64) unique not null,
        type venue_type default 'DEFAULT',
        description text not null default '',
        stages text[] not null default '{"main"}',
        artwork_id int references artworks (id) on delete set null,
        address text not null default '',
        metro_stations text[] not null default '{}',
        created_at timestamptz not null default now(),
        updated_at timestamptz not null default now(),
        is_public boolean not null default false
    );

create table if not exists
    events (
        id serial primary key,
        title varchar(64) not null,
        description text not null default '',
        genres text[] not null default '{}',
        artwork_id int references artworks (id) on delete set null,
        date daterange not null default daterange(NULL, NULL, '[]'),
        venue_id int references venues (id) on delete set null,
        age_restriction int not null default 18,
        promoter varchar(128) not null default '',
        tickets_url text not null default '',
        tickets jsonb not null default '{}',
        updated_at timestamptz not null default now(),
        created_at timestamptz not null default now(),
        is_draft boolean not null default true,
        published_at timestamptz,
        unique (title, date)
    );

create index if not exists events_date_idx on events using gist (date);
create index if not exists events_is_public_idx on events (id, published_at)
where 
    published_at is not null;

create table if not exists
    lineups (
        id serial primary key,
        event_id int not null references events (id) on delete cascade,
        stage text not null default 'main',
        timetable jsonb not null default '{}',
        unique (event_id, stage)
    );

create index if not exists lineups_event_id_idx on lineups (event_id);

create table if not exists
    genres (
        id serial primary key,
        name varchar(32) unique not null,
        description varchar(128) default ''
    );