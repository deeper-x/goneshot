package goneshot

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

	isAuthorized(testPage)
	res := w.Result()

	defer res.Body.Close()

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
}
