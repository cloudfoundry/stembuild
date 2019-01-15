package construct_test

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	rand.Seed(time.Now().UnixNano())
})

func TestCommandParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Construct Suite")
}
