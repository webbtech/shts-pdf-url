package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestPingHandler(t *testing.T) {

	os.Setenv("Stage", "test")

	var msg string
	t.Run("Successful ping", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer ts.Close()

		r, err := handler(events.APIGatewayProxyRequest{Path: "/"})

		expectedMsg := "Healthy"
		msg = extractMessage(r.Body, "message")
		if msg != expectedMsg {
			t.Fatalf("Expected error message: %s received: %s", expectedMsg, msg)
		}
		if err != nil {
			t.Fatal("Everything should be ok")
		}
	})
}

func TestEnvVars(t *testing.T) {
	os.Setenv("PARAM1", "VALUE12")
	t.Run("Successful ping", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer ts.Close()

		_, _ = handler(events.APIGatewayProxyRequest{Path: "/"})
		p, exists := os.LookupEnv("PARAM1")
		if !exists {
			t.Fatalf("Expected value for PARAM1 to be: %s", p)
		}
	})
}

func TestURLHandler(t *testing.T) {
	os.Setenv("Stage", "test")

	var msg string
	t.Run("succesful url", func(t *testing.T) {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(201)
		}))
		defer ts.Close()

		requestBody := `{"number": 1011, "requestType": "estimate"}` // These values must be from a valid s3 object
		r, err := handler(events.APIGatewayProxyRequest{Path: "/url", Body: requestBody})
		if err != nil {
			t.Fatal("Everything should be ok")
		}

		msg = extractMessage(r.Body, "data")

		expectedStrStart := "https://shts-pdf.s3.ca-central-1.amazonaws.com/estimate/est-1011.pdf"
		if !strings.HasPrefix(msg, expectedStrStart) {
			t.Fatalf("Expected message to start with: %s, message was: %s", expectedStrStart, msg)
		}
	})
}

func extractMessage(b, key string) (msg string) {
	var dat map[string]string
	_ = json.Unmarshal([]byte(b), &dat)
	return dat[key]
}
