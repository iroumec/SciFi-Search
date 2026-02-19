-- Modify "historic_searches" table
ALTER TABLE "public"."historic_searches" ADD COLUMN "search_date" date NOT NULL DEFAULT now();
