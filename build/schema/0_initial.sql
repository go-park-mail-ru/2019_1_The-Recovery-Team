CREATE extension IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS profile (
    id SERIAL PRIMARY KEY,
    email citext UNIQUE NOT NULL CHECK (email <> ''),
    password varchar(64) NOT NULL CHECK (password <> ''),
    nickname citext UNIQUE NOT NULL CHECK (nickname <> ''),
    avatar TEXT DEFAULT '',

    record INTEGER DEFAULT 0,
    win INTEGER DEFAULT 0,
    loss INTEGER DEFAULT 0
) WITH (autovacuum_enabled = true, autovacuum_analyze_threshold = 3);
