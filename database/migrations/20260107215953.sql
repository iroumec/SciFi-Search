-- Modify "documents" table
ALTER TABLE "public"."documents" ADD COLUMN "currency" text NOT NULL, ADD COLUMN "amount" text NOT NULL;
