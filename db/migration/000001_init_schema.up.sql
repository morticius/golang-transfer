CREATE TABLE "currencies" (
  "id" bigserial PRIMARY KEY,
  "code" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "currency_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "time_placed" timestamp NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "transfers" ADD FOREIGN KEY ("currency_id") REFERENCES "currencies" ("id");

CREATE INDEX ON "currencies" ("code");

CREATE INDEX ON "transfers" ("user_id");

CREATE INDEX ON "transfers" ("currency_id");

CREATE INDEX ON "transfers" ("user_id", "currency_id");

COMMENT ON COLUMN "currencies"."code" IS 'ISO 4217';

COMMENT ON COLUMN "transfers"."amount" IS 'can be negative(can be negative(withdrawal) or positive(deposit)';