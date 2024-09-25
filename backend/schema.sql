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

    "owner_id" INTEGER REFERENCES "accounts",

    "title" VARCHAR(100) NOT NULL,
    "description" TEXT
);

CREATE TABLE "questions" (
    "id" SERIAL PRIMARY KEY,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,

    "form_id" INTEGER REFERENCES "forms",

    "stem" TEXT NOT NULL
);

CREATE TABLE "topics" (
    "id" SERIAL PRIMARY KEY,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,

    "name" VARCHAR(50) NOT NULL
);

CREATE TABLE "question_topics" (
    "question_id" INTEGER REFERENCES "questions",
    "topic_id" INTEGER REFERENCES "topics",
    PRIMARY KEY ("question_id", "topic_id"),

    "created_at" TIMESTAMP NOT NULL
);

CREATE TABLE "options" (
    "id" SERIAL PRIMARY KEY,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,

    "question_id" INTEGER REFERENCES "questions" NOT NULL,

    "description" TEXT NOT NULL,
    "correct" BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE "optional_answers" (
    "id" SERIAL PRIMARY KEY,
    "submitted_at" TIMESTAMP NOT NULL,

    "participent_id" INTEGER REFERENCES "accounts",
    "question_id" INTEGER REFERENCES "questions" NOT NULL UNIQUE,
    "option_id" INTEGER REFERENCES "options" NOT NULL
);
