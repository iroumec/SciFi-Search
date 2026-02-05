-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "avatar_url" text NULL, ADD CONSTRAINT "users_avatar_url_key" UNIQUE ("avatar_url");
