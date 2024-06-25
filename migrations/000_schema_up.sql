CREATE TABLE artwork
(
    id           serial PRIMARY KEY,
    url          text UNIQUE NOT NULL,
    width        int,
    height       int,
    bg_color     char(6),
    text_color_1 char(6),
    text_color_2 char(6),
    text_color_3 char(6),
    text_color_4 char(6)
);

CREATE TYPE venue_type AS ENUM ('club', 'bar', 'cafe', 'space','concert_hall');

CREATE TABLE venue
(
    id             serial PRIMARY KEY,
    name           varchar(64) UNIQUE NOT NULL,
    type           venue_type,
    description    varchar(256)       NOT NULL,
    artwork_id     int REFERENCES artwork (id),
    address        text               NOT NULL,
    metro_stations text               NOT NULL,
    latitude       int,
    longitude      int,
    is_public      boolean            NOT NULL DEFAULT false,
    created_at     timestamptz        NOT NULL DEFAULT now(),
    updated_at     timestamptz        NOT NULL DEFAULT now()
);

CREATE TABLE venue_stages
(
    venue_id   int          NOT NULL,
    stage_id   int          NOT NULL,
    stage_name varchar(256) NOT NULL,
    PRIMARY KEY (venue_id, stage_id),
    FOREIGN KEY (venue_id) REFERENCES venue (id),
    UNIQUE (venue_id, stage_id, stage_name)
);

CREATE TABLE artist
(
    id   serial PRIMARY KEY,
    name varchar(64) UNIQUE NOT NULL
);

CREATE TABLE event
(
    id              serial PRIMARY KEY,
    title           varchar(64)  NOT NULL,
    description     text         NOT NULL,
    date            date         NOT NULL,
    start_time      time         NOT NULL,
    end_time        time         NOT NULL,
    age_restriction int          NOT NULL,
    promoter        varchar(128) NOT NULL,
    tickets_url     text         NOT NULL,
    is_public       boolean      NOT NULL DEFAULT false,
    artwork_id      int REFERENCES artwork (id),
    venue_id        int REFERENCES venue (id),
    updated_at      timestamp    NOT NULL DEFAULT now(),
    created_at      timestamp    NOT NULL DEFAULT now(),
    UNIQUE (title, date)
);

CREATE INDEX "index_event_date" on event (date);
CREATE INDEX "index_event_public" on event (id) where is_public is true;

CREATE TABLE event_lineup
(
    event_id    int,
    stage_name  int,
    artist_name int,
    is_live     boolean,
    start_at    time,
    FOREIGN KEY (event_id) REFERENCES event (id),
    PRIMARY KEY (event_id, stage_name, artist_name)
);

CREATE TABLE event_tickets
(
    event_id    int REFERENCES event (id),
    title       varchar(64) NOT NULL,
    price       int         NOT NULL,
    description text,
    PRIMARY KEY (event_id, title)
);

CREATE TABLE event_promo
(
    event_id int REFERENCES event (id),
    title    varchar(64) NOT NULL,
    terms    text        NOT NULL,
    PRIMARY KEY (event_id, title)
);

CREATE TABLE genre
(
    id          serial PRIMARY KEY,
    name        varchar(32) UNIQUE NOT NULL,
    description varchar(128)
);

CREATE TABLE event_genres
(
    event_id int,
    genre_id int,
    PRIMARY KEY (event_id, genre_id)
);



-- CREATE TABLE event_likes
-- (
--     event_id int REFERENCES event (id),
--     user_id  int REFERENCES "user" (id),
--     PRIMARY KEY (event_id, user_id)
-- );
-- CREATE TYPE event_status AS ENUM ('editing','review', 'published');
-- CREATE TABLE event_reviews
-- (
--     event_id    char(26) REFERENCES event (id)              NOT NULL,
--     user_id     int REFERENCES "user" (id)             NOT NULL,
--     stars       smallint CHECK (stars BETWEEN 1 and 5) NOT NULL,
--     title       varchar(128)                           NOT NULL,
--     description TEXT,
--     created_at  timestamp DEFAULT CURRENT_TIMESTAMP,
--     PRIMARY KEY (event_id, user_id)
-- );

-- CREATE TABLE venue_type
-- (
--     id   int PRIMARY KEY,
--     name varchar(16) NOT NULL UNIQUE
-- );
--
-- INSERT INTO venue_type (id, name)
-- VALUES (1, 'default'),
--        (2, 'club'),
--        (3, 'bar'),
--        (4, 'cafe'),
--        (5, 'space'),
--        (6, 'concert_hall');

-- INSERT INTO location (id, name, description, artwork_url, stages, address, metro_stations)
-- VALUES (-1, 'DEFAULT', 'This location used as default when no location provided', '', '', '', '');