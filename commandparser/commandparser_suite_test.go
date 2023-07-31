package commandparser_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCommandParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Command Parser Suite")
}
