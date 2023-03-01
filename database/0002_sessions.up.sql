create table if not exists sessions
(
    id          bigserial primary key,
    token       text not null,
    user_id     int references users (id),
    start_time  timestamp with time zone not null,
    end_time    timestamp with time zone,
    device_id   text,
    platform    text,
    model_name  text,
    os_version  text
)