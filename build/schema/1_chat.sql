CREATE TABLE IF NOT EXISTS message (
    id SERIAL PRIMARY KEY,
    author INTEGER,
    receiver INTEGER,
    created timestamp with time zone DEFAULT now() NOT NULL,
    edited boolean DEFAULT FALSE NOT NULL,
    text text NOT NULL
)