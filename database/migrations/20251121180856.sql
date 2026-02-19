-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "auth_id" text NOT NULL, ADD CONSTRAINT "users_auth_id_key" UNIQUE ("auth_id");
