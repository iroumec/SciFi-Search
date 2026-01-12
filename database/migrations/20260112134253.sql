-- Create index "idx_user_preferences_preference_trgm" to table: "user_preferences"
CREATE INDEX "idx_user_preferences_preference_trgm" ON "public"."user_preferences" USING gin ("preference" gin_trgm_ops);
