CREATE TABLE IF NOT EXISTS "file" (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS app_user (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS author (
    id UUID PRIMARY KEY,
    has_rights BOOLEAN NOT NULL,
    grant_at TIMESTAMPTZ DEFAULT NOW(),
    revoke_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (id) REFERENCES app_user(id),
    CONSTRAINT rights_timing_check CHECK (
        (has_rights AND grant_at > revoke_at) OR
        (NOT has_rights AND revoke_at > grant_at)
    )
);

CREATE TABLE IF NOT EXISTS plant_category (
    name TEXT PRIMARY KEY,
    attributes JSONB NOT NULL,
    photo_id UUID,
    CONSTRAINT photo_id_fk FOREIGN KEY (photo_id) REFERENCES "file"(id)
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
    CONSTRAINT main_photo_id_fk FOREIGN KEY (main_photo_id) REFERENCES "file"(id),
    CONSTRAINT category_fk FOREIGN KEY (category) REFERENCES plant_category(name),
    CONSTRAINT created_before_updated CHECK (created_at <= updated_at)
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
    FOREIGN KEY (author_id) REFERENCES author(id),
    CONSTRAINT created_before_updated CHECK (created_at <= updated_at)
);

CREATE TABLE IF NOT EXISTS post_photo (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    file_id UUID NOT NULL,
    place_number int NOT NULL,
    FOREIGN KEY (post_id) REFERENCES post(id),
    FOREIGN KEY (file_id) REFERENCES "file"(id),
    CONSTRAINT place_number_positive CHECK (place_number > 0),
    CONSTRAINT place_number_unique UNIQUE (post_id, place_number)
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
    FOREIGN KEY (owner_id) REFERENCES app_user(id),
    CONSTRAINT created_before_updated CHECK (created_at <= updated_at)
);

CREATE TABLE IF NOT EXISTS plant_album (
    id UUID PRIMARY KEY,
    plant_id UUID NOT NULL,
    album_id UUID NOT NULL,
    FOREIGN KEY (plant_id) REFERENCES plant(id),
    FOREIGN KEY (album_id) REFERENCES album(id)
);

-- check that user is created before grant_at in author
CREATE OR REPLACE FUNCTION check_user_created_before_grant()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM app_user 
        WHERE id = NEW.id AND created_at >= NEW.grant_at
    ) THEN
        RAISE EXCEPTION 'User must be created before grant_at';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER author_created_before_grant
BEFORE INSERT OR UPDATE ON author
FOR EACH ROW EXECUTE FUNCTION check_user_created_before_grant();
