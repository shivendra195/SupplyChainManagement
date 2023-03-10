CREATE TABLE IF NOT EXISTS country
(
    id           BIGSERIAL PRIMARY KEY,
    country      text NOT NULL,
    country_code text NOT NULL,
    archived_at  timestamp
);

CREATE TABLE IF NOT EXISTS state
(
    id          BIGSERIAL PRIMARY KEY,
    country_id  int references country (id) NOT NULL,
    state       text                        NOT NULL,
    archived_at timestamp
);