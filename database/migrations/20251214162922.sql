-- Modify "user_preferences" table
ALTER TABLE "public"."user_preferences" DROP CONSTRAINT "fk_user_preferences_preferences";
-- Drop "preferences" table
DROP TABLE "public"."preferences";
