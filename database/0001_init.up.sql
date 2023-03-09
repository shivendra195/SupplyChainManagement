CREATE TABLE if not exists users
(
    id           BIGSERIAL PRIMARY KEY,
    name         text                      NOT NULL,
    password     text                      NOT NULL,
    address      text                      NOT NULL,
    country_code text                      NOT NULL,
    email        text                      NOT NULL,
    phone        text                      NOT NULL,
    created_by   int references users (id) ,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone default now(),
    archived_at  timestamp
);

CREATE TABLE IF NOT EXISTS profile_image
(
    id          BIGSERIAL PRIMARY KEY,
    url         TEXT                     NOT NULL,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE if not exists user_profiles
(
    id               BIGSERIAL PRIMARY KEY,
    user_id          int references users (id)         NOT NULL,
    company_name     text,
    date_of_birth    text                              NOT NULL,
    gender           text,
    country          text,
    profile_image_id int references profile_image (id) NOT NULL,
    state            text,
    created_at       timestamp with time zone,
    updated_at       timestamp with time zone default now(),
    archived_at      timestamp
);


CREATE TYPE roles AS ENUM ('super admin','admin', 'dealer', 'retailer');


CREATE TABLE if not exists user_roles
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     int references users (id) NOT NULL,
    role        roles                     NOT NULL,
    archived_at timestamp
);

create table if not exists sessions
(
    id         bigserial primary key,
    token      text                     not null,
    user_id    int references users (id),
    start_time timestamp with time zone not null,
    end_time   timestamp with time zone,
    device_id  text,
    platform   text,
    model_name text,
    os_version text
);

CREATE TYPE order_status AS ENUM ('open','in stock','in transfer','sold out','out of stock');

CREATE TABLE IF NOT EXISTS orders
(
    id               BIGSERIAL PRIMARY KEY,
    ordered_by       INT REFERENCES users (id),
    quantity         TEXT NOT NULL,
    reference_no     TEXT,
    shipping_address TEXT,
    order_status     order_status,
    completed_at     TIMESTAMP WITH TIME ZONE,
    created_at       TIMESTAMP WITH TIME ZONE,
    updated_at       TIMESTAMP WITH TIME ZONE default now(),
    archived_at      TIMESTAMP
);



CREATE TABLE IF NOT EXISTS ordered_items_scan_history
(
    id          BIGSERIAL PRIMARY KEY,
    order_id    INT REFERENCES orders (id) NOT NULL,
    data        json                       NOT NULL,
    scanned_at  TIMESTAMP WITH TIME ZONE   NOT NULL,
    archived_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS items_qr_data
(
    id          BIGSERIAL PRIMARY KEY,
    data        json NOT NULL,
    address     TEXT NOT NULL,
    scanned_at  TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP
);



