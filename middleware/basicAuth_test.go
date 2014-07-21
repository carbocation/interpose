package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carbocation/interpose"
)

var comparetests = []struct {
	a   string
	b   string
	val bool
}{
	{"foo", "foo", true},
	{"bar", "bar", true},
	{"password", "password", true},
	{"Foo", "foo", false},
	{"foo", "foobar", false},
	{"password", "pass", false},
}

func Test_SecureCompare(t *testing.T) {
	for _, tt := range comparetests {
		if secureCompare(tt.a, tt.b) != tt.val {
			t.Errorf("Expected SecureCompare(%v, %v) to return %v but did not", tt.a, tt.b, tt.val)
		}
	}
}

func Test_BasicAuth(t *testing.T) {
	recorder := httptest.NewRecorder()

	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte("foo:bar"))

	i := interpose.New()

	i.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("hello"))
			next.ServeHTTP(w, req)
		})
	})
	i.Use(BasicAuth("foo", "bar"))

	r, _ := http.NewRequest("GET", "foo", nil)
	i.ServeHTTP(recorder, r)

	if recorder.Code != 401 {
		t.Errorf("recorder.Code wrong. Got %d wanted 401", recorder.Code)
	}

	if recorder.Body.String() == "hello" {
		t.Error("Auth block failed")
	}

	recorder = httptest.NewRecorder()
	r.Header.Set("Authorization", auth)
	i.ServeHTTP(recorder, r)

	if recorder.Code == 401 {
		t.Error("Response is 401")
	}

	if recorder.Body.String() != "hello" {
		t.Error("Auth failed, got: ", recorder.Body.String())
	}
}
