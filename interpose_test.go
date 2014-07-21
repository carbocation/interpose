package interpose

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	middle := New()

	middle.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		result += "0"
	}))

	middle.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			result += "1"
			next.ServeHTTP(rw, req)
			result += "1"
		})
	})

	middle.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			result += "2"
			next.ServeHTTP(rw, req)
			result += "2"
		})
	})

	middle.ServeHTTP(response, (*http.Request)(nil))
	expect(t, result, "21012")
}

func TestEmptyMiddleware(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	middle := New()

	middle.ServeHTTP(response, (*http.Request)(nil))
	expect(t, result, "")
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
