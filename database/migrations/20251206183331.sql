-- Create "documents" table
CREATE TABLE "public"."documents" (
  "id" serial NOT NULL,
  "name" text NOT NULL,
  "description" text NOT NULL,
  "first_area" text NOT NULL,
  "second_area" text NULL,
  "type" text NOT NULL,
  "link" text NULL,
  PRIMARY KEY ("id")
);
