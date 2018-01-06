package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "", nil)

	if err != nil {
		t.Fatal(err)
	}

	// We use Go's httptest library to create an http recorder. This recorder
	// will act as the target of our http request
	// (you can think of it as a mini-browser, which will accept the result of
	// the http request that we make)
	recorder := httptest.NewRecorder()

	// Create an HTTP handler from our handler function. "handler" is the handler
	// function defined in our main.go file that we want to test
	handlerFunction := http.HandlerFunc(handler)

	// Serve the HTTP request to our recorder. This is the line that actually
	// executes our the handler that we want to test
	handlerFunction.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `Hello World!`
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestRouter(t *testing.T) {
	router := newRouter()
	mockServer := httptest.NewServer(router)

	response, err := http.Get(mockServer.URL + "/hello")

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Status should be ok, got %d", response.StatusCode)
	}

	defer response.Body.Close()

	// read the body into a bunch of bytes (b)
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	// convert the bytes to a string
	responseString := string(b)
	expected := "Hello World!"

	// We want our response to match the one defined in our handler.
	// If it does happen to be "Hello world!", then it confirms, that the
	// route is correct
	if responseString != expected {
		t.Errorf("Response should be %s, got %s", expected, responseString)
	}
}

func TestRouterForNonExistentRoute(t *testing.T) {
	router := newRouter()
	mockServer := httptest.NewServer(router)

	// Most of the code is similar. The only difference is that now we make a
	//request to a route we know we didn't define, like the `POST /hello` route.
	response, err := http.Post(mockServer.URL+"/hello", "", nil)

	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 405 (method not allowed)
	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status should be 405, got %d", response.StatusCode)
	}

	// The code to test the body is also mostly the same, except this time, we
	// expect an empty body
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)

	if err != nil {
		t.Fatal(err)
	}

	// convert the bytes to a string
	responseString := string(b)
	expected := ""

	if responseString != expected {
		t.Errorf("Response should be %s, got %s", expected, responseString)
	}
}

func TestStaticFileServer(t *testing.T) {
	router := newRouter()
	mockServer := httptest.NewServer(router)

	response, err := http.Get(mockServer.URL + "/assets/")
	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 200 (ok)
	if response.StatusCode != http.StatusOK {
		t.Errorf("Status should be 200, got %d", response.StatusCode)
	}

	// It isn't wise to test the entire content of the HTML file.
	// Instead, we test that the content-type header is "text/html; charset=utf-8"
	// so that we know that an html file has been served
	contentType := response.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
	}
}
