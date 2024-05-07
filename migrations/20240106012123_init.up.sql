create table if not exists goals
(
    id        integer  not null
        primary key autoincrement,
    usr_id    integer  not null,
    chat_id   text     not null,
    message   text     not null,
    status_id integer  not null,
    deadline  DATETIME not null,
    timer_enabled bool not null,
    timer     integer,
    last_updated DATETIME not null
);
