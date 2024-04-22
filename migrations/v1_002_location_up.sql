CREATE TABLE location_type
(
    id   int PRIMARY KEY,
    name varchar(16) NOT NULL UNIQUE
);

INSERT INTO location_type (id, name)
VALUES (1, 'UNKNOWN'),
       (2, 'CLUB'),
       (3, 'BAR'),
       (4, 'CAFE'),
       (5, 'SPACE'),
       (6, 'CONCERT_HALL');

CREATE TABLE location
(
    id               serial       NOT NULL,
    name             varchar(64)  NOT NULL,
    location_type_id int          NOT NULL DEFAULT 1,
    description      varchar(256) NOT NULL,
    artwork_url      text         NOT NULL,
    stages           text         NOT NULL,
    address          text         NOT NULL,
    metro_stations   text         NOT NULL,
    latitude         int,
    longitude        int,
    is_public        boolean      NOT NULL DEFAULT FALSE,
    created_at       timestamptz  NOT NULL DEFAULT now(),
    updated_at       timestamptz  NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    UNIQUE (name),
    FOREIGN KEY (location_type_id) REFERENCES location_type (id)
);

INSERT INTO location (id, name, description, artwork_url, stages, address, metro_stations)
VALUES (-1, 'DEFAULT', 'This location used as default when no location provided', '', '', '', '');