-- Create "users" table
CREATE TABLE "public"."users" (
  "user_id" serial NOT NULL,
  "name" character varying(32) NOT NULL,
  "surname" character varying(32) NOT NULL,
  CONSTRAINT "pk_user" PRIMARY KEY ("user_id")
);
-- Create "historic_searches" table
CREATE TABLE "public"."historic_searches" (
  "historic_search_id" serial NOT NULL,
  "user_id" integer NOT NULL,
  "search_string" text NOT NULL,
  PRIMARY KEY ("historic_search_id"),
  CONSTRAINT "fk_historic_searches_users" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("user_id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create "preferences" table
CREATE TABLE "public"."preferences" (
  "preference" text NOT NULL,
  PRIMARY KEY ("preference")
);
-- Create "user_preferences" table
CREATE TABLE "public"."user_preferences" (
  "user_id" integer NOT NULL,
  "preference" text NOT NULL,
  CONSTRAINT "pk_user_preferences" PRIMARY KEY ("user_id", "preference"),
  CONSTRAINT "fk_user_preferences_preferences" FOREIGN KEY ("preference") REFERENCES "public"."preferences" ("preference") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_user_preferences_users" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("user_id") ON UPDATE CASCADE ON DELETE CASCADE
);
