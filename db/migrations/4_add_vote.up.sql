BEGIN;
    CREATE TABLE vote (
        ver UUID NOT NULL REFERENCES ver (id) ON DELETE CASCADE,
        doc UUID NOT NULL REFERENCES doc (id) ON DELETE CASCADE,
        vord_num INT NOT NULL,
        account INT NOT NULL REFERENCES account (id),

        PRIMARY KEY (ver, account),
        FOREIGN KEY (doc, vord_num) REFERENCES vord (doc, num) ON UPDATE CASCADE
    );

    CREATE INDEX ON vote (doc, vord_num, account);
COMMIT;
