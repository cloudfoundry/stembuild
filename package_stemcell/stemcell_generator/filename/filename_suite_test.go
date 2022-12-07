package filename_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFilename(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filename Suite")
}
