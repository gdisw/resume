package view

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gdisw/resume/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouteConfigPush(t *testing.T) {
	{
		parent := RouteConfig{PathPrefix: "/parent"}
		child := RouteConfig{PathPrefix: "/child"}
		assert.Equal(t, "/parent/child", parent.push(child).PathPrefix)
	}

	{
		parent := RouteConfig{PathPrefix: "/parent"}
		child := RouteConfig{PathPrefix: "child"}
		assert.Equal(t, "/parent/child", parent.push(child).PathPrefix)
	}

	{
		parent := RouteConfig{}
		child := RouteConfig{PathPrefix: "/path"}
		assert.Equal(t, "/path", parent.push(child).PathPrefix)
	}
}

func TestRouteConfigExtraction(t *testing.T) {
	assert.Equal(t, "/hello", route(ViewData{}, "/hello"))
	assert.Equal(t, "/hello", route(ViewData{"RouteConfig": 12}, "/hello"))
	assert.Equal(t, "/super/hello", route(ViewData{
		"RouteConfig": RouteConfig{PathPrefix: "/super"},
	}, "/hello"))
}

func TestRouteConfigMiddlewareEmptyPrefix(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	PutRouteConfig(RouteConfig{})
}

func TestRouteConfigMiddlewareInvalidPrefix(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	PutRouteConfig(RouteConfig{PathPrefix: "hello"})
}

func TestRouteConfigMiddleware(t *testing.T) {
	handler := new(testutil.MockHandler)
	handler.
		On("ServeHTTP", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			r := args.Get(1).(*http.Request)
			viewData := viewDataFromContext(r.Context())
			assert.Equal(t, "/section/subsection/paragraph/text", route(viewData, "/text"))
		}).
		Once()

	chain := PutRouteConfig(
		RouteConfig{PathPrefix: "/section"},
	)(
		PutRouteConfig(
			RouteConfig{PathPrefix: "/subsection"},
		)(
			PutRouteConfig(
				RouteConfig{PathPrefix: "/paragraph"},
			)(handler),
		),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	chain.ServeHTTP(resp, req)
	handler.AssertExpectations(t)
}
