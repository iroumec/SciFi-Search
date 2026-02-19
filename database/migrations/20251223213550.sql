-- Modify "document_loads" table
ALTER TABLE "public"."document_loads" DROP CONSTRAINT "fk_document_loads_documents", ADD CONSTRAINT "fk_document_loads_documents" FOREIGN KEY ("document_id") REFERENCES "public"."documents" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
