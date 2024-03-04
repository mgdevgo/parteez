CREATE TABLE location_type
(
    id   int PRIMARY KEY,
    name varchar(16) UNIQUE NOT NULL
);

INSERT INTO location_type (id, name)
VALUES
    (1, 'UNKNOWN'),
    (2, 'CLUB'),
    (3, 'BAR'),
    (4, 'CAFE'),
    (5, 'SPACE'),
    (6, 'CONCERT_HALL');

CREATE TABLE location
(
    id               char(10)    NOT NULL PRIMARY KEY,
    name             varchar(64) NOT NULL UNIQUE,
    location_type_id int         NOT NULL REFERENCES location_type (id) DEFAULT 1,
    description      varchar(256),
    image_url        text,
    stages           text,
    music_genres     text,
    address          text,
    metro_stations   text,
    latitude         int,
    longitude        int,
    is_public        boolean     NOT NULL,
    created_at       timestamptz NOT NULL DEFAULT now(),
    update_at        timestamptz NOT NULL DEFAULT now()
);