package integration_test

import (
	"math/rand"
	"os"
	"time"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/ovftool"
	"github.com/cloudfoundry-incubator/stembuild/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Integration Suite")
}

var stembuildExecutable string

var _ = SynchronizedBeforeSuite(func() []byte {
	rand.Seed(time.Now().UnixNano())
	Expect(helpers.CopyRecursive(".", "../test/data")).To(Succeed())
	Expect(CheckOVFToolOnPath()).To(Succeed())

	return nil
}, func(_ []byte) {
})

var _ = SynchronizedAfterSuite(func() {
}, func() {
	Expect(os.RemoveAll("./data")).To(Succeed())
})

func CheckOVFToolOnPath() error {
	searchPaths, err := ovftool.SearchPaths()
	if err != nil {
		return err
	}
	if _, err := ovftool.Ovftool(searchPaths); err != nil {
		return err
	}
	return nil
}
