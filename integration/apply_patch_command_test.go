package integration

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	"github.com/pivotal-cf-experimental/stembuild/helpers"
	"github.com/pivotal-cf-experimental/stembuild/stembuildoptions"
)

var _ = Describe("Apply Patch", func() {
	var manifestStruct stembuildoptions.StembuildOptions
	BeforeEach(func() {
		manifestStruct = stembuildoptions.StembuildOptions{}
	})

	Context("when valid manifest file", func() {
		var stemcellFilename string
		var manifestFilename string

		BeforeEach(func() {
			manifestStruct.Version = "1200.0"
			manifestStruct.VHDFile = "testdata/original.vhd"
			manifestStruct.PatchFile = "testdata/diff.patch"
		})

		JustBeforeEach(func() {
			manifestFile, err := ioutil.TempFile("", "")
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				Expect(manifestFile.Close()).To(Succeed())
			}()

			contents, err := helpers.StringFromManifest(helpers.ManifestTemplate, manifestStruct)
			Expect(err).NotTo(HaveOccurred())
			_, err = manifestFile.Write([]byte(contents))
			Expect(err).NotTo(HaveOccurred())

			manifestFilename = manifestFile.Name()
		})

		Context("when no output directory is specified on the command line", func() {
			BeforeEach(func() {
				osVersion := "2012R2"
				stemcellFilename = fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz", manifestStruct.Version, osVersion)
				manifestStruct.VHDFile = "testdata/original.vhd"
				manifestStruct.PatchFile = "testdata/diff.patch"
			})

			AfterEach(func() {
				Expect(os.Remove(stemcellFilename)).To(Succeed())
			})

			Context("current working directory has no stemcell tgz in it", func() {
				It("creates a stemcell in current working directory", func() {
					session := helpers.Stembuild("apply-patch", manifestFilename)
					Eventually(session, 5).Should(Exit(0))
					Eventually(session).Should(Say(`created stemcell: .*\.tgz`))
				})
			})

			Context("current working directory has stemcell tgz in it", func() {
				BeforeEach(func() {
					stemcellFile, err := os.Create(stemcellFilename)
					Expect(err).NotTo(HaveOccurred())
					stemcellFile.Close()
				})

				It("displays an error", func() {
					session := helpers.Stembuild("apply-patch", manifestFilename)
					Eventually(session).Should(Exit(1))
					Eventually(session.Err).Should(Say("file may already exist"))
					Eventually(session.Err).Should(Say(`\n\nfor usage: stembuild -h`))
				})
			})
		})

		Context("when output directory specified with -o flag", func() {
			var tmpDir string

			BeforeEach(func() {
				var err error
				tmpDir, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				Expect(os.RemoveAll(tmpDir)).To(Succeed())
			})

			Context("when manifest does not specify output directory", func() {
				Context("directory already exists", func() {
					It("creates stemcell in output directory", func() {
						session := helpers.Stembuild("-o", tmpDir, "apply-patch", manifestFilename)
						Eventually(session, 5).Should(Exit(0))
						Eventually(session).Should(Say(`created stemcell: .*%s.*\.tgz`, tmpDir))
					})
				})

				Context("directory does not exist", func() {
					AfterEach(func() {
						Expect(os.RemoveAll("idontexist")).To(Succeed())
					})
					It("creates directory and puts stemcell in it", func() {
						session := helpers.Stembuild("-o", "idontexist", "apply-patch", manifestFilename)
						Eventually(session, 5).Should(Exit(0))
						Eventually(session).Should(Say(`created stemcell: .*idontexist.*\.tgz`))
					})
				})
			})

			Context("when output directory specified only in manifest", func() {
				BeforeEach(func() {
					tmpDir, err := ioutil.TempDir("", "")
					Expect(err).NotTo(HaveOccurred())

					manifestStruct.OutputDir = tmpDir
				})

				AfterEach(func() {
					Expect(os.RemoveAll(manifestStruct.OutputDir)).To(Succeed())
				})

				It("creates stemcell in directory from manifest", func() {
					session := helpers.Stembuild("apply-patch", manifestFilename)
					Eventually(session, 5).Should(Exit(0))
					Eventually(session).Should(Say(`created stemcell: .*%s.*\.tgz`, manifestStruct.OutputDir))
				})
			})

			Context("when manifest does specify output directory", func() {
				It("creates stemcell in dir specified by -o flag", func() {
				})
			})
		})
	})
})

// func TestMissingOutputDirectoryCreatesDirectory(t *testing.T) {
// 	// Setup output directory
// 	testOutputDir, err := ioutil.TempDir("", "testOutputDir-")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	os.RemoveAll(testOutputDir)
// 	if helpers.Exists(testOutputDir) {
// 		t.Errorf("%s already exists, not a valid test", testOutputDir)
// 	}

// 	// Setup input vhd and vmdk
// 	testInputDir, err := ioutil.TempDir("", "testInputDir-")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.RemoveAll(testInputDir)
// 	testEmptyFilePath := filepath.Join(testInputDir, "testEmptyFile.txt")
// 	testEmptyFile, err := os.Create(testEmptyFilePath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	testEmptyFile.Close()

// 	testCommand := fmt.Sprintf(
// 		"stembuild -vhd %s -patch %s -v 1200.666 -output %s",
// 		testEmptyFilePath,
// 		testEmptyFilePath,
// 		testOutputDir,
// 	)
// 	testArgs := strings.Split(testCommand, " ")
// 	os.Args = testArgs
// 	runInit()
// 	ParseFlags()

// 	errs := ValidateFlags()

// 	if len(errs) != 0 {
// 		t.Errorf("expected no errors, but got errors: %s", errs)
// 	}

// 	if !helpers.Exists(testOutputDir) {
// 		t.Errorf("%s was not created", testOutputDir)
// 	}
// }
