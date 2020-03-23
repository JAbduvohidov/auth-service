package services

const usersDDL = `CREATE TABLE IF NOT EXISTS users
(
    id       BIGSERIAL PRIMARY KEY,
    name     TEXT NOT NULL,
    surname  TEXT NOT NULL,
    login    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    avatar   TEXT NOT NULL,
    role     TEXT    DEFAULT 'USER',
    removed  BOOLEAN DEFAULT FALSE
);`

const moderatorDML = `INSERT INTO users (id, name, surname, login, password, avatar, role)
VALUES (1, 'Moderator', 'Moderator', 'moderator', 'moderator', 'https://i.pravatar.cc/50', 'MODERATOR') ON CONFLICT DO NOTHING;`
