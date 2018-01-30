package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/pivotal-cf-experimental/stembuild/utils"
)

var _ = Describe("Utils", func() {
	DescribeTable("ValidateVersion",
		func(version string, isErrorExpected bool) {
			Expect(ValidateVersion(version) == nil).To(Equal(isErrorExpected))
		},
		Entry("failed to validate version", "1.2", true),
		Entry("expected error for version", "-1.2", false),
		Entry("expected error for version", "1.-2", false),
		Entry("failed to validate version", "001.002", true),
		Entry("expected error for version", "0a1.002", false),
		Entry("expected error for version", "1.a", false),
		Entry("expected error for version", "a1.2", false),
		Entry("expected error for version", "a.2", false),
		Entry("expected error for version", "1.2 a", false),
		Entry("failed to validate version", "1200.0.3-build.2", true),
		Entry("expected error for version", "1200.0.3-build.a", false),
		Entry("failed to validate version", "1.2-build.1", true),
	)
})
