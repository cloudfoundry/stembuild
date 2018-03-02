package utils_test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/stembuild/helpers"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	rand.Seed(time.Now().UnixNano())
	Expect(helpers.CopyRecursive(".", "../testdata")).To(Succeed())
	return nil
}, func(_ []byte) {
})

var _ = SynchronizedAfterSuite(func() {
}, func() {
	Expect(os.RemoveAll("./testdata")).To(Succeed())
})
