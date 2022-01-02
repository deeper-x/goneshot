package goneshot

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Super Secret test Information")
}

func IndexDemo(w http.ResponseWriter, r *http.Request) {
	client := r.URL.Query().Get("client")
	validToken, err := GetJWT(client)
	if err != nil {
		fmt.Println("Failed to generate token")
	}

	fmt.Fprint(w, string(validToken))
}

func TestGetJWT(t *testing.T) {
	_, err := GetJWT("client")
	if err != nil {
		t.Error(err)
	}
}

func TestIndex(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	IndexDemo(w, req)
	res := w.Result()

	defer res.Body.Close()

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
}

func TestIsAuthorized(t *testing.T) {
	w := httptest.NewRecorder()

	IsAuthorized(testPage)
	res := w.Result()

	defer res.Body.Close()

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
}

func TestValidateToken(t *testing.T) {
	cookie := &http.Cookie{
		Name:  "Token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}
	w := httptest.NewRecorder()

	http.SetCookie(w, cookie)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Add(cookie.Name, cookie.Value)

	tknstr := w.Header().Get("Set-Cookie")

	tknw := strings.Split(tknstr, "=")[1]
	tknreq := r.Header["Token"][0]

	if tknw != tknreq {
		t.Error("token not valid")
	}

	_, err := validateToken(w, r)

	if err.Error() != "signature is invalid" {
		t.Fatal(err)
	}
}
