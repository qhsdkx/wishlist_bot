create table if not exists wishes
(
    id BIGSERIAL PRIMARY KEY,
    wish_text text,
    user_id BIGINT REFERENCES users(id)
);

CREATE INDEX wish_idx ON wishes (wish_text);