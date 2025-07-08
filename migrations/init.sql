CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL
);

CREATE TABLE todos (
    id     UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    text   TEXT NOT NULL,
    done   BOOLEAN NOT NULL DEFAULT false,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO users (id, name)
VALUES ('f47ac10b-58cc-4372-a567-0e02b2c3d479', 'Alice');

INSERT INTO todos (id, text, done, user_id)
VALUES ('f47ac10b-58cc-4372-a567-0e02b2c3d479', 'Finish writing the doc', false, 'f47ac10b-58cc-4372-a567-0e02b2c3d479');

