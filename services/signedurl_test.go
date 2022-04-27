package services

import (
	"fmt"
	"testing"

	"github.com/webbtech/shts-pdf-url/config"
)

func TestCreateSignedURL(t *testing.T) {
	t.Run("successfully return signed url", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Init()

		fileObject := "estimate/est-1005.pdf"
		url, err := CreateSignedURL(cfg, fileObject)
		if err != nil {
			t.Fatalf("Expected null error received: %s", err)
		}

		fmt.Printf("url: %+v\n", url)

	})
}
