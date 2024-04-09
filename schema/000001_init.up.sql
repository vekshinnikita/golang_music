CREATE TABLE users (
    id serial NOT NULL UNIQUE,
    name VARCHAR(255) not null,
    username VARCHAR(255) not null UNIQUE,
    password_hash VARCHAR(255) not null
);