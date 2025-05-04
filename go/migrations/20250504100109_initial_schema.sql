-- Create "user" table
CREATE TABLE "user" (
  "id" uuid NOT NULL,
  "email" character varying(255) NOT NULL,
  "password" character varying(255) NOT NULL,
  "name" character varying(255) NOT NULL,
  "role" character varying(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "user_email_key" UNIQUE ("email"),
  CONSTRAINT "user_role_check" CHECK ((role)::text = ANY ((ARRAY['admin'::character varying, 'user'::character varying])::text[]))
);
