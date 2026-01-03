CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" decimal NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "valid_balance" CHECK (balance>=0)
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigserial NOT NULL
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "from" bigserial NOT NULL,
  "to" bigserial NOT NULL,
  "amount" decimal NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "positive_amount" CHECK (amount>0)
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transactions" ("from");

CREATE INDEX ON "transactions" ("to");

CREATE INDEX ON "transactions" ("from", "to");

COMMENT ON COLUMN "accounts"."balance" IS 'decimal类型，精确小数';

COMMENT ON COLUMN "accounts"."currency" IS '货币类型';

COMMENT ON TABLE "entries" IS '取/存款记录';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("from") REFERENCES "accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("to") REFERENCES "accounts" ("id");
