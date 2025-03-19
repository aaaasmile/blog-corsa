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
CREATE TABLE IF NOT EXISTS "post" (
	"id"	INTEGER,
	"title"	TEXT,
	"post_id"	TEXT NOT NULL,
	"timestamp"	NUMERIC,
	"abstract"	TEXT,
	"uri"	TEXT,
	-- "next_post_id"	TEXT,
	-- "prev_post_id"	TEXT,
	-- "content"	TEXT,
	-- "status" INTEGER,
	-- "tags"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE VIRTUAL TABLE postsearch USING fts5(post_rowid, content);
COMMIT;
