package stemcell_generator_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStemcellGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StemcellGenerator Suite")
}
