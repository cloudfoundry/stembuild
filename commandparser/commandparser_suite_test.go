package commandparser_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCommandParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Command Parser Suite")
}
