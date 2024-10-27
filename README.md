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
   BE_JWT_SECRET=secret

   FE_PORT=3000
   ```

2. Configure Atlas at `backend/atlas.hcl`:

   ```hcl
   env "local" {
     src = "file://schema.sql"
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
   docker compose up
   ```

5. Apply application's schema to the database:

   ```bash
   (cd backend && atlas schema apply --env=local)
   ```
