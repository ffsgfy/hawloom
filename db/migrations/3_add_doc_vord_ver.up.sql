-- NOTE: vord = VOting RounD

BEGIN;
    CREATE TABLE doc (
        id UUID NOT NULL PRIMARY KEY,
        title TEXT NOT NULL,
        flags INT NOT NULL,
        created_by INT NOT NULL REFERENCES account (id),
        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        vord_duration INT NOT NULL CHECK (vord_duration > 0)
    );

    CREATE TABLE vord (
        doc UUID NOT NULL REFERENCES doc (id) ON DELETE CASCADE,
        num INT NOT NULL CHECK (num >= -1), -- -1 is the currently active vord
        flags INT NOT NULL,
        start_at TIMESTAMP WITH TIME ZONE NOT NULL,
        finish_at TIMESTAMP WITH TIME ZONE NOT NULL,

        PRIMARY KEY (doc, num)
    );

    CREATE TABLE ver (
        id UUID NOT NULL PRIMARY KEY,
        doc UUID NOT NULL REFERENCES doc (id) ON DELETE CASCADE,
        vord_num INT NOT NULL,
        votes INT NOT NULL DEFAULT 0,
        votes_updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        created_by INT REFERENCES account (id),
        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        summary TEXT NOT NULL,
        content TEXT NOT NULL,

        FOREIGN KEY (doc, vord_num) REFERENCES vord (doc, num) ON UPDATE CASCADE
    );

    CREATE INDEX ON vord (finish_at)
    WHERE num = -1;

    CREATE INDEX ON ver (doc, vord_num, votes);
COMMIT;
