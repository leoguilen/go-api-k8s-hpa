package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServerInfoHandler_GivenHttpCallWithDifferentMethodOfGET_ThenReturnStatusCode405(t *testing.T) {
	// Arrange
	r := httptest.NewRequest(http.MethodPost, "http://example.com", nil)
	w := httptest.NewRecorder()

	// Act
	Handler(w, r)

	// Assert
	if want, got := http.StatusMethodNotAllowed, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestServerInfoHandler_GivenHttpCallWithGETMethod_ThenReturnStatusCode200(t *testing.T) {
	// Arrange
	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	w := httptest.NewRecorder()

	// Act
	Handler(w, r)

	// Assert
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestServerInfoHandler_GivenHttpCallWithSuccessfulGETMethod_ThenResponseBodyShouldContainAJsonWithTheExpectedFormat(t *testing.T) {
	// Arrange
	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	w := httptest.NewRecorder()
	expected := NewContextInfo(r)

	// Act
	Handler(w, r)

	// Assert
	var result *ContextDetails
	if err := json.NewDecoder(w.Result().Body).Decode(&result); err != nil {
		t.Fatalf("decoding response body failed")
	}

	result.Request.Identifier = expected.Request.Identifier
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("expected a %v, instead got: %v", expected, result)
	}
}
