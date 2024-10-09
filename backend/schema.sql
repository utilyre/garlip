CREATE TABLE "accounts" (
    "id" SERIAL PRIMARY KEY,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,

    "username" VARCHAR(50) NOT NULL UNIQUE,
    "password" BYTEA NOT NULL,
    "fullname" VARCHAR(100),
    "verified" BOOLEAN NOT NULL DEFAULT false,
    "bio" TEXT
);

CREATE TABLE "forms" (
    "id" SERIAL PRIMARY KEY,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,

    "owner_id" INTEGER REFERENCES "accounts" ON DELETE SET NULL,

    "title" VARCHAR(100) NOT NULL,
    "description" TEXT
);

CREATE TABLE "questions" (
    "id" SERIAL PRIMARY KEY,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,

    "form_id" INTEGER REFERENCES "forms" ON DELETE SET NULL,

    "stem" TEXT NOT NULL
);

CREATE TABLE "topics" (
    "id" SERIAL PRIMARY KEY,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,

    "name" VARCHAR(50) NOT NULL
);

CREATE TABLE "question_topics" (
    "question_id" INTEGER REFERENCES "questions" ON DELETE CASCADE,
    "topic_id" INTEGER REFERENCES "topics" ON DELETE CASCADE,
    PRIMARY KEY ("question_id", "topic_id"),

    "created_at" TIMESTAMP NOT NULL
);

CREATE TABLE "submissions" (
    "id" SERIAL PRIMARY KEY,
    "created_at" TIMESTAMP NOT NULL,

    "respondent_id" INTEGER REFERENCES "accounts" ON DELETE SET NULL,
    "form_id" INTEGER NOT NULL REFERENCES "forms" ON DELETE CASCADE,

    "note" TEXT
);

CREATE TABLE "options" (
    "id" SERIAL PRIMARY KEY,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,

    "question_id" INTEGER NOT NULL REFERENCES "questions" ON DELETE CASCADE,

    "description" TEXT NOT NULL,
    "correct" BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE "answers" (
    "id" SERIAL PRIMARY KEY,

    "submission_id" INTEGER NOT NULL REFERENCES "submissions" ON DELETE CASCADE,
    "question_id" INTEGER NOT NULL REFERENCES "questions" ON DELETE CASCADE,
    "option_id" INTEGER NOT NULL REFERENCES "options" ON DELETE CASCADE,

    UNIQUE ("submission_id", "question_id")
);
