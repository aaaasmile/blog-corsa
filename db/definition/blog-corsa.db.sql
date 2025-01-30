BEGIN TRANSACTION;
DROP TABLE IF EXISTS "comment";
CREATE TABLE IF NOT EXISTS "comment" (
	"id"	INTEGER,
	"parent_id"	INTEGER,
	"name"	TEXT,
	"email"	TEXT,
	"comment"	TEXT,
	"timestamp"	NUMERIC,
	"post_id"	TEXT NOT NULL,
	"status"	INTEGER,
	"req_id"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
);
COMMIT;
