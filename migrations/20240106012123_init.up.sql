CREATE TABLE IF NOT EXISTS goals (
                                     id SERIAL PRIMARY KEY,
                                     usr_id INTEGER NOT NULL,
                                     chat_id TEXT NOT NULL,
                                     message TEXT NOT NULL,
                                     status_id TEXT NOT NULL,
                                     deadline timestamp WITH TIME ZONE DEFAULT NOW(),
                                     timer_enabled BOOLEAN NOT NULL,
                                     timer BIGINT,
                                     last_updated timestamp WITH TIME ZONE DEFAULT NOW()
);
