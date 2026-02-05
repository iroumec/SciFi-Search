-- Modify "documents" table
ALTER TABLE "public"."documents" DROP COLUMN "first_area", DROP COLUMN "second_area", ADD COLUMN "main_area" text NOT NULL, ADD COLUMN "secondary_area" text NULL;
