--Version: 1
--Description: Initial migration
CREATE TABLE IF NOT EXISTS events
(
    id         UUID PRIMARY KEY,
    title      TEXT,
    location   TEXT,
    event_date TEXT
);
