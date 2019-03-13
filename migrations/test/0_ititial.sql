-- +migrate Up
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
);

INSERT INTO profile (id, email, nickname, password) VALUES
(100, 'test@mail.ru', 'test', '$2a$10$GuVAH77ae3Ni0DlLbe2IQ.e3cIZxFkOOR0ztnAKDuEnxCt7OE8qee'),
(101, 'test3@mail.ru', 'test3', 'e3cIZxFkOOR0ztnAKDuEnxCt7OE8qee'),
(102, 'test5@mail.ru', 'test5', '$2a$10$GuVAH77ae3Ni0DlLbe2IQ.e3cIZxFkOOR0ztnAKDuEnxCt7OE8qee')
RETURNING id, email, nickname;

-- +migrate Down
DROP TABLE IF EXISTS profile;


