package colorlogger_test

import (
	"bytes"

	. "github.com/cloudfoundry/stembuild/colorlogger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stdout", func() {

	It("write debug output when log level is debug", func() {
		buf := bytes.Buffer{}

		Expect(DEBUG).To(Equal(0))
		logger := ConstructLogger(DEBUG, false, &buf)

		message := "This is a test"

		logger.Logf(DEBUG, message)
		Expect(buf.String()).To(Equal("debug: " + message + "\n"))
	})

	It("write no debug output when log level is NONE", func() {
		buf := bytes.Buffer{}

		logger := ConstructLogger(NONE, false, &buf)

		message := "This is a test"

		logger.Logf(DEBUG, message)
		Expect(buf.String()).To(BeEmpty())
	})

	It("write no none output when log level is NONE", func() {
		buf := bytes.Buffer{}

		logger := ConstructLogger(NONE, false, &buf)

		message := "This is a test"

		logger.Logf(NONE, message)
		Expect(buf.String()).To(BeEmpty())
	})

	It("write colored debug output when log level is DEBUG and color is true", func() {
		buf := bytes.Buffer{}

		logger := ConstructLogger(DEBUG, true, &buf)

		message := "This is a test"

		logger.Logf(DEBUG, message)
		Expect(buf.String()).To(Equal("\033[32mdebug:\033[0m " + message + "\n"))
	})

})
