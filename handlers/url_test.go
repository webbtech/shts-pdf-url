package handlers

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/webbtech/shts-pdf-url/config"
)

func TestValidateInput(t *testing.T) {
	t.Run("Missing request body", func(t *testing.T) {
		su := &SignedURL{}

		expectedErr := ERR_MISSING_REQUEST_BODY
		err := su.validateInput()
		if err.Msg != expectedErr {
			t.Fatalf("Error should be: %s, have: %s", expectedErr, err.Msg)
		}
	})

	t.Run("Missing struct fields", func(t *testing.T) {
		su := &SignedURL{}
		su.request.Body = `{}`

		err := su.validateInput()
		nLines := strings.Split(err.Msg, "\n")
		expectedNumErrs := 2
		haveLines := len(nLines)
		if expectedNumErrs != haveLines {
			t.Fatalf("Number of Msg errors should be: %d, have: %d", expectedNumErrs, haveLines)
		}

		expectedError1 := "Field validation for 'EstimateNumber' failed on the 'required' tag"
		haveError1 := nLines[0]
		if expectedError1 != haveError1 {
			t.Fatalf("First error message should be: %s, have: %s", expectedError1, haveError1)
		}
	})

	t.Run("Invalid type input", func(t *testing.T) {
		su := &SignedURL{}
		su.request.Body = `{"number": 900, "requestType": "estimat"}`

		err := su.validateInput()
		if !strings.HasPrefix(err.Msg, ERR_INVALID_TYPE) {
			t.Fatalf("Expected error to start with: %s", ERR_INVALID_TYPE)
		}
	})
}

func TestProcess(t *testing.T) {
	t.Run("Successfully creates S3Object string", func(t *testing.T) {
		cfg := &config.Config{DefaultsFilePath: "../config/defaults.yml"}
		err := cfg.Init()
		if err != nil {
			t.Fatalf("Expected null error received: %s", err)
		}
		su := &SignedURL{cfg: cfg}
		su.request.Body = `{"number": 1011, "requestType": "estimate"}`

		su.process()
		expectedObjStr := "estimate/est-1011.pdf"
		if expectedObjStr != su.s3ObjectKey {
			t.Fatalf("s3Object should be: %s, have: %s", expectedObjStr, su.s3ObjectKey)
		}
	})

	// TODO: require test here to successfully open a signed url
	t.Run("create signed url", func(t *testing.T) {
		cfg := &config.Config{DefaultsFilePath: "../config/defaults.yml"}
		err := cfg.Init()
		if err != nil {
			t.Fatalf("Expected null error received: %s", err)
		}
		su := &SignedURL{cfg: cfg}
		su.request.Body = `{"number": 1011, "requestType": "estimate"}`
		su.process()

		expectedStatusCode := 201
		if expectedStatusCode != su.response.StatusCode {
			t.Fatalf("Status should be: %d, have: %d", expectedStatusCode, su.response.StatusCode)
		}

		expectedSuccessCode := CODE_SUCCESS
		var responseBody = &responseBody{}
		json.Unmarshal([]byte(su.response.Body), &responseBody)

		if expectedSuccessCode != responseBody.Code {
			t.Fatalf("SuccessCode should be: %s, have: %s", expectedSuccessCode, responseBody.Code)
		}

		expectedMessageStart := "https://shts-pdf.s3.ca-central-1.amazonaws.com/"
		if !strings.HasPrefix(responseBody.Message, expectedMessageStart) {
			t.Fatalf("Expected message to start with: %s", expectedMessageStart)
		}
	})
}
