CREATE TABLE role
(
    id          serial PRIMARY KEY,
    name        varchar(64),
    description varchar(128)
);
INSERT INTO role (name)
VALUES ('owner'),
       ('editor'),
       ('manager'),
       ('service');
CREATE TABLE "user"
(
    id                serial PRIMARY KEY,
    name              varchar(32) UNIQUE NOT NULL,
    email             varchar(64) UNIQUE NOT NULL,
    email_verified_at timestamp,
    password          varchar(256)       NOT NULL
);
CREATE TABLE admin
(
    id      serial PRIMARY KEY,
    user_id int REFERENCES "user" (id),
    role    int REFERENCES role (id) DEFAULT 2
);
CREATE TYPE location_type as ENUM ('CONCERT_HALL', 'CLUB','BAR', 'CAFE','OPEN_AIR');
CREATE TABLE location
(
    id        serial PRIMARY KEY,
    name      varchar(64) UNIQUE NOT NULL,
    type      location_type,
    address   text,
    latitude  int,
    longitude int,
    stages    json
);
CREATE TYPE event_status AS ENUM ('editing','review', 'published');
CREATE TABLE event
(
    id          serial PRIMARY KEY,
    name        varchar(64) UNIQUE NOT NULL,
    image_url   text,
    description text,
    starts_at   timestamp,
    ends_at     timestamp,
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
CREATE INDEX "event_starts_at_idx" on event (starts_at);
CREATE TABLE event_reviews
(
    event_id    int REFERENCES event (id)              NOT NULL,
    user_id     int REFERENCES "user" (id)             NOT NULL,
    stars       smallint CHECK (stars BETWEEN 1 and 5) NOT NULL,
    title       varchar(128)                           NOT NULL,
    description TEXT,
    created_at  timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, user_id)
);
CREATE TABLE event_interests
(
    event_id int REFERENCES event (id),
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
    event_id int REFERENCES event (id),
    genre_id int REFERENCES genre (id),
    PRIMARY KEY (event_id, genre_id)
);
CREATE TABLE user_sessions
(
    session_id    serial,
    user_id       int REFERENCES "user" (id),
    refresh_token varchar(1024) NOT NULL,
    expires_at    timestamp     NOT NULL,
    PRIMARY KEY (session_id, user_id)
)