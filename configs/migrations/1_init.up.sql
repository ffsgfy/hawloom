BEGIN;
    CREATE TABLE account (
        id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        name VARCHAR(36) UNIQUE NOT NULL,
        pwd_hash BYTEA NOT NULL
    );

    CREATE TABLE tablet (
        id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        owner INT REFERENCES account(id),
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        content BYTEA NOT NULL,

        CONSTRAINT updated_after_created CHECK (updated_at >= created_at)
    );
COMMIT;
