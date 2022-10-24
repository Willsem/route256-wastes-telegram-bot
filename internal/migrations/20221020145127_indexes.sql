-- create index "user_id" to table: "users"
CREATE UNIQUE INDEX "user_id" ON "users" ("id");
-- create index "waste_category" to table: "wastes"
CREATE INDEX "waste_category" ON "wastes" ("category");
