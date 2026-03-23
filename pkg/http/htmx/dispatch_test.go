package htmx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDispatch(t *testing.T) {
	tests := []struct {
		name           string
		headers        map[string]string
		wantHtmxCalled bool
	}{
		{
			name:           "HTMX request",
			headers:        map[string]string{"HX-Request": "true"},
			wantHtmxCalled: true,
		},
		{
			name:           "Non-HTMX request",
			headers:        map[string]string{},
			wantHtmxCalled: false,
		},
		{
			name:           "HTMX request with wrong value",
			headers:        map[string]string{"HX-Request": "false"},
			wantHtmxCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request with the specified headers
			req := httptest.NewRequest("GET", "/", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			// Create a recorder to capture the response
			rec := httptest.NewRecorder()

			// Track which handler was called
			htmxCalled := false
			otherwiseCalled := false

			// Create the handlers
			htmxHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				htmxCalled = true
			})
			otherwiseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				otherwiseCalled = true
			})

			// Create the dispatch handler and serve the request
			handler := Dispatch(htmxHandler, otherwiseHandler)
			handler.ServeHTTP(rec, req)

			// Verify the correct handler was called
			if tt.wantHtmxCalled && !htmxCalled {
				t.Error("HTMX handler was not called when it should have been")
			}
			if !tt.wantHtmxCalled && htmxCalled {
				t.Error("HTMX handler was called when it should not have been")
			}
			if !tt.wantHtmxCalled && !otherwiseCalled {
				t.Error("Otherwise handler was not called when it should have been")
			}
			if tt.wantHtmxCalled && otherwiseCalled {
				t.Error("Otherwise handler was called when it should not have been")
			}
		})
	}
}
