package stdapi

// EnsureHTTPS is middleware that redirects HTTP requests to HTTPS.
//
// This checks the X-Forwarded-Proto header and issues a 301 permanent redirect
// if the protocol is http. Use this when running behind a load balancer or proxy.
func EnsureHTTPS(fn HandlerFunc) HandlerFunc {
	return func(c *Context) error {
		if c.Request().Header.Get("X-Forwarded-Proto") == "http" {
			u := *(c.Request().URL)
			u.Host = c.Request().Host
			u.Scheme = "https"
			return c.Redirect(301, u.String())
		}
		return fn(c)
	}
}
