CREATE TABLE event
(
    id              serial PRIMARY KEY,
    name            varchar(64)  NOT NULL,
    description     text         NOT NULL,
    artwork_url     text         NOT NULL,
    genre           text         NOT NULL,
    start_date      timestamptz  NOT NULL,
    end_date        timestamptz  NOT NULL,
    line_up         jsonb        NOT NULL,
    age_restriction int          NOT NULL,
    price           jsonb        NOT NULL,
    tickets_url     text         NOT NULL,
    promoter        varchar(128) NOT NULL,
    location_id     int          NOT NULL REFERENCES location (id) DEFAULT -1,
    is_public       boolean      NOT NULL,
    updated_at      timestamp    NOT NULL                          DEFAULT now(),
    created_at      timestamp    NOT NULL                          DEFAULT now(),
    UNIQUE (name, start_time)
);

-- CREATE TABLE line_up
-- (
--     event_id    int          NOT NULL REFERENCES event (id),
--     artist_name varchar(128) NOT NULL,
--     start_time  char(5)      NOT NULL,
--     live        boolean      NOT NULL,
--     stage_id    int          NOT NULL REFERENCES location_stage (id) DEFAULT 1,
--     PRIMARY KEY (event_id, artist_name),
-- )


-- CREATE INDEX "event_start_time_idx" on event (start_time);
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