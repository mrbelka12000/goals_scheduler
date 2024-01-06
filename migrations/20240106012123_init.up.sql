CREATE TABLE IF NOT EXISTS goals (
    id integer primary key autoincrement not null,
    usr_id integer not null,
    notifier_id integer not null,
    message text not null,
    status_id integer not null,
    deadline DATETIME not null,
    FOREIGN KEY(notifier_id) REFERENCES notifier(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS notifier (
    id integer primary key autoincrement not null,
    usr_id integer not null,
    status_id integer not null,
    ticker text not null,
    last_updated DATETIME not null,
    expires DATETIME not null
)