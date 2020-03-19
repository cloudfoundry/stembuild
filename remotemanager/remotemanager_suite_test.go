package remotemanager_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRemotemanager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Remotemanager Suite")
}
