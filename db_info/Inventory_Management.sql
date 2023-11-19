CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "category_name" varchar NOT NULL
);

CREATE TABLE "units" (
  "id" bigserial PRIMARY KEY,
  "unit_name" varchar NOT NULL
);

CREATE TABLE "goods" (
  "id" bigserial PRIMARY KEY,
  "category" bigint NOT NULL,
  "model" varchar NOT NULL,
  "unit" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "desc" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "categories" ("category_name");

CREATE INDEX ON "units" ("unit_name");

CREATE INDEX ON "goods" ("category");

CREATE INDEX ON "goods" ("model");

CREATE INDEX ON "goods" ("category", "model");

COMMENT ON COLUMN "goods"."amount" IS 'must be positive and bigger than zero';

ALTER TABLE "goods" ADD FOREIGN KEY ("category") REFERENCES "categories" ("id");

ALTER TABLE "goods" ADD FOREIGN KEY ("unit") REFERENCES "units" ("id");
