package requests

import (
	"net/http"
	"testing"
)

func TestDefaultModifer(t *testing.T) {
	t.Run("application/json headers added", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodConnect, "https://google.com", nil)
		if err != nil {
			t.Errorf("error generating http request: %v", err)
			return
		}

		DefaultModifer(req)

		ctGot := req.Header.Get("Content-Type")
		aGot := req.Header.Get("Accept")
		if ctGot != applicationJson {
			t.Errorf("request header Content-Type not set correctly, ct = %s", ctGot)
		}
		if aGot != applicationJson {
			t.Errorf("request header Accept not set correctly, app = %s", aGot)
		}
	})
}
