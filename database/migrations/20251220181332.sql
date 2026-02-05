-- Modify "documents" table
ALTER TABLE "public"."documents" ALTER COLUMN "description" DROP NOT NULL, ADD COLUMN "based_on" text NULL, ADD COLUMN "grantor" text NULL, ADD COLUMN "deadline" text NOT NULL;
