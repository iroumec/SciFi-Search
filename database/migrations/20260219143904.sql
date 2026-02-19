-- Modify "documents" table
ALTER TABLE "public"."documents" ADD CONSTRAINT "fk_documents_users" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("user_id") ON UPDATE CASCADE ON DELETE CASCADE;
