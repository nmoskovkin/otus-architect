CREATE TABLE users
(
    id         varchar(36),
    first_name VARCHAR(255)      NOT NULL,
    last_name  VARCHAR(255)      NOT NULL,
    age        SMALLINT UNSIGNED NOT NULL,
    interests  TEXT              NOT NULL,
    city       VARCHAR(255)      NOT NULL
);
