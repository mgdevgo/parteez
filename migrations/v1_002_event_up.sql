CREATE TYPE location_type as ENUM ('CONCERT_HALL', 'CLUB','BAR', 'CAFE','OPEN_AIR');
CREATE TABLE location
(
    id        serial PRIMARY KEY,
    name      varchar(64) UNIQUE NOT NULL,
    description varchar(256),
    type      location_type,
    address   text,
    latitude  int,
    longitude int,
    stages    json
);
CREATE TYPE event_status AS ENUM ('editing','review', 'published');
CREATE TABLE event
(
    id          char(26) UNIQUE PRIMARY KEY,
    name        varchar(64) UNIQUE NOT NULL,
    image_url   text,
    description text,
    start_time   timestamp,
    end_time     timestamp,
    line_up     json,
    location_id int REFERENCES location (id),
    promoter    varchar(128),
    tickets_url text,
    price       json,
    min_age     int,
    status      event_status DEFAULT 'editing',
    updated_at  timestamp    DEFAULT CURRENT_TIMESTAMP,
    created_at  timestamp    DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX "event_start_time_idx" on event (start_time);
CREATE TABLE event_reviews
(
    event_id    char(26) REFERENCES event (id)              NOT NULL,
    user_id     int REFERENCES "user" (id)             NOT NULL,
    stars       smallint CHECK (stars BETWEEN 1 and 5) NOT NULL,
    title       varchar(128)                           NOT NULL,
    description TEXT,
    created_at  timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, user_id)
);
CREATE TABLE event_interests
(
    event_id char(26) REFERENCES event (id),
    user_id  int REFERENCES "user" (id),
    PRIMARY KEY (event_id, user_id)
);
CREATE TABLE genre
(
    id          serial PRIMARY KEY,
    name        varchar(16) UNIQUE,
    description varchar(128)
);
CREATE TABLE event_genres
(
    event_id char(26) REFERENCES event (id),
    genre_id int REFERENCES genre (id),
    PRIMARY KEY (event_id, genre_id)
);