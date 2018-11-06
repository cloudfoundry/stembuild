package stemcell_test

import (
	"crypto/sha1"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pivotal-cf-experimental/stembuild/pack/options"
	"github.com/pivotal-cf-experimental/stembuild/pack/stemcell"
	"github.com/pivotal-cf-experimental/stembuild/test/helpers"
)

var _ = Describe("Stemcell", func() {
	var tmpDir string
	var stembuildConfig options.StembuildOptions
	var c stemcell.Config

	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		stembuildConfig = options.StembuildOptions{
			OSVersion: "2012R2",
			Version:   "1200.1",
		}

		c = stemcell.Config{
			Stop:         make(chan struct{}),
			Debugf:       func(format string, a ...interface{}) {},
			BuildOptions: stembuildConfig,
		}
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	Describe("CreateImage", func() {

		It("successfully creates an image tarball", func() {
			c.BuildOptions.VMDKFile = filepath.Join("..", "..", "test", "data", "expected.vmdk")
			err := c.CreateImage()
			Expect(err).NotTo(HaveOccurred())

			// the image will be saved to the Config's temp directory
			tmpdir, err := c.TempDir()
			Expect(err).NotTo(HaveOccurred())

			outputImagePath := filepath.Join(tmpdir, "image")
			Expect(c.Image).To(Equal(outputImagePath))

			// Make sure the sha1 sum is correct
			h := sha1.New()
			f, err := os.Open(c.Image)
			Expect(err).NotTo(HaveOccurred())

			_, err = io.Copy(h, f)
			Expect(err).NotTo(HaveOccurred())

			actualShasum := fmt.Sprintf("%x", h.Sum(nil))
			Expect(c.Sha1sum).To(Equal(actualShasum))

			// expect the image ova to contain only the following file names
			expectedNames := []string{
				"image.ovf",
				"image.mf",
				"image-disk1.vmdk",
			}

			imageDir, err := helpers.ExtractGzipArchive(c.Image)
			Expect(err).NotTo(HaveOccurred())
			list, err := ioutil.ReadDir(imageDir)
			Expect(err).NotTo(HaveOccurred())

			var names []string
			infos := make(map[string]os.FileInfo)
			for _, fi := range list {
				names = append(names, fi.Name())
				infos[fi.Name()] = fi
			}
			Expect(names).To(ConsistOf(expectedNames))

			// the vmx template should generate an ovf file that
			// does not contain an ethernet section.
			ovf := filepath.Join(imageDir, "image.ovf")
			ovfFile, err := helpers.ReadFile(ovf)
			Expect(err).NotTo(HaveOccurred())
			Expect(ovfFile).NotTo(MatchRegexp(`(?i)ethernet`))
		})
	})

	Describe("CreateManifest", func() {
		It("Creates a manifest correctly", func() {
			expectedManifest := `---
name: bosh-vsphere-esxi-windows1-go_agent
version: 'version'
sha1: sha1sum
operating_system: windows1
cloud_properties:
  infrastructure: vsphere
  hypervisor: esxi
stemcell_formats:
- vsphere-ovf
- vsphere-ova
`
			result := stemcell.CreateManifest("1", "version", "sha1sum")
			Expect(result).To(Equal(expectedManifest))
		})
	})

	Describe("catchInterruptSignal", func() {

		It("cleans up on one interrupt", func() {
			if runtime.GOOS == "windows" {
				Skip("Skipping, test not supported on Windows.")
			}

			inputVmdk := filepath.Join("..", "..", "test", "data", "expected.vmdk")
			session := helpers.Stembuild("package", "--vmdk", inputVmdk, "--os", "2012R2", "--version", "1200.0", "--outputDir", tmpDir)
			time.Sleep(1 * time.Second)

			err := session.Command.Process.Signal(os.Interrupt)
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(1 * time.Second)

			stdErr := session.Err.Contents()
			Expect(string(stdErr)).To(ContainSubstring("received ("))
		})

		// Tried to create test to handle 2 interrupts in a row, but timing of processes makes it difficult
		// to test
	})
})
