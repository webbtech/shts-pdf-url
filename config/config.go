package config

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"reflect"
	"time"

	"gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	config
	IsDefaultsLocal bool
}

// StageEnvironment string
type StageEnvironment string

// DB type constants
const (
	DevEnv   StageEnvironment = "dev"
	StageEnv StageEnvironment = "stage"
	TestEnv  StageEnvironment = "test"
	ProdEnv  StageEnvironment = "prod"
)

const (
	defaultFileName    = "url-defaults.yml"
	defaultsRemotePath = "https://shts-pdf.s3.ca-central-1.amazonaws.com/public/url-defaults.yml"
)

var (
	defs             = &defaults{}
	defaultsFilePath string
)

// ========================== Public Methods =============================== //

// Init method
func (c *Config) Init() (err error) {

	if err = c.setDefaults(); err != nil {
		return err
	}

	if err = c.setEnvVars(); err != nil {
		return err
	}

	c.setFinal()

	return nil
}

// GetStageEnv method
func (c *Config) GetStageEnv() StageEnvironment {
	return c.Stage
}

// SetStageEnv method
func (c *Config) SetStageEnv(env string) (err error) {
	defs.Stage = env
	return c.validateStage()
}

// ========================== Private Methods =============================== //

func (c *Config) setDefaults() (err error) {

	var file []byte
	if c.IsDefaultsLocal == true { // DefaultsRemote is explicitly set to true

		dir, _ := os.Getwd()
		defaultsFilePath = path.Join(dir, defaultFileName)
		if _, err = os.Stat(defaultsFilePath); os.IsNotExist(err) {
			return err
		}

		file, err = ioutil.ReadFile(defaultsFilePath)
		if err != nil {
			return err
		}

	} else { // using remote file path
		res, err := http.Get(defaultsRemotePath)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		file, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
	}

	err = yaml.Unmarshal([]byte(file), &defs)
	if err != nil {
		return err
	}

	err = c.validateStage()
	if err != nil {
		return err
	}

	return err
}

// validateStage method to validate Stage value
func (c *Config) validateStage() (err error) {

	validEnv := true

	switch defs.Stage {
	case "dev":
	case "development":
		c.Stage = DevEnv
	case "stage":
		c.Stage = StageEnv
	case "test":
		c.Stage = TestEnv
	case "prod":
		c.Stage = ProdEnv
	case "production":
		c.Stage = ProdEnv
	default:
		validEnv = false
	}

	if !validEnv {
		return errors.New("Invalid StageEnvironment requested")
	}

	return err
}

// setEnvVars sets any environment variables that match the default struct fields
func (c *Config) setEnvVars() (err error) {

	vals := reflect.Indirect(reflect.ValueOf(defs))
	for i := 0; i < vals.NumField(); i++ {
		nm := vals.Type().Field(i).Name
		if e := os.Getenv(nm); e != "" {
			vals.Field(i).SetString(e)
		}
		// If field is Stage, validate and return error if required
		if nm == "Stage" {
			err = c.validateStage()
			if err != nil {
				return err
			}
		}
	}

	return err
}

// Copies required fields from the defaults to the Config struct
func (c *Config) setFinal() {
	c.AwsRegion = defs.AwsRegion
	c.S3Bucket = defs.S3Bucket
	c.UrlExpireTime = time.Duration(time.Duration(defs.ExpireHrs) * time.Hour)
}
