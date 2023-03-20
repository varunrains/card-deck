package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var theTests = []struct {
		name               string
		url                string
		expectedStatusCode int
		httpMethod         string
	}{
		{"CreateDeck", "/createDeck", http.StatusOK, "POST"},
		{"CreateDeckWithShuffled", "/createDeck?shuffled=true", http.StatusOK, "POST"},
		{"CreateDeckWithSpecifiedCards", "/createDeck?cards=AH,KH", http.StatusOK, "POST"},
		{"OpenDeck_Without_Parameters", "/openDeck", http.StatusNotFound, "GET"},
		{"OpenDeck_With_Parameters", "/openDeck/1234", http.StatusOK, "GET"},
		{"DrawDeck_Without_Parameters", "/drawDeck", http.StatusNotFound, "GET"},
		{"DrawDeck_Without_count", "/drawDeck/1235", http.StatusBadRequest, "GET"},
		{"DrawDeck_Without_Integercount", "/drawDeck/1235?count=a", http.StatusBadRequest, "GET"},
		{"DrawDeck_With_Parameters", "/drawDeck/1234?count=1", http.StatusOK, "GET"},
		{"404", "/something-unknown", http.StatusNotFound, "GET"},
	}

	routes := app.routes()

	//create a test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	//range through test data
	for _, e := range theTests {
		var resp *http.Response
		var err error
		//resp, err := ts.Client().Get(ts.URL + e.url)
		if e.httpMethod == "GET" {
			resp, err = ts.Client().Get(ts.URL + e.url)
		} else if e.httpMethod == "POST" {
			resp, err = ts.Client().Post(ts.URL+e.url, "application/json", nil)
		}

		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s: expected status %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}
