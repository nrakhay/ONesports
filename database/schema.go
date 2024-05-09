package database

const Schema = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name TEXT,
    last_name TEXT,
    email TEXT,
    phone_number TEXT
);

CREATE TABLE IF NOT EXISTS voice_channel_recordings (
    id SERIAL PRIMARY KEY,
    channel_id TEXT,
    file_path TEXT
);
`
