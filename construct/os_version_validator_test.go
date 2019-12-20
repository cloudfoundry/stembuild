package construct_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/stembuild/construct"

	"github.com/cloudfoundry-incubator/stembuild/version"

	"github.com/cloudfoundry-incubator/stembuild/construct/constructfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("OsVersionValidator", func() {
	var (
		validator        *construct.OSVersionValidator
		fakeGuestManager *constructfakes.FakeGuestManager
		fakeMessenger    *constructfakes.FakeOSValidatorMessenger
	)

	BeforeEach(func() {
		fakeGuestManager = &constructfakes.FakeGuestManager{}
		fakeMessenger = &constructfakes.FakeOSValidatorMessenger{}

		validator = &construct.OSVersionValidator{
			GuestManager: fakeGuestManager,
			Messenger:    fakeMessenger,
		}

		buildBuffer := gbytes.NewBuffer()
		_, err := buildBuffer.Write([]byte(version.VersionDev.Build))
		Expect(err).NotTo(HaveOccurred())

		fakeGuestManager.DownloadFileInGuestReturns(buildBuffer, 3, nil)
	})

	Describe("Validate version", func() {

		DescribeTable("permutations of stembuild and guest os versions",
			func(stembuildVersion string, guestOSVersion string, expectedMatch bool) {
				buildBuffer := gbytes.NewBuffer()
				_, err := buildBuffer.Write([]byte(guestOSVersion))
				Expect(err).NotTo(HaveOccurred())

				fakeGuestManager.DownloadFileInGuestReturns(buildBuffer, int64(len(guestOSVersion)), nil)

				err = validator.Validate(stembuildVersion)
				if expectedMatch {
					Expect(err).ToNot(HaveOccurred())
				} else {
					Expect(err).To(HaveOccurred())
				}
			},
			Entry("stembuild version 'dev' does NOT match guest os build number '17763'", "dev", version.Version2019.Build, false),
			Entry("stembuild version 'dev' does NOT match guest os version '1803'", "dev", version.Version1803.Build, false),
			Entry("stembuild version '2019.4.9' does NOT match guest os version 'dev'", "2019.4.9", version.VersionDev.Build, false),
			Entry("stembuild version '2019.4.9' does NOT match guest os version '1803'", "2019.4.9", version.Version1803.Build, false),
			Entry("stembuild version '1803.10.11' does NOT match guest os version 'dev'", "1803.10.11", version.VersionDev.Build, false),
			Entry("stembuild version '1803.10.11' does NOT match guest os version '2019'", "1803.10.11", version.Version2019.Build, false),

			Entry("stembuild version '2019.4' does NOT match guest os version 'dev'", "2019.4", version.VersionDev.Build, false),
			Entry("stembuild version '2019.4' does NOT match guest os version '1803'", "2019.4", version.Version1803.Build, false),
			Entry("stembuild version '1803.10' does NOT match guest os version 'dev'", "1803.10", version.VersionDev.Build, false),
			Entry("stembuild version '1803.10' does NOT match guest os version '2019'", "1803.10", version.Version2019.Build, false),

			Entry("stembuild version 'dev' does match guest os version 'dev", "dev", version.VersionDev.Build, true),
			Entry("stembuild version '2019.8.88' does match guest os version '2019'", "2019.8.88", version.Version2019.Build, true),
			Entry("stembuild version '1803.4.6' does match guest os version '1803'", "1803.4.6", version.Version1803.Build, true),
			Entry("stembuild version '2019.8' does match guest os version '2019'", "2019.8", version.Version2019.Build, true),
			Entry("stembuild version '1803.4' does match guest os version '1803'", "1803.4", version.Version1803.Build, true),
		)

		It("returns nil even if get OS version file creation fails", func() {
			fakePid := 123
			fakeGuestManager.StartProgramInGuestReturnsOnCall(0, int64(fakePid), errors.New("failed to create blah"))

			err := validator.Validate("fakeVersion")

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeGuestManager.DownloadFileInGuestCallCount()).To(Equal(0))
			Expect(fakeMessenger.OSVersionFileCreationFailedCallCount()).To(Equal(1))
		})

		It("returns nil if the exit code for OS version file creation process cannot be retrieved", func() {
			fakeExitCode := 123
			fakeGuestManager.ExitCodeForProgramInGuestReturnsOnCall(0, int32(fakeExitCode), errors.New("failed to get exit code for process"))

			err := validator.Validate("fakeVersion")

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeGuestManager.DownloadFileInGuestCallCount()).To(Equal(0))
			Expect(fakeMessenger.ExitCodeRetrievalFailedCallCount()).To(Equal(1))
		})

		It("returns nil if the exit code for OS version file creation process is non-zero", func() {
			fakeExitCode := 123
			fakeGuestManager.StartProgramInGuestReturns(123, nil)
			fakeGuestManager.ExitCodeForProgramInGuestReturnsOnCall(0, int32(fakeExitCode), nil)

			err := validator.Validate("fakeVersion")

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeGuestManager.DownloadFileInGuestCallCount()).To(Equal(0))
			Expect(fakeMessenger.OSVersionFileCreationFailedCallCount()).To(Equal(1))
		})

		It("returns nil if the os version file that was created cannot be downloaded", func() {
			fakeGuestManager.DownloadFileInGuestReturnsOnCall(0, nil, 0, errors.New("could not download"))

			err := validator.Validate("fakeVersion")

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeMessenger.DownloadFileFailedCallCount()).To(Equal(1))
		})

		It("removes non-alphanumeric characters from the build number", func() {

			buildBuffer := gbytes.NewBuffer()
			_, err := buildBuffer.Write([]byte("&&&&&17763@@"))
			Expect(err).NotTo(HaveOccurred())

			fakeGuestManager.DownloadFileInGuestReturns(buildBuffer, 50, nil)

			err = validator.Validate(version.Version2019.Name)
			Expect(err).NotTo(HaveOccurred())

		})
		It("removes non-alphanumeric characters from the dev build number", func() {

			buildBuffer := gbytes.NewBuffer()
			_, err := buildBuffer.Write([]byte("&&&&&dev@@"))
			Expect(err).NotTo(HaveOccurred())

			fakeGuestManager.DownloadFileInGuestReturns(buildBuffer, 50, nil)

			err = validator.Validate(version.VersionDev.Name)
			Expect(err).NotTo(HaveOccurred())

		})
	})

})

func funcName() {
	for _, a := range version.AllVersions {
		Expect(a).To(Equal(a))
	}
}
