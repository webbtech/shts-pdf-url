package config

import "time"

type defaults struct {
	AwsRegion string `yaml:"AwsRegion"`
	ExpireHrs int    `yaml:"ExpireHrs"`
	S3Bucket  string `yaml:"S3Bucket"`
	Stage     string `yaml:"Stage"`
}

type config struct {
	AwsRegion     string
	S3Bucket      string
	Stage         StageEnvironment
	UrlExpireTime time.Duration
}
