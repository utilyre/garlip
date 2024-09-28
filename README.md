# Garlip

Garlip is a question bank capable of creating forms and analyzing answers
according to respondents, questions, and topics.

## Development

### Prerequisites

- [Docker Engine](https://docs.docker.com/engine)

- [Go](https://go.dev)

- [Node](https://nodejs.org)

- [pnpm](https://pnpm.io)

  ```bash
  npm install -g pnpm
  ```

- [Atlas](https://atlasgo.io)

  ```bash
  pnpm add -g @ariga/atlas
  ```

- [sqlc](https://sqlc.dev)

  ```bash
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
  ```

### Steps

1. Create a `.env` file at the root of project:

   ```bash
   DB_USER=admin
   DB_PASS=secret

   BE_PORT=8080
   FE_PORT=3000
   ```

2. Configure Atlas at `backend/atlas.hcl`:

   ```hcl
   env "local" {
     src = "file://backend/schema.sql"
     url = "postgres://admin:secret@localhost:5432?search_path=public&sslmode=disable"
     dev = "docker://postgres/16.4-alpine3.20/dev?search_path=public"
   }
   ```

3. Generate ORM code using sqlc:

   ```bash
   (cd backend && sqlc generate)
   ```

4. Spin up all services:

   ```bash
   docker compose -f compose.dev.yaml up
   ```

5. Apply application's schema to the database:

   ```bash
   (cd backend && atlas schema apply --env=local)
   ```

## Stack

### Back-End

- Go
- PostgreSQL
- Atlas
- sqlc
- Chi

### Front-End

- Next
- Tailwind
- shadcn/ui
- Recharts

## Features

- [ ] auth
- [ ] creating forms
- [ ] sharing forms
- [ ] publising/closing forms
- [ ] bar charts
- [ ] dashboard
  - [ ] created forms (from newest to oldest)
  - [ ] participated forms (from newest to oldest)
- [ ] analytics
  - [ ] How many participents answered correctly to a specific topic in a specific timeline?
  - [ ] How many participents answered to each choice of a question?

## Ideas

- Use AI to score descriptive questions
