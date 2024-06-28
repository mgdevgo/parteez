CREATE TABLE "role"
(
    id          serial PRIMARY KEY,
    name        varchar(64),
    description varchar(128)
);
INSERT INTO "role" (name)
VALUES ('owner'),
       ('editor'),
       ('manager'),
       ('service');

CREATE TABLE "user"
(
    "id"                serial PRIMARY KEY,
    "name"              varchar(32) UNIQUE NOT NULL,
    "email"             varchar(64) UNIQUE NOT NULL,
    "email_verified"    boolean            NOT NULL,
    "email_verified_at" timestamp,
    "password"          varchar(256)       NOT NULL
);
CREATE TABLE admin
(
    id      serial PRIMARY KEY,
    user_id int REFERENCES "user" (id),
    role    int REFERENCES role (id) DEFAULT 2
);
CREATE TABLE user_sessions
(
    session_id    serial,
    user_id       int REFERENCES "user" (id),
    refresh_token varchar(1024) NOT NULL,
    expires_at    timestamp     NOT NULL,
    PRIMARY KEY (session_id, user_id)
);