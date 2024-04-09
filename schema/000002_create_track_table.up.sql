CREATE TABLE tracks (
    id serial NOT NULL UNIQUE,
    user_id INTEGER not null REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) not null,
    description VARCHAR(255),
    author VARCHAR(255) not null,
    track_file_name VARCHAR(255),
    poster_file_name VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
