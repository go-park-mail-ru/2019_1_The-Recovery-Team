CREATE extension IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS profile (
    id SERIAL PRIMARY KEY,
    email citext UNIQUE,
    password varchar(64),
    nickname citext UNIQUE NOT NULL,
    avatar TEXT DEFAULT '',

    record INTEGER DEFAULT 1500,
    win INTEGER DEFAULT 0,
    loss INTEGER DEFAULT 0
) WITH (autovacuum_enabled = true, autovacuum_analyze_threshold = 3);

CREATE TABLE IF NOT EXISTS token (
    id SERIAL PRIMARY KEY,
    oauth citext NOT NULL,
    user_id citext NOT NULL,
    profile_id INTEGER REFERENCES profile(id) UNIQUE,
    token TEXT NOT NULL
)
