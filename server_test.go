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
package tstsrv_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/go-essentials/assert"
	"github.com/go-essentials/tstsrv"
)

// UT: Send requests to the server and check the response.
func TestServer(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	// FAKE SETUP.
	srvFake := tstsrv.New(map[string]tstsrv.RespConfiguration{
		"/test?v=10": {
			Responses: []tstsrv.Response{
				{StatusCode: http.StatusOK, Body: "response 1"},
				{StatusCode: http.StatusCreated, Body: "response 2"},
				{StatusCode: http.StatusOK, DropConnection: true},
			},
		},
	})

	defer srvFake.Close()

	// ARRANGE.
	srvURL := srvFake.URL()

	// SCENARIO #1.
	t.Run("Request a URL for the first (configured) time.", func(t *testing.T) {
		// ACT & ASSERT.
		resp, err := http.Get(srvURL + "/test?v=10")

		assert.Nilf(t, err, "\n\n"+
			"UT Name:  The first response should NOT return an error.\n"+
			"\033[32mExpected: <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equalf(t, resp.StatusCode, http.StatusOK, "\n\n"+
			"UT Name:  The first response should return the 200 status code.\n"+
			"\033[32mExpected: %d\033[0m\n"+
			"\033[31mActual:   %d\033[0m\n\n", http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		assert.Nilf(t, err, "\n\n"+
			"UT Name:  Reading the first response should NOT return an error.\n"+
			"\033[32mExpected: <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equalf(t, string(body), "response 1", "", "\n\n"+
			"UT Name:  The first response should math 'response 1'.\n"+
			"\033[32mExpected: %s\033[0m\n"+
			"\033[31mActual:   %s\033[0m\n\n", "response 1", string(body))
	})

	// SCENARIO #2.
	t.Run("Request a URL for the second (configured) time.", func(t *testing.T) {
		// ACT & ASSERT .
		resp, err := http.Get(srvURL + "/test?v=10")

		assert.Nilf(t, err, "\n\n"+
			"UT Name:  The second response should NOT return an error.\n"+
			"\033[32mExpected: <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equalf(t, resp.StatusCode, http.StatusCreated, "\n\n"+
			"UT Name:  The second response should return the 200 status code.\n"+
			"\033[32mExpected: %d\033[0m\n"+
			"\033[31mActual:   %d\033[0m\n\n", http.StatusCreated, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		assert.Nilf(t, err, "\n\n"+
			"UT Name:  Reading the second response should NOT return an error.\n"+
			"\033[32mExpected: <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equalf(t, string(body), "response 2", "\n\n"+
			"UT Name:  The second response should math 'response 2'.\n"+
			"\033[32mExpected: %s\033[0m\n"+
			"\033[31mActual:   %s\033[0m\n\n", "response 2", string(body))
	})

	// SCENARIO #3.
	t.Run("Request a URL for the third (configured) time.", func(t *testing.T) {
		// ACT & ASSERT .
		resp, err := http.Get(srvURL + "/test?v=10")

		assert.Nilf(t, err, "\n\n"+
			"UT Name:  The third response should NOT return an error.\n"+
			"\033[32mExpected: <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equalf(t, resp.StatusCode, http.StatusOK, "\n\n"+
			"UT Name:  The third response should return the 200 status code.\n"+
			"\033[32mExpected: %d\033[0m\n"+
			"\033[31mActual:   %d\033[0m\n\n", http.StatusOK, resp.StatusCode)

		_, err = io.ReadAll(resp.Body)
		resp.Body.Close()

		assert.NotNilf(t, err, "\n\n"+
			"UT Name:  The third response should return an error when being read.\n"+
			"\033[32mExpected: NOT <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)
	})

	// SCENARIO #4.
	t.Run("Request a URL for the fourth (exhausted) time.", func(t *testing.T) {
		// ACT & ASSERT .
		resp, err := http.Get(srvURL + "/test?v=10")

		assert.Nilf(t, err, "\n\n"+
			"UT Name:  The fourth response should NOT return an error.\n"+
			"\033[32mExpected: <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equalf(t, resp.StatusCode, http.StatusNotImplemented, "\n\n"+
			"UT Name:  The fourth response should return the 501 status code.\n"+
			"\033[32mExpected: %d\033[0m\n"+
			"\033[31mActual:   %d\033[0m\n\n", http.StatusNotImplemented, resp.StatusCode)
	})

	// SCENARIO #5.
	t.Run("Request a NON existing URL.", func(t *testing.T) {
		// ACT & ASSERT .
		resp, err := http.Get(srvURL + "/")

		assert.Nilf(t, err, "\n\n"+
			"UT Name:  The response for a NON existing URL should NOT return an error.\n"+
			"\033[32mExpected: <nil>\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", err)

		assert.Equalf(t, resp.StatusCode, http.StatusNotImplemented, "\n\n"+
			"UT Name:  The response for a NON existing URL should return the 501 status code.\n"+
			"\033[32mExpected: %d\033[0m\n"+
			"\033[31mActual:   %d\033[0m\n\n", http.StatusNotImplemented, resp.StatusCode)

		resp.Body.Close()
	})
}
