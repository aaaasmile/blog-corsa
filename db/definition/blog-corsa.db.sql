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
CREATE TABLE "race" (
	"id"	INTEGER,
	"name"	TEXT,
	"title"	TEXT,
	"meter_length"	INTEGER,
	"ascending_meter"	INTEGER,
	"descending_meter"	INTEGER,
	"rank_global"	INTEGER,
	"rank_gender"	INTEGER,
	"rank_class"	INTEGER,
	"class_name"	TEXT,
	"race_start_datetime"	TEXT,
	"sport_type_id"	INTEGER,
	"result_time"	TEXT,
	"pace_kmh"	REAL,
	"pace_minkm"	TEXT,
	"comment"	TEXT,
	"race_subtype_id"	INTEGER,
	"km_length"	REAL,
	"loop_number"	INTEGER,
	"loop_length"	REAL,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("race_subtype_id") REFERENCES "race_subtype"("id"),
	FOREIGN KEY("sport_type_id") REFERENCES "sport_type"("id")
)
COMMIT;
