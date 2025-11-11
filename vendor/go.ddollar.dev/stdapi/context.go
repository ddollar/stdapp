package stdapi

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"go.ddollar.dev/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sebest/xff"
)

var (
	// SessionExpiration is the session cookie max age in seconds (default: 30 days).
	SessionExpiration = 86400 * 30

	// SessionName is the name of the session cookie.
	// Must be set before using session features.
	SessionName = ""

	// SessionSecret is the secret key used to sign session cookies.
	// Must be set before using session features.
	SessionSecret = ""
)

// Context wraps an HTTP request and response with convenience methods.
//
// Context provides a clean API for accessing request data, rendering responses,
// managing sessions, and logging. It can transparently handle both HTTP and
// WebSocket connections.
type Context struct {
	context  context.Context
	id       string
	logger   *logger.Logger
	name     string
	request  *http.Request
	response *Response
	rvars    map[string]string
	session  sessions.Store
	vars     map[string]interface{}
	ws       *websocket.Conn
}

// Flash represents a one-time notification message stored in the session.
//
// Flash messages are automatically deleted after being retrieved, making them
// ideal for displaying success/error messages after redirects.
type Flash struct {
	// Kind is the type of flash message (e.g., "success", "error", "warning").
	Kind string

	// Message is the flash message content.
	Message string
}

func init() {
	gob.Register(Flash{})
}

// NewContext creates a new Context wrapping the given response writer and request.
//
// This is typically called internally by the router, but can be used directly
// for testing or custom request handling.
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	s := sessions.NewCookieStore([]byte(SessionSecret))
	s.Options.MaxAge = SessionExpiration
	s.Options.SameSite = http.SameSiteLaxMode

	return &Context{
		context:  r.Context(),
		logger:   logger.New(""),
		request:  r,
		response: &Response{ResponseWriter: w},
		rvars:    map[string]string{},
		session:  s,
		vars:     map[string]interface{}{},
	}
}

// Ajax returns true if the request was made via XMLHttpRequest.
//
// This checks for the X-Requested-With: XMLHttpRequest header.
func (c *Context) Ajax() bool {
	return c.request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// Body returns the request body as a ReadCloser.
func (c *Context) Body() io.ReadCloser {
	return c.request.Body
}

// BodyJSON reads and unmarshals the request body as JSON into v.
func (c *Context) BodyJSON(v interface{}) error {
	data, err := ioutil.ReadAll(c.Body())
	if err != nil {
		return errors.WithStack(err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// func (c *Context) Close() error {
//   return nil
// }

// Context returns the request's context.Context.
func (c *Context) Context() context.Context {
	return c.context
}

// Flash adds a flash message to the session.
//
// The kind parameter typically indicates the message type (e.g., "success", "error").
// Flash messages are automatically deleted when retrieved via Flashes().
func (c *Context) Flash(kind, message string) error {
	s, err := c.session.Get(c.request, SessionName)
	if err != nil {
		return errors.WithStack(err)
	}

	s.AddFlash(Flash{Kind: kind, Message: message})

	return s.Save(c.request, c.response)
}

// Flashes retrieves and clears all flash messages from the session.
func (c *Context) Flashes() ([]Flash, error) {
	s, err := c.session.Get(c.request, SessionName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fs := []Flash{}

	for _, f := range s.Flashes() {
		if ff, ok := f.(Flash); ok {
			fs = append(fs, ff)
		}
	}

	if err := s.Save(c.request, c.response); err != nil {
		return nil, errors.WithStack(err)
	}

	return fs, nil
}

// Form returns the named form value from POST or PUT body parameters.
func (c *Context) Form(name string) string {
	return c.request.FormValue(name)
}

// Get retrieves a value from the context's variable store.
//
// Values are stored via Set() and persist for the lifetime of the request.
func (c *Context) Get(name string) interface{} {
	v, ok := c.vars[name]
	if !ok {
		return nil
	}

	return v
}

// Header returns the named request header value.
func (c *Context) Header(name string) string {
	return c.request.Header.Get(name)
}

// IP returns the client's IP address.
//
// This respects X-Forwarded-For headers when behind a proxy.
func (c *Context) IP() string {
	return strings.Split(xff.GetRemoteAddr(c.Request()), ":")[0]
}

// Logger returns the structured logger for this request.
//
// The logger is automatically tagged with the request ID and handler name.
func (c *Context) Logger() *logger.Logger {
	return c.logger
}

// Logf logs a formatted message using the request's logger.
func (c *Context) Logf(format string, args ...interface{}) {
	c.logger.Logf(format, args...)
}

// Name returns the handler function name for this request.
func (c *Context) Name() string {
	return c.name
}

// Protocol returns the request protocol (http or https).
//
// This respects the X-Forwarded-Proto header when behind a proxy, defaulting to https.
func (c *Context) Protocol() string {
	if h := c.Header("X-Forwarded-Proto"); h != "" {
		return h
	}

	return "https"
}

// Query returns the named URL query parameter value.
func (c *Context) Query(name string) string {
	return c.request.URL.Query().Get(name)
}

// Read reads data from the request body or WebSocket connection.
//
// For HTTP requests, this reads from the request body.
// For WebSocket connections, this reads the next WebSocket message.
// Returns io.EOF when the connection is closed or a binary message is received.
func (c *Context) Read(data []byte) (int, error) {
	if c.ws == nil {
		return c.Body().Read(data)
	}

	t, r, err := c.ws.NextReader()
	if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		return 0, io.EOF
	}
	if err != nil {
		return 0, errors.WithStack(err)
	}

	switch t {
	case websocket.TextMessage:
		return r.Read(data)
	case websocket.BinaryMessage:
		return 0, io.EOF
	default:
		return 0, errors.WithStack(fmt.Errorf("unknown message type: %d\n", t))
	}
}

// Redirect sends an HTTP redirect response to the target URL with the given status code.
func (c *Context) Redirect(code int, target string) error {
	http.Redirect(c.response, c.request, target, code)
	return nil
}

// RenderJSON renders v as indented JSON with Content-Type: application/json.
func (c *Context) RenderJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return errors.WithStack(err)
	}

	c.response.Header().Set("Content-Type", "application/json")

	if _, err := c.response.Write(data); err != nil {
		return errors.WithStack(err)
	}

	if _, err := c.response.Write([]byte{10}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// RenderOK renders a simple "ok\n" response with status 200.
func (c *Context) RenderOK() error {
	fmt.Fprintf(c.response, "ok\n")
	return nil
}

// RenderTemplate renders an HTML template with the given parameters.
//
// See LoadTemplates for information about template configuration and layout resolution.
func (c *Context) RenderTemplate(path string, params interface{}) error {
	return RenderTemplate(c, path, params)
}

// RenderTemplatePart renders a specific named template block.
func (c *Context) RenderTemplatePart(path, part string, params interface{}) error {
	return RenderTemplatePart(c, path, part, params)
}

// RenderText writes plain text to the response.
func (c *Context) RenderText(t string) error {
	_, err := c.response.Write([]byte(t))
	return errors.WithStack(err)
}

// Request returns the underlying *http.Request.
func (c *Context) Request() *http.Request {
	return c.request
}

// Required validates that all named form parameters are present and non-empty.
//
// Returns an error listing any missing parameters.
func (c *Context) Required(names ...string) error {
	missing := []string{}

	for _, n := range names {
		if c.Form(n) == "" {
			missing = append(missing, n)
		}
	}

	if len(missing) > 0 {
		return errors.WithStack(fmt.Errorf("parameter required: %s", strings.Join(missing, ", ")))
	}

	return nil
}

// Response returns the wrapped response writer.
func (c *Context) Response() *Response {
	return c.response
}

// SessionGet retrieves a string value from the session.
//
// Returns an empty string if the key doesn't exist.
// Requires SessionName and SessionSecret to be configured.
func (c *Context) SessionGet(name string) (string, error) {
	if SessionName == "" {
		return "", fmt.Errorf("no session name set")
	}

	if SessionSecret == "" {
		return "", fmt.Errorf("no session secret set")
	}

	s, _ := c.session.Get(c.request, SessionName)

	vi, ok := s.Values[name]
	if !ok {
		return "", nil
	}

	vs, ok := vi.(string)
	if !ok {
		return "", errors.WithStack(fmt.Errorf("session value is not string"))
	}

	return vs, nil
}

// SessionSet stores a string value in the session.
//
// Requires SessionName and SessionSecret to be configured.
func (c *Context) SessionSet(name, value string) error {
	if SessionName == "" {
		return fmt.Errorf("no session name set")
	}

	if SessionSecret == "" {
		return fmt.Errorf("no session secret set")
	}

	s, _ := c.session.Get(c.request, SessionName)

	s.Values[name] = value

	return s.Save(c.request, c.response)
}

// Set stores a value in the context's variable store.
//
// Values persist for the lifetime of the request and can be retrieved via Get().
func (c *Context) Set(name string, value interface{}) {
	c.vars[name] = value
}

// Tag appends structured tags to the request logger.
//
// Example: c.Tag("user_id=%d", userID)
func (c *Context) Tag(format string, args ...interface{}) {
	c.logger = c.logger.Append(format, args...)
}

// SetVar sets a route variable value.
//
// This can be used by middleware to override route variables.
func (c *Context) SetVar(name, value string) {
	c.rvars[name] = value
}

// Value returns the named parameter from form data or headers.
//
// This checks form parameters first, then falls back to headers.
func (c *Context) Value(name string) string {
	if v := c.Form(name); v != "" {
		return v
	}

	if v := c.Header(name); v != "" {
		return v
	}

	return ""
}

// Var returns the named route variable value from the URL path.
//
// For a route like "/users/{id}", Var("id") returns the captured value.
func (c *Context) Var(name string) string {
	if v, ok := c.rvars[name]; ok {
		return v
	}
	return mux.Vars(c.request)[name]
}

// Websocket returns the underlying WebSocket connection, or nil for HTTP requests.
func (c *Context) Websocket() *websocket.Conn {
	return c.ws
}

// Write writes data to the response or WebSocket connection.
//
// For HTTP requests, this writes to the response body.
// For WebSocket connections, this sends a text message.
func (c *Context) Write(data []byte) (int, error) {
	if c.ws == nil {
		return c.response.Write(data)
	}

	w, err := c.ws.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	defer w.Close()

	return w.Write(data)
}
