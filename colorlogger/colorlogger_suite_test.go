package colorlogger_test

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

func TestColorLogger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Color Logger Suite")
}
