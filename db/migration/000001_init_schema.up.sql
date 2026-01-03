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
  "account_id" bigserial NOT NULL,
  -- transfer_id is null when type is debit or credit
  "transfer_id" bigserial,
  "amount" decimal NOT NULL,
  -- entry_type: debit/credit, or NULL for a transfer
  "entry_type" varchar(32),
  "currency" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "optional_transfer" CHECK (transfer_id IS NOT NULL OR entry_type IS NOT NULL)
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from" bigserial NOT NULL,
  "to" bigserial NOT NULL,
  "amount" decimal NOT NULL,
  "instruction" varchar(500),
  "created_at" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "positive_amount" CHECK (amount>0)
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from");

CREATE INDEX ON "transfers" ("to");

CREATE INDEX ON "transfers" ("from", "to");

COMMENT ON COLUMN "accounts"."balance" IS 'decimal类型，精确小数';

COMMENT ON COLUMN "accounts"."currency" IS '货币类型';

COMMENT ON TABLE "entries" IS '取/存款记录';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "entries" ADD FOREIGN KEY ("transfer_id") REFERENCES "transfers" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to") REFERENCES "accounts" ("id");
