package csrf

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gdisw/resume/pkg/testutil"
	"github.com/justinas/nosurf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerator_GET_NonHTMX(t *testing.T) {
	handler := new(testutil.MockHandler)
	generator := Generator(handler)

	handler.
		On("ServeHTTP", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			r := args.Get(1).(*http.Request)
			token := nosurf.Token(r)
			assert.NotEmpty(t, token, "CSRF token should be available for GET non-HTMX requests")
		}).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	nosurfHandler := nosurf.New(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		generator.ServeHTTP(w, r)
	}))
	nosurfHandler.SetBaseCookie(http.Cookie{Path: "/"})

	nosurfHandler.ServeHTTP(rec, req)
	handler.AssertExpectations(t)
}

func TestGenerator_GET_HTMX(t *testing.T) {
	handler := new(testutil.MockHandler)
	generator := Generator(handler)

	handler.
		On("ServeHTTP", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			r := args.Get(1).(*http.Request)
			assert.NotNil(t, r.Context(), "Context should be present")
		}).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()

	generator.ServeHTTP(rec, req)
	handler.AssertExpectations(t)
}

func TestGenerator_POST_NonHTMX(t *testing.T) {
	handler := new(testutil.MockHandler)
	generator := Generator(handler)

	handler.
		On("ServeHTTP", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			r := args.Get(1).(*http.Request)
			assert.NotNil(t, r.Context(), "Context should be present")
		}).
		Once()

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	generator.ServeHTTP(rec, req)
	handler.AssertExpectations(t)
}

func TestGenerator_POST_HTMX(t *testing.T) {
	handler := new(testutil.MockHandler)
	generator := Generator(handler)

	handler.
		On("ServeHTTP", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			r := args.Get(1).(*http.Request)
			assert.NotNil(t, r.Context(), "Context should be present")
		}).
		Once()

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()

	generator.ServeHTTP(rec, req)
	handler.AssertExpectations(t)
}

func TestGenerator_PUT_NonHTMX(t *testing.T) {
	handler := new(testutil.MockHandler)
	generator := Generator(handler)

	handler.
		On("ServeHTTP", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			r := args.Get(1).(*http.Request)
			assert.NotNil(t, r.Context(), "Context should be present")
		}).
		Once()

	req := httptest.NewRequest(http.MethodPut, "/", nil)
	rec := httptest.NewRecorder()

	generator.ServeHTTP(rec, req)
	handler.AssertExpectations(t)
}

func TestGenerator_HTMXHeaderCheck(t *testing.T) {
	testCases := []struct {
		name           string
		headerValue    string
		shouldAddToken bool
		method         string
	}{
		{"GET with lowercase true", "true", false, http.MethodGet},
		{"GET with uppercase TRUE", "TRUE", false, http.MethodGet},
		{"GET with mixed case True", "True", false, http.MethodGet},
		{"GET with false value", "false", true, http.MethodGet},
		{"GET with empty header", "", true, http.MethodGet},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := new(testutil.MockHandler)
			generator := Generator(handler)

			handler.
				On("ServeHTTP", mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					r := args.Get(1).(*http.Request)
					if tc.shouldAddToken && tc.method == http.MethodGet {
						assert.NotNil(t, r.Context())
					}
				}).
				Once()

			req := httptest.NewRequest(tc.method, "/", nil)
			if tc.headerValue != "" {
				req.Header.Set("HX-Request", tc.headerValue)
			}
			rec := httptest.NewRecorder()

			if tc.shouldAddToken && tc.method == http.MethodGet {
				nosurfHandler := nosurf.New(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					generator.ServeHTTP(w, r)
				}))
				nosurfHandler.SetBaseCookie(http.Cookie{Path: "/"})
				nosurfHandler.ServeHTTP(rec, req)
			} else {
				generator.ServeHTTP(rec, req)
			}

			handler.AssertExpectations(t)
		})
	}
}

func TestGenerator_AlwaysCallsNextHandler(t *testing.T) {
	handler := new(testutil.MockHandler)
	generator := Generator(handler)

	handler.
		On("ServeHTTP", mock.Anything, mock.Anything).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()

	generator.ServeHTTP(rec, req)
	handler.AssertExpectations(t)
}

func TestGenerator_DifferentHTTPMethods(t *testing.T) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodHead,
		http.MethodOptions,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			handler := new(testutil.MockHandler)
			generator := Generator(handler)

			handler.
				On("ServeHTTP", mock.Anything, mock.Anything).
				Once()

			req := httptest.NewRequest(method, "/", nil)
			rec := httptest.NewRecorder()

			generator.ServeHTTP(rec, req)
			handler.AssertExpectations(t)
		})
	}
}
