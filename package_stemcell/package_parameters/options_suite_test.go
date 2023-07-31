package package_parameters_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestStembuildOptions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VmdkPackageParameters Suite")
}
