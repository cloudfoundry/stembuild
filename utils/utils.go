package utils

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func DownloadFileFromURL(destPath string, patchURL string, debugf func(string, ...interface{})) (string, error) {
	patchPath := filepath.Join(destPath, fmt.Sprintf("patchfile-%d", rand.Intn(2000)))
	myFile, err := os.Create(patchPath)
	if err != nil {
		return "", fmt.Errorf("Could not create create downloaded file in directory %s", destPath)
	}

	defer myFile.Close()

	debugf("Downloading patch file from %s", patchURL)
	response, err := http.Get(patchURL)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Could not create stemcell from %s\\nUnexpected response code: %d", patchURL, response.StatusCode)
	}

	_, err = io.Copy(myFile, response.Body)
	if err != nil {
		return "", err
	}
	debugf("Finished downloading patchfile")
	return patchPath, nil
}

func ValidateVersion(s string) error {
	if s == "" {
		return errors.New("missing required argument 'version'")
	}
	patterns := []string{
		`^\d{1,}\.\d{1,}$`,
		`^\d{1,}\.\d{1,}-build\.\d{1,}$`,
		`^\d{1,}\.\d{1,}\.\d{1,}$`,
		`^\d{1,}\.\d{1,}\.\d{1,}-build\.\d{1,}$`,
	}
	for _, pattern := range patterns {
		if regexp.MustCompile(pattern).MatchString(s) {
			return nil
		}
	}
	return fmt.Errorf("invalid version (%s) expected format [NUMBER].[NUMBER] or "+
		"[NUMBER].[NUMBER].[NUMBER]", s)
}
