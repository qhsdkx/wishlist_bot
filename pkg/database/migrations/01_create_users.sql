create table if not exists users
(
    id         bigserial
        primary key,
    name       text
        constraint users_name_check
            check (length(name) <= 64),
    surname    text
        constraint users_surname_check
            check (length(surname) <= 64),
    username   text
        constraint users_username_check
            check (length(username) <= 64),
    birthdate  date,
    status     text not null,
    created_at timestamp default now(),
    updated_at timestamp,
    deleted_at timestamp
);
