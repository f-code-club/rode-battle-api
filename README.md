# R.ODE Battle API

## Tools

- [**Go**](https://golang.org/dl/): Version 1.26.2 or later.
- [**PostgreSQL**](https://www.postgresql.org/download/).
- **migrate**: Used for managing migrations.
  ```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  ```
- **sqlc**: Used for generating type-safe Go code from SQL.
  ```bash
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
  ```
- **Air**: A live-reloading tool for Go applications.
  ```bash
  go install github.com/cosmtrek/air@latest
  ```

## Getting Started

### Clone the repository

```bash
git clone https://github.com/f-code-club/rode-battle-api.git
cd rode-battle-api
```

### Setup Environment Variables

```
DATABASE_URL="postgres://user:password@localhost:5432/rode?sslmode=disable"
```

### Setup Database

```bash
createdb rode
migrate -database $DATABASE_URL -path migration up
```

### Run the Application

```bash
air
```
