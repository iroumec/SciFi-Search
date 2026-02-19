-- Modify "document_loads" table
ALTER TABLE "public"."document_loads" ADD CONSTRAINT "fk_document_loads_users" FOREIGN KEY ("loader_id") REFERENCES "public"."users" ("user_id") ON UPDATE CASCADE ON DELETE CASCADE;
