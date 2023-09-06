package package_parameters_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStembuildOptions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VmdkPackageParameters Suite")
}
