-- Create "document_loads" table
CREATE TABLE "public"."document_loads" (
  "loader_id" integer NOT NULL,
  "document_id" integer NOT NULL,
  CONSTRAINT "pk_document_loads" PRIMARY KEY ("loader_id", "document_id")
);
