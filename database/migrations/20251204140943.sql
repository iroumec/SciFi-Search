-- Modify "historic_searches" table
ALTER TABLE "public"."historic_searches" DROP COLUMN "earch_datetime", ADD COLUMN "search_datetime" timestamptz NOT NULL DEFAULT now();
