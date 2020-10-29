package package_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Package Suite")
}

const (
	VcenterCACert = "VCENTER_CA_CERT"
)

var (
	pathToCACert string
)

var _ = SynchronizedBeforeSuite(func() []byte {
	rawCA := envMustExist(VcenterCACert)
	t, err := ioutil.TempFile("", "ca-cert")
	Expect(err).ToNot(HaveOccurred())
	pathToCACert = t.Name()
	Expect(t.Close()).To(Succeed())
	err = ioutil.WriteFile(pathToCACert, []byte(rawCA), 0666)
	Expect(err).ToNot(HaveOccurred())
}, func(_ []byte) {
})

var _ = SynchronizedAfterSuite(func() {
}, func() {
	if pathToCACert != "" {
		os.RemoveAll(pathToCACert)
	}
})

func envMustExist(variableName string) string {
	result := os.Getenv(variableName)
	if result == "" {
		Fail(fmt.Sprintf("%s must be set", variableName))
	}

	return result
}
