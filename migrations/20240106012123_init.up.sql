create table goals
(
    id        integer  not null
        primary key autoincrement,
    usr_id    integer  not null,
    chat_id   text     not null,
    message   text     not null,
    status_id integer  not null,
    deadline  DATETIME not null
);

create table notifier
(
    id           integer  not null
        primary key autoincrement,
    usr_id       integer  not null,
    chat_id      text     not null,
    last_updated DATETIME not null,
    goal_id      integer  not null
        references goals
            on delete cascade,
    status_id    integer  not null,
    notify       integer  not null
);

