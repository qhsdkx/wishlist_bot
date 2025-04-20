create table if not exists users
(
    id         BIGSERIAL PRIMARY KEY,
    name       text check ( length(name) <= 64 ),
    surname    text check ( length(surname) <= 64 ),
    nickname   text check ( length(nickname) <= 64 ),
    birthdate  date,
    status     text NOT NULL,
    created_at timestamp default now(),
    deleted_at timestamp
);