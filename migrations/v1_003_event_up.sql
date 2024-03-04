-- CREATE TYPE event_status AS ENUM ('editing','review', 'published');

CREATE TABLE event
(
    id           char(21)    NOT NULL PRIMARY KEY,
    name         varchar(64) NOT NULL UNIQUE,
    description  text,
    image_url    text,
    music_genres text,
    line_up      json,
    start_time   timestamptz,
    end_time     timestamptz,
    tickets_url  text,
    price        json,
    min_age      int,
    promoter     varchar(128),
    location_id  char(10) REFERENCES location (id),
    is_public    boolean      NOT NULL,
    updated_at   timestamp    NOT NULL DEFAULT now(),
    created_at   timestamp    NOT NULL DEFAULT now()
);

-- CREATE INDEX "event_start_time_idx" on event (start_time);
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
-- CREATE TABLE event_interests
-- (
--     event_id char(26) REFERENCES event (id),
--     user_id  int REFERENCES "user" (id),
--     PRIMARY KEY (event_id, user_id)
-- );
-- CREATE TABLE genre
-- (
--     id          serial PRIMARY KEY,
--     name        varchar(16) UNIQUE,
--     description varchar(128)
-- );
-- CREATE TABLE event_genres
-- (
--     event_id char(26) REFERENCES event (id),
--     genre_id int REFERENCES genre (id),
--     PRIMARY KEY (event_id, genre_id)
-- );