# stdapi

A lightweight, opinionated HTTP API framework for Go built on top of [gorilla/mux](https://github.com/gorilla/mux).

stdapi provides a cleaner, more ergonomic API for building web applications and HTTP APIs by wrapping standard Go HTTP handlers with a custom Context object and structured error handling.

## Features

- **Context-based request/response handling** - Clean API for accessing request data and rendering responses
- **Structured error handling** - Handlers return errors with HTTP status codes instead of manually writing error responses
- **Built-in middleware support** - Composable middleware with parent router inheritance
- **WebSocket support** - Handle WebSocket connections alongside HTTP endpoints with the same handler interface
- **Session management** - Cookie-based sessions with flash message support
- **Template rendering** - HTML templates with hierarchical layout resolution
- **Automatic request parameter unmarshaling** - Parse query params, form data, and headers into structs
- **TLS/HTTP2 support** - Auto-generated self-signed certificates for development
- **Request ID generation** - Automatic request tracking with structured logging
- **Health checks** - Built-in `/check` endpoint with custom health check support

## Installation

```bash
go get github.com/ddollar/stdapi
```

Requires Go 1.23.0 or later.

## Quick Start

```go
package main

import (
    "github.com/ddollar/stdapi"
)

func main() {
    // Create a new server
    s := stdapi.New("myapp", "localhost")

    // Register routes
    s.Route("GET", "/users", listUsers)
    s.Route("GET", "/users/{id}", getUser)
    s.Route("POST", "/users", createUser)

    // Start the server
    s.Listen("http", ":8080")
}

func listUsers(c *stdapi.Context) error {
    users := []User{
        {ID: 1, Name: "Alice"},
        {ID: 2, Name: "Bob"},
    }
    return c.RenderJSON(users)
}

func getUser(c *stdapi.Context) error {
    id := c.Var("id")
    user := findUser(id)
    if user == nil {
        return stdapi.Errorf(404, "user not found")
    }
    return c.RenderJSON(user)
}

func createUser(c *stdapi.Context) error {
    var user User
    if err := c.BodyJSON(&user); err != nil {
        return stdapi.Errorf(400, "invalid request body")
    }

    if err := c.Required("name", "email"); err != nil {
        return stdapi.Errorf(400, "missing required fields")
    }

    saveUser(&user)
    return c.RenderJSON(user)
}
```

## Core Concepts

### Server

The `Server` is the main entry point. It embeds a `Router` and provides server-level configuration:

```go
s := stdapi.New("namespace", "hostname")

// Optional: Custom health check
s.Check = func(c *stdapi.Context) error {
    // Verify database connection, etc.
    return nil
}

// Optional: Panic recovery
s.Recover = func(err error) {
    log.Printf("panic: %v", err)
}

// Optional: Wrapper middleware
s.Wrapper = func(h http.Handler) http.Handler {
    return someMiddleware(h)
}
```

### HandlerFunc

Unlike standard `http.HandlerFunc`, stdapi handlers receive a `Context` and return an `error`:

```go
func handler(c *stdapi.Context) error {
    // Access request data
    id := c.Var("id")           // Route variables
    name := c.Query("name")     // Query parameters
    value := c.Form("value")    // Form data
    token := c.Header("X-Token") // Headers

    // Render response
    return c.RenderJSON(data)
}
```

If a handler returns an error:
- Errors implementing `stdapi.Error` use the error's status code
- Other errors result in a 500 Internal Server Error
- Error messages are logged and sent to the client

### Context

The `Context` provides access to request/response data and utilities:

#### Request Data
```go
c.Body()              // io.ReadCloser
c.BodyJSON(&v)        // Unmarshal JSON body
c.Form("name")        // Form value
c.Query("name")       // Query parameter
c.Header("name")      // Header value
c.Var("name")         // Route variable
c.Value("name")       // Form or header
c.IP()                // Client IP (respects X-Forwarded-For)
c.Protocol()          // "http" or "https"
c.Ajax()              // true if XMLHttpRequest
c.Request()           // *http.Request
c.Context()           // context.Context
```

#### Responses
```go
c.RenderJSON(v)              // JSON response
c.RenderText("text")         // Plain text
c.RenderOK()                 // "ok\n"
c.RenderTemplate(path, data) // HTML template
c.Redirect(code, url)        // HTTP redirect
c.Response()                 // *stdapi.Response
```

#### Sessions
```go
// Configure session (before starting server)
stdapi.SessionName = "session"
stdapi.SessionSecret = "secret-key"
stdapi.SessionExpiration = 86400 * 30 // 30 days

// Use in handlers
c.SessionSet("user_id", "123")
userID, _ := c.SessionGet("user_id")

c.Flash("success", "User created!")
flashes, _ := c.Flashes()
```

#### Logging
```go
c.Logf("processing user %s", userID)
c.Tag("user_id=%s", userID)
logger := c.Logger()
```

#### Variable Storage
```go
c.Set("key", value)
v := c.Get("key")
```

### Routing

#### Basic Routes
```go
s.Route("GET", "/users", listUsers)
s.Route("POST", "/users", createUser)
s.Route("PUT", "/users/{id}", updateUser)
s.Route("DELETE", "/users/{id}", deleteUser)
s.Route("ANY", "/catch-all", handleAny)
```

#### Path Variables
```go
s.Route("GET", "/users/{id}", getUser)
s.Route("GET", "/posts/{year:[0-9]+}/{month:[0-9]+}", getArchive)

func getUser(c *stdapi.Context) error {
    id := c.Var("id")
    // ...
}
```

#### Redirects
```go
s.Redirect("GET", "/old", 301, "/new")
```

#### Static Files
```go
s.Static("/assets", http.Dir("./public"))
```

#### Subrouters
```go
api := s.Subrouter("/api/v1")
api.Route("GET", "/users", listUsers)  // Handles /api/v1/users

// Inline subrouter
s.SubrouterFunc("/admin", func(admin *stdapi.Router) {
    admin.Route("GET", "/users", adminListUsers)
    admin.Route("POST", "/users", adminCreateUser)
})
```

#### WebSocket Routes
```go
s.Route("SOCKET", "/ws", handleWebSocket)

func handleWebSocket(c *stdapi.Context) error {
    buf := make([]byte, 1024)
    for {
        n, err := c.Read(buf)
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }

        // Process message
        response := process(buf[:n])

        if _, err := c.Write(response); err != nil {
            return err
        }
    }
    return nil
}
```

### Middleware

Middleware wraps handlers to add pre- or post-processing logic:

```go
func authMiddleware(fn stdapi.HandlerFunc) stdapi.HandlerFunc {
    return func(c *stdapi.Context) error {
        token := c.Header("Authorization")
        if token == "" {
            return stdapi.Errorf(401, "unauthorized")
        }

        user, err := validateToken(token)
        if err != nil {
            return stdapi.Errorf(401, "invalid token")
        }

        c.Set("user", user)
        return fn(c)
    }
}

// Apply to entire server
s.Use(authMiddleware)

// Apply to subrouter
api := s.Subrouter("/api")
api.Use(authMiddleware)

// Built-in middleware
s.Use(stdapi.EnsureHTTPS)
```

Child routers inherit parent middleware, executing in order from outermost to innermost.

#### Standard HTTP Middleware
```go
// Use standard http.HandlerFunc as middleware
s.UseHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("X-Custom-Header", "value")
})
```

### Error Handling

Create errors with HTTP status codes:

```go
return stdapi.Errorf(400, "invalid input")
return stdapi.Errorf(401, "unauthorized")
return stdapi.Errorf(404, "user %d not found", userID)
return stdapi.Errorf(500, "internal error")
```

Regular errors are treated as 500 Internal Server Error:

```go
if err := db.Query(); err != nil {
    return err  // Results in 500 response
}
```

### Templates

Configure templates with a filesystem and optional helpers:

```go
stdapi.LoadTemplates(http.Dir("./templates"), func(c *stdapi.Context) template.FuncMap {
    return template.FuncMap{
        "formatDate": func(t time.Time) string {
            return t.Format("2006-01-02")
        },
        "currentUser": func() *User {
            return c.Get("user").(*User)
        },
    }
})
```

Render templates in handlers:

```go
func showUser(c *stdapi.Context) error {
    user := findUser(c.Var("id"))
    return c.RenderTemplate("users/show", user)
}
```

Templates support hierarchical layouts. Rendering `admin/users/list` loads:
1. `layout.tmpl` (root)
2. `admin/layout.tmpl`
3. `admin/users/layout.tmpl`
4. `admin/users/list.tmpl`

Templates should define a `main` block:

```html
{{ define "main" }}
<h1>{{ .Title }}</h1>
{{ end }}
```

### Request Parameter Unmarshaling

Automatically parse request parameters into structs:

```go
type ListOptions struct {
    Page   *int    `query:"page" default:"1"`
    Limit  *int    `query:"limit" default:"20"`
    Sort   *string `query:"sort" default:"created_at"`
    APIKey *string `header:"X-API-Key"`
}

func listUsers(c *stdapi.Context) error {
    var opts ListOptions
    if err := stdapi.UnmarshalOptions(c.Request(), &opts); err != nil {
        return stdapi.Errorf(400, "invalid parameters")
    }

    users := db.FindUsers(*opts.Page, *opts.Limit, *opts.Sort)
    return c.RenderJSON(users)
}
```

Supported tags:
- `param` - Form/POST data
- `query` - URL query parameters
- `header` - HTTP headers
- `default` - Default value if missing

Supported types:
- `*bool`, `*int`, `*int64`, `*string`
- `*time.Duration`, `*time.Time`
- `[]string` (comma-separated)
- `map[string]string` (URL query encoded)

### Protocols

Start the server with different protocols:

```go
// Plain HTTP
s.Listen("http", ":8080")

// HTTPS with auto-generated self-signed certificate
s.Listen("https", ":443")
s.Listen("tls", ":443")

// HTTP/2
s.Listen("h2", ":443")
```

For production, provide your own certificates via `Wrapper` or reverse proxy.

### Graceful Shutdown

```go
server := stdapi.New("myapp", "localhost")
server.Route("GET", "/", handler)

go func() {
    if err := server.Listen("http", ":8080"); err != nil {
        log.Fatal(err)
    }
}()

// Wait for interrupt signal
quit := make(chan os.Signal, 1)
signal.Notify(quit, os.Interrupt)
<-quit

// Graceful shutdown with 5 second timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := server.Shutdown(ctx); err != nil {
    log.Fatal(err)
}
```

## Advanced Usage

### Custom Health Checks

```go
s.Check = func(c *stdapi.Context) error {
    if err := db.Ping(); err != nil {
        return err
    }
    if err := redis.Ping(); err != nil {
        return err
    }
    return nil
}

// GET /check returns 500 if health check fails
```

### Custom 404 Handler

```go
s.HandleNotFound(func(c *stdapi.Context) error {
    return c.RenderJSON(map[string]string{
        "error": "not found",
        "path":  c.Request().URL.Path,
    })
})
```

### Panic Recovery

```go
s.Recover = func(err error) {
    log.Printf("panic recovered: %+v", err)
    // Send to error tracking service
}
```

### Custom Route Matching

```go
// Match requests based on custom logic
apiV2 := s.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
    return r.Header.Get("API-Version") == "2"
})

apiV2.Route("GET", "/users", listUsersV2)
```

### Accessing Underlying Types

```go
// Get *http.Request
req := c.Request()

// Get *http.Response
resp := c.Response()

// Get WebSocket connection (or nil)
ws := c.Websocket()

// Get route from returned Route
route := s.Route("GET", "/users/{id}", getUser)
route.Host("api.example.com")
route.Schemes("https")
```

## Examples

### RESTful API

```go
package main

import (
    "github.com/ddollar/stdapi"
    "log"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var users = []User{
    {ID: 1, Name: "Alice", Email: "alice@example.com"},
    {ID: 2, Name: "Bob", Email: "bob@example.com"},
}

func main() {
    s := stdapi.New("api", "localhost")

    // Routes
    s.Route("GET", "/users", listUsers)
    s.Route("GET", "/users/{id}", getUser)
    s.Route("POST", "/users", createUser)
    s.Route("PUT", "/users/{id}", updateUser)
    s.Route("DELETE", "/users/{id}", deleteUser)

    log.Fatal(s.Listen("http", ":8080"))
}

func listUsers(c *stdapi.Context) error {
    return c.RenderJSON(users)
}

func getUser(c *stdapi.Context) error {
    id := c.Var("id")
    for _, user := range users {
        if fmt.Sprintf("%d", user.ID) == id {
            return c.RenderJSON(user)
        }
    }
    return stdapi.Errorf(404, "user not found")
}

func createUser(c *stdapi.Context) error {
    var user User
    if err := c.BodyJSON(&user); err != nil {
        return stdapi.Errorf(400, "invalid JSON")
    }

    user.ID = len(users) + 1
    users = append(users, user)

    return c.RenderJSON(user)
}

func updateUser(c *stdapi.Context) error {
    id := c.Var("id")
    var updated User
    if err := c.BodyJSON(&updated); err != nil {
        return stdapi.Errorf(400, "invalid JSON")
    }

    for i, user := range users {
        if fmt.Sprintf("%d", user.ID) == id {
            users[i] = updated
            users[i].ID = user.ID
            return c.RenderJSON(users[i])
        }
    }

    return stdapi.Errorf(404, "user not found")
}

func deleteUser(c *stdapi.Context) error {
    id := c.Var("id")
    for i, user := range users {
        if fmt.Sprintf("%d", user.ID) == id {
            users = append(users[:i], users[i+1:]...)
            return c.RenderOK()
        }
    }
    return stdapi.Errorf(404, "user not found")
}
```

### Authentication Middleware

```go
func authRequired(fn stdapi.HandlerFunc) stdapi.HandlerFunc {
    return func(c *stdapi.Context) error {
        token := c.Header("Authorization")
        if token == "" {
            return stdapi.Errorf(401, "missing authorization header")
        }

        user, err := validateToken(token)
        if err != nil {
            return stdapi.Errorf(401, "invalid token")
        }

        c.Set("user", user)
        c.Tag("user_id=%d", user.ID)

        return fn(c)
    }
}

// Apply to routes
s.Use(authRequired)

// Or to specific subrouters
api := s.Subrouter("/api")
api.Use(authRequired)
api.Route("GET", "/profile", getProfile)

func getProfile(c *stdapi.Context) error {
    user := c.Get("user").(*User)
    return c.RenderJSON(user)
}
```

### WebSocket Chat

```go
s.Route("SOCKET", "/chat", handleChat)

func handleChat(c *stdapi.Context) error {
    username := c.Query("username")
    if username == "" {
        return stdapi.Errorf(400, "username required")
    }

    c.Logf("user %s connected", username)

    buf := make([]byte, 1024)
    for {
        n, err := c.Read(buf)
        if err == io.EOF {
            c.Logf("user %s disconnected", username)
            break
        }
        if err != nil {
            return err
        }

        message := fmt.Sprintf("%s: %s", username, string(buf[:n]))
        broadcast(message)

        if _, err := c.Write([]byte(message)); err != nil {
            return err
        }
    }

    return nil
}
```

## Documentation

Full API documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/ddollar/stdapi).

Generate local documentation:

```bash
godoc -http=:6060
```

Then visit http://localhost:6060/pkg/github.com/ddollar/stdapi/

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Dependencies

- [github.com/ddollar/logger](https://github.com/ddollar/logger) - Structured logging
- [github.com/gorilla/mux](https://github.com/gorilla/mux) - HTTP router
- [github.com/gorilla/sessions](https://github.com/gorilla/sessions) - Session management
- [github.com/gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket support
- [github.com/pkg/errors](https://github.com/pkg/errors) - Error handling with stack traces
- [github.com/sebest/xff](https://github.com/sebest/xff) - X-Forwarded-For header parsing
