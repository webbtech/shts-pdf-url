package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	lerrors "github.com/webbtech/shts-pdf-gen/errors"
	"github.com/webbtech/shts-pdf-url/awsservices"
	"github.com/webbtech/shts-pdf-url/config"
	"github.com/webbtech/shts-pdf-url/model"
)

/*
Steps (methods required) to process Signed URL include:
1. validate input
2. create S3 object string
3. create signed url
*/

type SignedURL struct {
	Cfg         *config.Config
	input       *model.UrlRequest
	request     events.APIGatewayProxyRequest
	response    events.APIGatewayProxyResponse
	s3ObjectKey string
	signedUrl   string
}

const (
	ERR_INVALID_TYPE         = "Invalid request type in input"
	ERR_MISSING_REQUEST_BODY = "Missing request body"
	CODE_SUCCESS             = "SUCCESS"
)

// FIXME: these don't (shouldn't be) uppercase vars
var (
	Stage             string
	TypePrefixMap     = map[string]string{"estimate": "est", "invoice": "inv"}
	ValidRequestTypes = []string{"estimate", "invoice"}
)

// ========================== Public Methods =============================== //
func (su *SignedURL) Response(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	su.request = request
	su.process()
	return su.response, nil
}

// ========================== Private Methods ============================== //

func (su *SignedURL) process() {

	rb := responseBody{}
	var body []byte
	var err error
	var statusCode int = 201
	var stdError *lerrors.StdError
	var url string

	// Validate input
	if err := su.validateInput(); err != nil {
		errors.As(err, &stdError)
	}

	// Create s3ObjectKey and signed url
	if stdError == nil {
		reqType := *su.input.RequestType
		prefix := TypePrefixMap[reqType]
		su.s3ObjectKey = fmt.Sprintf("%s/%s-%d.pdf", reqType, prefix, *su.input.EstimateNumber)

		url, err = awsservices.CreateSignedURL(su.Cfg, su.s3ObjectKey)
		if err != nil {
			stdError = &lerrors.StdError{
				Caller:     "awsservices.CreateSignedURL",
				Code:       lerrors.CodeApplicationError,
				Err:        err,
				Msg:        "Failed to create signed URL",
				StatusCode: 400,
			}
		}
		// fmt.Printf("url: %+v\n", url)
		// fmt.Printf("err: %+v\n", err)
	}

	// Process any errors
	if stdError != nil {
		rb.Code = stdError.Code
		rb.Message = stdError.Msg
		statusCode = stdError.StatusCode
		logError(stdError)
	} else {
		rb.Code = CODE_SUCCESS
		rb.Message = "Success"
		rb.Data = url
	}

	// Create the response object
	body, _ = json.Marshal(&rb)
	su.response = events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    headers,
		StatusCode: statusCode,
	}
}

func (su *SignedURL) validateInput() (err *lerrors.StdError) {

	var inputErrs []string
	validate := validator.New()

	json.Unmarshal([]byte(su.request.Body), &su.input)

	if su.input == nil {
		err = &lerrors.StdError{
			Caller:     "handlers.validateInput",
			Code:       lerrors.CodeBadInput,
			Err:        errors.New(ERR_MISSING_REQUEST_BODY),
			Msg:        ERR_MISSING_REQUEST_BODY,
			StatusCode: 400,
		}
		return err
	}

	// validate struct based on tags
	// see https://github.com/go-playground/validator
	valErr := validate.Struct(su.input)
	if valErr != nil {
		// for more on usage, see: https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
		for _, err := range valErr.(validator.ValidationErrors) {
			inputErrs = append(inputErrs, fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", err.Field(), err.Tag()))
		}
	}

	// validate RequestType
	if su.input.RequestType != nil {
		_, found := findString(ValidRequestTypes, *su.input.RequestType)
		if !found {
			errStr := fmt.Sprintf(ERR_INVALID_TYPE+": \"%s\"", *su.input.RequestType)
			inputErrs = append(inputErrs, errStr)
		}
	}

	if len(inputErrs) > 0 {
		err = &lerrors.StdError{
			Caller:     "handlers.validateInput",
			Code:       lerrors.CodeBadInput,
			Err:        errors.New("Failed request input validation"),
			Msg:        strings.Join(inputErrs, "\n"),
			StatusCode: 400,
		}
		return err
	}

	return nil
}

// NOTE: these could go into it's own package
func logError(err *lerrors.StdError) {
	if Stage == "" {
		Stage = os.Getenv("Stage")
	}

	if Stage != "test" {
		log.Error(err)
	}
}

func findString(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}

	return -1, false
}
