CREATE TABLE IF NOT EXISTS  users ON auth (
    username        VARCHAR(30) PRIMARY KEY NOT NULL,
    email           VARCHAR(320) NOT NULL,
    "password"      VARCHAR(100) NOT NULL
);