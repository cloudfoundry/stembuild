package vcenter_manager_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVcenterManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VcenterManager Suite")
}
