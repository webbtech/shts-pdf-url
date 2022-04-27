package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
	"github.com/webbtech/shts-pdf-url/config"
)

// CreateSignedURL function
func CreateSignedURL(cfg *config.Config, fileObject string) (url string, err error) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AwsRegion),
	})
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(cfg.S3Bucket),
		Key:    aws.String(fileObject),
	})

	urlStr, err := req.Presign(cfg.UrlExpireTime)
	if err != nil {
		log.Errorf("Failed to sign request: %s", err.Error())
		return "", err
	}

	return urlStr, nil
}
