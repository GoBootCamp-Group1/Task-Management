-- Create sequences
CREATE SEQUENCE users_id_seq;
CREATE SEQUENCE boards_id_seq;
CREATE SEQUENCE roles_id_seq;
CREATE SEQUENCE board_users_id_seq;
CREATE SEQUENCE columns_id_seq;
CREATE SEQUENCE tasks_id_seq;

CREATE TABLE "users" (
  "id" bigint PRIMARY KEY DEFAULT nextval('users_id_seq'),
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "name" varchar,
  "email" varchar,
  "password" varchar
);

CREATE TABLE "notifications" (
  "id" uuid PRIMARY KEY ,
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "read_at" timestamp,
  "user_id" bigint,
  "type" varchar,
  "data" json
);

CREATE TABLE "boards" (
  "id" bigint PRIMARY KEY DEFAULT nextval('boards_id_seq'),
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "created_by" bigint,
  "name" varchar,
  "is_private" bool
);

CREATE TABLE "roles" (
  "id" bigint PRIMARY KEY DEFAULT nextval('roles_id_seq'),
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "name" varchar,
  "weight" int
);

CREATE TABLE "board_users" (
  "id" bigint PRIMARY KEY DEFAULT nextval('board_users_id_seq'),
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "user_id" bigint,
  "role_id" bigint,
  "board_id" bigint
);

CREATE TABLE "columns" (
  "id" bigint PRIMARY KEY DEFAULT nextval('columns_id_seq'),
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "board_id" bigint,
  "name" varchar,
  "order_position" int
);

CREATE TABLE "tasks" (
  "id" bigint PRIMARY KEY DEFAULT nextval('tasks_id_seq'),
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "created_by" bigint,
  "board_id" bigint,
  "parent_id" bigint,
  "assignee_id" bigint,
  "name" varchar,
  "description" text,
  "start_datetime" timestamp,
  "end_datetime" timestamp,
  "story_point" int,
  "additional" json,
  "column_id" bigint,
  "order_position" int
);

CREATE TABLE "task_dependencies" (
  "task_id" bigint,
  "dependent_task_id" bigint
);

CREATE TABLE "task_comments" (
  "id" uuid,
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "user_id" bigint,
  "task_id" bigint,
  "comment" text
);

ALTER TABLE "notifications" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "boards" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "board_users" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "board_users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "board_users" ADD FOREIGN KEY ("board_id") REFERENCES "boards" ("id");

ALTER TABLE "columns" ADD FOREIGN KEY ("board_id") REFERENCES "boards" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("board_id") REFERENCES "boards" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("parent_id") REFERENCES "tasks" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("assignee_id") REFERENCES "users" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("column_id") REFERENCES "columns" ("id");

ALTER TABLE "task_dependencies" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("id");

ALTER TABLE "task_dependencies" ADD FOREIGN KEY ("dependent_task_id") REFERENCES "tasks" ("id");

ALTER TABLE "task_comments" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("id");

ALTER TABLE "task_comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
