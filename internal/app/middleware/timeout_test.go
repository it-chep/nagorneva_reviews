package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type countingWriter struct {
	*httptest.ResponseRecorder
	calls int
}

func (w *countingWriter) WriteHeader(code int) { w.calls++; w.ResponseRecorder.WriteHeader(code) }
func TestTimeoutDoesNotWriteStatusTwice(t *testing.T) {
	for _, code := range []int{0, 200, 500, 504} {
		w := &countingWriter{ResponseRecorder: httptest.NewRecorder()}
		handler := Timeout(time.Nanosecond)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
			if code != 0 {
				w.WriteHeader(code)
			}
		}))
		handler.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
		want := code
		if want == 0 {
			want = 504
		}
		if w.calls != 1 || w.Code != want {
			t.Fatalf("code=%d calls=%d status=%d", code, w.calls, w.Code)
		}
	}
}
