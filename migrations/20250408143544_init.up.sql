CREATE TABLE notifications (
                               id SERIAL PRIMARY KEY,
                               telegram_chat_id INTEGER NOT NULL,

                               schedule_type VARCHAR(20) NOT NULL CHECK (schedule_type IN ('timer', 'daily', 'once')),
    -- For daily notifications
                               day INTEGER,
                               hour INTEGER,
                               minute INTEGER,

    -- For timer notifications, the interval in seconds (e.g., 300 for 5 minutes, 36000 for 10 hours)
                               interval_seconds DOUBLE PRECISION,
                               next_execution TIMESTAMP,

    -- For one-time notifications, the exact date and time to execute the notification
                               scheduled_time TIMESTAMP,

                               deadline TIMESTAMP,
                               message TEXT NOT NULL,
                               status VARCHAR(20) DEFAULT 'st',
                               created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Optional: Index to speed up lookup of notifications that are due for processing
CREATE INDEX idx_next_execution ON notifications(next_execution);

