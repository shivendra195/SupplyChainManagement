CREATE TABLE users
(
    id           BIGSERIAL PRIMARY KEY,
    name         text      NOT NULL,
    age          int       NOT NULL,
    password     text      NOT NULL,
    address      text      NOT NULL,
    country_code text      NOT NULL,
    email        text      NOT NULL,
    phone        text      NOT NULL,
    created_at   timestamp NOT NULL,
    updated_at   timestamp NOT NULL,
    archived_at  timestamp
);

CREATE TABLE user_profiles
(
    id           BIGSERIAL PRIMARY KEY,
    user_id      int references users (id) NOT NULL,
    company_name text,
    country      text,
    state        text,
    created_at   timestamp                 NOT NULL,
    updated_at   timestamp                 NOT NULL,
    archived_at  timestamp
);


CREATE TYPE roles AS ENUM ('admin', 'dealer', 'retailer');


CREATE TABLE user_roles
(
    id      BIGSERIAL PRIMARY KEY,
    user_id int references users (id) NOT NULL,
    role    roles                     NOT NULL,
    archived_at  timestamp
);



