# stdapp

A comprehensive Go framework for building production-ready full-stack applications with GraphQL APIs, PostgreSQL, Vue.js frontends, and Docker support.

## Features

- **GraphQL API** - Built-in GraphQL server with WebSocket subscriptions via [stdgraph](https://go.ddollar.dev/stdgraph)
- **Multi-Domain Architecture** - Support for multi-tenant applications with isolated database schemas
- **Database Migrations** - Automatic PostgreSQL migrations with [Bun ORM](https://bun.uptrace.dev)
- **Vue.js Integration** - Embedded SPA serving with hot reload in development
- **Development Mode** - File watching with automatic rebuilds
- **Cron Scheduling** - Docker label-based cron jobs for background tasks
- **CLI Framework** - Rich command-line interface via [stdcli](https://go.ddollar.dev/stdcli)
- **Docker Ready** - Complete Docker and docker-compose configuration
- **Middleware Support** - Extensible HTTP middleware chain
- **Project Scaffolding** - Initialize new projects with complete structure

## Installation

```bash
go install go.ddollar.dev/stdapp/cmd/stdapp@latest
```

## Quick Start

Create a new application:

```bash
stdapp init myapp
cd myapp
```

This generates a complete project structure:

```
myapp/
├── api/
│   ├── models/          # Database models
│   └── resolver/        # GraphQL resolvers
│       ├── mutation.go
│       ├── query.go
│       ├── resolver.go
│       ├── schema.graphql
│       └── subscription.go
├── db/
│   └── migrate/         # SQL migrations
├── web/
│   ├── src/            # Vue.js application
│   └── dist/           # Built assets
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── main.go
└── Makefile
```

Start development:

```bash
make dev
# or
docker-compose up
```

## Usage

### Basic Application

```go
package main

import (
    "embed"
    "io/fs"
    "os"
    "time"

    "yourapp/api/resolver"
    "go.ddollar.dev/stdapp"
)

//go:embed db/migrate/*.sql
var migrations embed.FS

//go:embed all:web/dist
var web embed.FS

func main() {
    sweb, _ := fs.Sub(web, "web/dist")

    opts := stdapp.Options{
        Database:     os.Getenv("DATABASE_URL"),
        Domains:      []string{"public"},  // Database schemas
        Migrations:   migrations,
        Name:         "myapp",
        Resolver:     resolver.New,
        Web:          sweb,
        WriteTimeout: 5 * time.Minute,
    }

    app, _ := stdapp.New(opts)
    os.Exit(app.Run(os.Args[1:]))
}
```

### GraphQL Resolver

```go
package resolver

import (
    _ "embed"
    "github.com/uptrace/bun"
    "go.ddollar.dev/stdapp"
)

//go:embed schema.graphql
var schema string

type Resolver struct {
    db *bun.DB
}

func New(db *bun.DB, domain string) (stdapp.Resolver, error) {
    return &Resolver{db: db}, nil
}

func (r *Resolver) Schema() string { return schema }
func (r *Resolver) Query() any { return &Query{r: r} }
func (r *Resolver) Mutation() any { return &Mutation{r: r} }
func (r *Resolver) Subscription() any { return &Subscription{r: r} }
```

### Custom Router

Add REST endpoints alongside GraphQL:

```go
opts := stdapp.Options{
    // ... other options
    Router: func(r *stdapi.Router) error {
        r.Route("GET", "/health", func(c *stdapi.Context) error {
            return c.RenderJSON(map[string]string{"status": "ok"})
        })
        return nil
    },
}
```

### Middleware

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Authentication logic
        next.ServeHTTP(w, r)
    })
}

opts := stdapp.Options{
    // ... other options
    Middleware: []stdapp.Middleware{authMiddleware},
}
```

### Multi-Domain Setup

Support multiple tenants with isolated database schemas:

```go
opts := stdapp.Options{
    Database: os.Getenv("DATABASE_URL"),
    Domains:  []string{"tenant1", "tenant2", "shared"},
    Resolver: resolver.New,
    // ... other options
}
```

Each domain gets its own GraphQL endpoint:
- `/api/tenant1` → PostgreSQL schema `tenant1`
- `/api/tenant2` → PostgreSQL schema `tenant2`
- `/api/shared` → PostgreSQL schema `shared`

## CLI Commands

The framework provides a complete CLI for managing your application:

```bash
# Start the API server
myapp api [--development] [--watch=go,graphql] [--port=8000]

# Run database migrations
myapp migrate [--dry]

# Create a new migration
myapp migration <name> [--dir=db/migrate]

# Start the web server (SPA)
myapp web [--development] [--port=8080]

# Run cron daemon
myapp cron [--development]

# Run arbitrary commands
myapp cmd [--development] <command>

# Database management
myapp pg console [--schema=public]
myapp pg export > backup.sql
myapp pg import < backup.sql
myapp pg reset

# Initialize new project
myapp init <name>
```

### Development Mode

The `--development` flag enables:
- File watching with automatic restarts
- Verbose logging
- Hot reload integration

Customize watched file extensions:

```bash
myapp api --development --watch=go,graphql,sql
```

## Cron Jobs

Schedule background tasks using Docker labels:

```yaml
# docker-compose.yml
services:
  api:
    labels:
      stdapp.cron.backup: "@daily /app/myapp cmd backup"
      stdapp.cron.cleanup: "@every 1h /app/myapp cmd cleanup"
      stdapp.cron.report: "0 9 * * MON /app/myapp cmd weekly-report"
```

Start the cron daemon:

```bash
myapp cron
```

Supported formats:
- Standard cron: `0 9 * * *`
- Descriptors: `@hourly`, `@daily`, `@weekly`, `@monthly`, `@yearly`
- Intervals: `@every 5m`, `@every 1h30m`, `@every 24h`

## Configuration

### Environment Variables

```bash
# Required
DATABASE_URL=postgres://user:pass@localhost/myapp?sslmode=disable

# Optional
PORT=8000
DEVELOPMENT=true
```

### Docker Compose

The generated `docker-compose.yml` includes:
- Application container with hot reload
- PostgreSQL database
- Redis (optional)
- Development overrides via `docker-compose.override.yml`

### Makefile Targets

```bash
make dev        # Start development environment
make build      # Build production binary
make test       # Run tests
make lint       # Run golangci-lint
make vendor     # Vendor dependencies
make migrate    # Run migrations
```

## Architecture

### Request Flow

```
HTTP Request
    ↓
Middleware Chain
    ↓
Router (stdapi)
    ├→ REST Endpoints (custom router)
    └→ GraphQL Handler (stdgraph)
        ↓
    Domain Router (multi-tenant)
        ↓
    Database (Bun ORM with schema)
        ↓
    Resolver (Query/Mutation/Subscription)
```

### Database Schemas

Each domain uses a separate PostgreSQL schema for data isolation:

```sql
-- Migrations run for each domain
CREATE SCHEMA IF NOT EXISTS tenant1;
CREATE SCHEMA IF NOT EXISTS tenant2;

SET search_path TO tenant1;
-- Tables created here
```

### GraphQL Subscriptions

WebSocket subscriptions are supported via [graphql-transport-ws](https://go.ddollar.dev/graphql-transport-ws):

```graphql
subscription {
  messageAdded {
    id
    content
    createdAt
  }
}
```

## Dependencies

stdapp integrates several companion libraries:

- **[stdapi](https://go.ddollar.dev/stdapi)** - HTTP server and routing
- **[stdcli](https://go.ddollar.dev/stdcli)** - CLI framework
- **[stdgraph](https://go.ddollar.dev/stdgraph)** - GraphQL implementation
- **[migrate](https://go.ddollar.dev/migrate)** - Database migrations
- **[logger](https://go.ddollar.dev/logger)** - Structured logging
- **[Bun](https://bun.uptrace.dev)** - SQL ORM
- **[graph-gophers/graphql-go](https://github.com/graph-gophers/graphql-go)** - GraphQL engine

## Examples

### Creating a Model

```go
// api/models/user.go
package models

import "time"

type User struct {
    ID        int64     `bun:",pk,autoincrement"`
    Email     string    `bun:",unique,notnull"`
    Name      string
    CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
```

### GraphQL Schema

```graphql
# api/resolver/schema.graphql
type Query {
  user(id: ID!): User
  users: [User!]!
}

type Mutation {
  createUser(email: String!, name: String!): User!
  deleteUser(id: ID!): Boolean!
}

type Subscription {
  userCreated: User!
}

type User {
  id: ID!
  email: String!
  name: String!
  createdAt: DateTime!
}

scalar DateTime
```

### Query Implementation

```go
// api/resolver/query.go
package resolver

import (
    "context"
    "yourapp/api/models"
    "github.com/graph-gophers/graphql-go"
)

type Query struct {
    r *Resolver
}

func (q *Query) User(ctx context.Context, args struct{ ID graphql.ID }) (*UserResolver, error) {
    var user models.User
    err := q.r.db.NewSelect().
        Model(&user).
        Where("id = ?", args.ID).
        Scan(ctx)
    if err != nil {
        return nil, err
    }
    return &UserResolver{user: user}, nil
}

type UserResolver struct {
    user models.User
}

func (u *UserResolver) ID() graphql.ID {
    return graphql.ID(fmt.Sprint(u.user.ID))
}

func (u *UserResolver) Email() string {
    return u.user.Email
}

func (u *UserResolver) Name() string {
    return u.user.Name
}

func (u *UserResolver) CreatedAt() stdgraph.DateTime {
    return stdgraph.DateTime{Time: u.user.CreatedAt}
}
```

## Testing

```go
func TestAPI(t *testing.T) {
    opts := stdapp.Options{
        Database:   testDatabaseURL(),
        Name:       "test",
        Resolver:   resolver.New,
        Migrations: migrations,
    }

    app, err := stdapp.New(opts)
    require.NoError(t, err)

    // Test your app
}
```

## Production Deployment

### Build

```bash
# Multi-stage Docker build
docker build -t myapp:latest .
```

### Run

```bash
# API server
myapp api --port=8000

# Web server
myapp web --port=8080

# Cron daemon
myapp cron
```

### Kubernetes

The Docker images work seamlessly with Kubernetes or [kip](https://github.com/ddollar/kip):

```yaml
# kip.yml
services:
  api:
    command: api
    environment:
      DATABASE_URL: ${DATABASE_URL}
    port: 8000
    scale:
      count: 3

  cron:
    command: cron
    environment:
      DATABASE_URL: ${DATABASE_URL}
    scale:
      count: 1

resources:
  postgres:
    type: postgres
```

## License

MIT

## Contributing

Issues and pull requests are welcome at the repository.

## Author

David Dollar
