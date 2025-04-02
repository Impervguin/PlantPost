-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS "user" (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS author (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    has_rights BOOLEAN NOT NULL,
    grant_at TIMESTAMPTZ DEFAULT NOW(),
    revoke_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES "user"(id)
);

CREATE TABLE IF NOT EXISTS plant_category (
    name TEXT PRIMARY KEY,
    attributes JSONB NOT NULL,
    photo_id UUID
);

CREATE TABLE IF NOT EXISTS "file" (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS plant (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    latin_name TEXT NOT NULL,
    description TEXT NOT NULL,
    category TEXT,
    main_photo_id UUID,
    specification JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (main_photo_id) REFERENCES file(id),
    FOREIGN KEY (category) REFERENCES plant_category(name)
);

CREATE TABLE IF NOT EXISTS plant_photo (
    id UUID PRIMARY KEY,
    plant_id UUID NOT NULL,
    file_id UUID NOT NULL,
    description TEXT NOT NULL,
    FOREIGN KEY (plant_id) REFERENCES plant(id),
    FOREIGN KEY (file_id) REFERENCES "file"(id)
);

CREATE TABLE IF NOT EXISTS post (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    author_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (author_id) REFERENCES author(id)
);

CREATE TABLE IF NOT EXISTS post_photo (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    file_id UUID NOT NULL,
    place_number int NOT NULL CHECK (place_number > 0),
    FOREIGN KEY (post_id) REFERENCES post(id),
    FOREIGN KEY (file_id) REFERENCES "file"(id)
);

CREATE TABLE IF NOT EXISTS plant_post (
    id UUID PRIMARY KEY,
    plant_id UUID NOT NULL,
    post_id UUID NOT NULL,
    FOREIGN KEY (plant_id) REFERENCES plant(id),
    FOREIGN KEY (post_id) REFERENCES post(id)
);



CREATE TABLE IF NOT EXISTS post_tag (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    tag TEXT NOT NULL,
    FOREIGN KEY (post_id) REFERENCES post(id)
);

CREATE TABLE IF NOT EXISTS album (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    owner_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (owner_id) REFERENCES "user"(id)
);

CREATE TABLE IF NOT EXISTS plant_album (
    id UUID PRIMARY KEY,
    plant_id UUID NOT NULL,
    album_id UUID NOT NULL,
    FOREIGN KEY (plant_id) REFERENCES plant(id),
    FOREIGN KEY (album_id) REFERENCES album(id)
);

-- +goose StatementEnd
-- +goose Down

DROP TABLE IF EXISTS plant_album;
DROP TABLE IF EXISTS album;
DROP TABLE IF EXISTS post_tag;
DROP TABLE IF EXISTS plant_post;
DROP TABLE IF EXISTS post_photo;
DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS plant_photo;
DROP TABLE IF EXISTS plant;
DROP TABLE IF EXISTS plant_category;
DROP TABLE IF EXISTS "file";
DROP TABLE IF EXISTS author;
DROP TABLE IF EXISTS "user";

-- +goose StatementBegin
-- +goose StatementEnd