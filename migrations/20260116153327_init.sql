-- Create "roles" table
CREATE TABLE "roles" (
  "id" uuid NOT NULL,
  "name" character varying(32) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "roles_name_key" UNIQUE ("name")
);
-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL,
  "role_id" uuid NULL,
  "first_name" character varying(32) NOT NULL,
  "last_name" character varying(32) NOT NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  "phone" text NOT NULL DEFAULT '',
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "users_email_key" UNIQUE ("email"),
  CONSTRAINT "users_phone_key" UNIQUE ("phone"),
  CONSTRAINT "role_id" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "posts" table
CREATE TABLE "posts" (
  "id" uuid NOT NULL,
  "title" character varying(255) NOT NULL,
  "content" text NULL,
  "published" boolean NOT NULL DEFAULT false,
  "user_id" uuid NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
