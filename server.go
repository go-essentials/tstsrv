// =====================================================================================================================
// == LICENSE:       Copyright (c) 2025 Kevin De Coninck
// ==
// ==                Permission is hereby granted, free of charge, to any person
// ==                obtaining a copy of this software and associated documentation
// ==                files (the "Software"), to deal in the Software without
// ==                restriction, including without limitation the rights to use,
// ==                copy, modify, merge, publish, distribute, sublicense, and/or sell
// ==                copies of the Software, and to permit persons to whom the
// ==                Software is furnished to do so, subject to the following
// ==                conditions:
// ==
// ==                The above copyright notice and this permission notice shall be
// ==                included in all copies or substantial portions of the Software.
// ==
// ==                THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// ==                EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// ==                OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// ==                NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// ==                HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// ==                WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// ==                FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// ==                OTHER DEALINGS IN THE SOFTWARE.
// =====================================================================================================================

// Package tstsrv implements a configurable development server, suitable for testing.
package tstsrv

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

// Server wraps Go's built-in "httptest.Server" but provides an API for configuring the responses.
type Server struct {
	httpServer *httptest.Server             // The actual "httptest.Server".
	routes     map[string]RespConfiguration // The map of routes and their configuration.
	lock       sync.Mutex                   // Protect concurrent access to call counts.
}

// RespConfiguration is the configuration for a Server.
type RespConfiguration struct {
	Responses []Response // A sequence of HTTP responses to return.
	callCount int        // Counter to track the number of calls.
}

// Response represents the response to an HTTP request.
type Response struct {
	StatusCode     int    // The HTTP status code to return.
	Body           string // The body to return.
	DropConnection bool   // Drop the connection. This is to simulate that the body can't be read.
}

// New returns a new Server with the given routes.
func New(routes map[string]RespConfiguration) *Server {
	server := &Server{
		routes: routes,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.lock.Lock()
		defer server.lock.Unlock()

		requestUri := rawUrl(r)

		if routeConfig, match := server.routes[requestUri]; !match || routeConfig.callCount >= len(routeConfig.Responses) {
			w.WriteHeader(http.StatusNotImplemented)

			return
		} else {
			response := routeConfig.Responses[routeConfig.callCount]
			routeConfig.callCount++
			server.routes[requestUri] = routeConfig

			w.WriteHeader(response.StatusCode)

			if response.DropConnection {
				conn, _, _ := w.(http.Hijacker).Hijack()
				conn.Close()
			} else {
				response.Body = strings.Replace(response.Body, "$$URI$$", server.httpServer.URL, -1)

				w.Write([]byte(response.Body))
			}
		}
	})

	server.httpServer = httptest.NewServer(handler)

	return server
}

// Close closes the underlying httptest.Server.
func (f *Server) Close() {
	f.httpServer.Close()
}

// URL returns the URL of the server.
func (f *Server) URL() string {
	return f.httpServer.URL
}

// Returns the complete URL (including the query string) of r.
func rawUrl(r *http.Request) string {
	if r.URL.RawQuery != "" {
		return fmt.Sprintf("%s?%s", r.URL.Path, r.URL.RawQuery)
	}

	return r.URL.Path
}
