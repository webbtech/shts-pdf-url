package services

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/webbtech/shts-pdf-url/config"
)

// Found much of this code here: https://aws.github.io/aws-sdk-go-v2/docs/code-examples/s3/generatepresignedurl/
// TODO: tests should include a mocked version of s3.NewPresignClient

// S3PresignGetObjectAPI defines the interface for the PresignGetObject function.
// We use this interface to test the function using a mocked service.
type S3PresignGetObjectAPI interface {
	PresignGetObject(
		ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// GetPresignedURL retrieves a presigned URL for an Amazon S3 bucket object.
// Inputs:
//    c is the context of the method call, which includes the AWS Region.
//    api is the interface that defines the method call.
//    input defines the input arguments to the service call.
// 		optFns is an array of s3.PresignOptions functions
// Output:
//     If successful, the presigned URL for the object and nil.
//     Otherwise, nil and an error from the call to PresignGetObject.
func GetPresignedURL(c context.Context, api S3PresignGetObjectAPI, input *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	return api.PresignGetObject(c, input, optFns...)
}

// CreateSignedURL function
func CreateSignedURL(cfg *config.Config, fileObject string) (url string, err error) {

	acfg, err := awscfg.LoadDefaultConfig(context.TODO(),
		awscfg.WithRegion(cfg.AwsRegion),
	)
	if err != nil {
		return url, err
	}

	client := s3.NewFromConfig(acfg)

	psClient := s3.NewPresignClient(client)
	input := &s3.GetObjectInput{
		Bucket: aws.String(cfg.S3Bucket),
		Key:    aws.String(fileObject),
	}
	expiresOpt := s3.WithPresignExpires(cfg.UrlExpireTime)

	resp, err := GetPresignedURL(context.TODO(), psClient, input, expiresOpt)
	if err != nil {
		return url, err
	}

	return resp.URL, nil
}
