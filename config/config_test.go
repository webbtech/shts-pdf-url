package config

import (
	"os"
	"path"
	"testing"
	"time"
)

var cfg *Config

func TestInitConfig(t *testing.T) {
	t.Run("Successful Init with local file", func(t *testing.T) {
		cfg = &Config{}
		err := cfg.Init()
		if err != nil {
			t.Fatalf("Expected null error received: %s", err)
		}
	})

	t.Run("Successful Init with remote file", func(t *testing.T) {
		// cfg = &Config{}
		cfg = &Config{DefaultsFilePath: "https://shts-pdf.s3.ca-central-1.amazonaws.com/public/defaults.yml"}
		err := cfg.Init()
		if err != nil {
			t.Fatalf("Expected null error received: %s", err)
		}
	})
}

func TestGetStageEnv(t *testing.T) {
	cfg = &Config{}
	cfg.setDefaults()

	t.Run("successfully get stage value", func(t *testing.T) {
		stg := cfg.GetStageEnv()
		if stg != ProdEnv {
			t.Fatalf("Stage should be %s, have: %s", ProdEnv, stg)
		}
	})

	t.Run("successfully set then get stage value", func(t *testing.T) {
		cfg.SetStageEnv("test")
		stg := cfg.GetStageEnv()
		if stg != TestEnv {
			t.Fatalf("Stage should be %s, have: %s", TestEnv, stg)
		}
	})
}

func TestSetDefaults(t *testing.T) {
	t.Run("test setting DefaultsFilePath", func(t *testing.T) {

		cfg = &Config{}
		cfg.setDefaults()
		dir, _ := os.Getwd()
		expectedFilePath := path.Join(dir, defaultFileName)
		if expectedFilePath != cfg.DefaultsFilePath {
			t.Fatalf("DefaultsFilePath should be %s, have: %s", expectedFilePath, cfg.DefaultsFilePath)
		}

		fp := path.Join("/tmp", defaultFileName)
		cfg = &Config{DefaultsFilePath: fp}
		err := cfg.setDefaults()
		expectedFilePath = fp

		if expectedFilePath != cfg.DefaultsFilePath {
			t.Fatalf("DefaultsFilePath should be %s, have: %s", expectedFilePath, cfg.DefaultsFilePath)
		}
		if err == nil {
			t.Fatalf("setDefaults should return error")
		}
	})
}

// TestValidateStage tests the validateStage method
// validateStage is called at various times including in setEnvVars
func TestValidateStage(t *testing.T) {
	cfg = &Config{}
	cfg.setDefaults()

	t.Run("stage set from defaults file", func(t *testing.T) {
		if cfg.Stage != ProdEnv {
			t.Fatalf("Stage value should be: %s, have: %s", ProdEnv, cfg.Stage)
		}
	})

	t.Run("stage set from environment", func(t *testing.T) {
		os.Setenv("Stage", "test")
		cfg.setEnvVars() // calls validateStage
		if cfg.Stage != TestEnv {
			t.Fatalf("Stage value should be: %s, have: %s", TestEnv, cfg.Stage)
		}
	})

	t.Run("stage set from invalid environment variable", func(t *testing.T) {
		os.Setenv("Stage", "testit")
		err := cfg.setEnvVars()
		if err == nil {
			t.Fatalf("Expected validateStage to return error")
		}
	})

	t.Run("stage set with SetStageEnv method", func(t *testing.T) {
		err := cfg.SetStageEnv("stage")
		if err != nil {
			t.Fatalf("Expected null error received: %s", err)
		}
	})

	t.Run("invalid stage set with SetStageEnv method", func(t *testing.T) {
		err := cfg.SetStageEnv("stageit")
		if err == nil {
			t.Fatalf("Expected validateStage error")
		}
	})
}

func TestUrlExpireTime(t *testing.T) {
	t.Run("sets expireTime", func(t *testing.T) {
		cfg = &Config{}
		cfg.Init()

		expectedHrs := time.Duration(time.Duration(defs.ExpireHrs) * time.Hour)
		if expectedHrs != cfg.UrlExpireTime {
			t.Fatalf("UrlExpireTime should be: %d, have: %d", expectedHrs, cfg.UrlExpireTime)
		}
	})
}
