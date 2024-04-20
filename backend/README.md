# Backend 💻

## Run locally

### Set environment variables (.env) - example:
```bash
SERVER_PORT=8080
SERVER_TIMEOUT_READ=3s
SERVER_TIMEOUT_WRITE=5s
SERVER_TIMEOUT_IDLE=5s
SERVER_DEBUG=true

POSTGRES_HOST=db # If you are running it from Docker Compose 
POSTGRES_PORT=5432
POSTGRES_DB=users_service
POSTGRES_USER=admin
POSTGRES_PASSWORD=1234
POSTGRES_DEBUG=true
```

### Running API

```bash
# locally
go run /backend/cmd/api/main.go

# or using docker-compose
docker-compose up
```

### Running migrations

```bash
# Creates new migration file with the current timestamp
go run /backend/cmd/migrate/main.go create NAME [sql|go] 

# Migrate the DB to the most recent version available
go run /backend/cmd/migrate/main.go up

# Roll back the version by 1
go run /backend/cmd/migrate/main.go down

# More commands can be found in Help
go run /backend/cmd/migrate/main.go -h
```

### Genearting docs

```bash
swag init -g ./cmd/api/main.go
```

## Folder structure
```shell
myapp
├── cmd
│  ├── api
│  │  └── main.go
│  └── migrate
│     └── main.go
│
├── api
│  ├── resource
│  │  ├── users
│  │  │  ├── handler.go
│  │  │  ├── model.go
│  │  │  ├── repository.go
│  │  │  └── repository_test.go
│  │  ├── common
│  │  │  └── error
│  │  │     └── error.go
│  │  └── health
│  │     └── handler.go
│  │
│  └── router
│     ├── middleware
│     │  ├── content_type.go
│     │  └── request_logger.go
│     └── router.go
│
├── migrations
│  └── 00001_create_users_table.sql
│
├── config
│  └── config.go
│
├── util
│  ├── logger
│  │  └── logger.go
│  ├── mock
│  │  └── db.go
│  └── validator
│     └── validator.go
│
├── .env
│
├── go.mod
├── go.sum
│
├── docker-compose.yml
└── Dockerfile
```
