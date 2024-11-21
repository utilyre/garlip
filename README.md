# Garlip

Garlip is a question bank capable of creating forms and analyzing answers
with respect to participant, question, and topic.

## Development

### Prerequisites

- [Docker](https://docker.com)

- [Task](https://taskfile.dev)

- [Go](https://go.dev)

- [Node](https://nodejs.org)

- [pnpm](https://pnpm.io)

  ```bash
  npm install -g pnpm
  ```

- [migrate](https://github.com/golang-migrate/migrate)

  ```bash
  go install -tags=postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest
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
   DB_NAME=garlip

   BE_PORT=8080
   JWT_SECRET=secret

   FE_PORT=3000
   ```

2. Generate Go source code from SQL queries:

   ```bash
   task queries
   ```

4. Spin up all services:

   ```bash
   task start
   ```

5. Apply migrations:

   ```bash
   task migrate:up
   ```

### Tips

- Run `task -l` to get a list of all available commands.
