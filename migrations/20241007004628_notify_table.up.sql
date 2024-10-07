CREATE TABLE IF NOT EXISTS notify (
    id serial primary key,
    hour integer not null ,
    minute integer not null,
    weekday integer not null,
    goal_id integer not null
);

ALTER TABLE goals
    ADD COLUMN IF NOT EXISTS notify_enabled boolean default false;

