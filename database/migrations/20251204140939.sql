-- Modify "historic_searches" table
ALTER TABLE "public"."historic_searches" DROP COLUMN "search_date", ADD COLUMN "earch_datetime" timestamptz NOT NULL DEFAULT now();
