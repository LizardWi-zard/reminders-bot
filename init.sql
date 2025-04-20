-- init.sql
create table users
(
    id serial primary key,
    username text not null,
    chat_id bigint not null
);

create table reminders
(
    id serial primary key,
    user_id bigint not null references users,
    content text not null,
    interval interval not null,
    is_active boolean default true not null,
    last_checked timestamp default (now() AT TIME ZONE 'utc'::text) not null,
    unique (user_id, content, interval)
);