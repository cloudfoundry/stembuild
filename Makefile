GOSRC = $(shell find . -name "*.go" ! -name "*test.go" ! -name "*fake*" ! -path "./integration/*")
COMMAND = out/stembuild
AUTOMATION_PATH = StemcellAutomation.zip
AUTOMATION_PREFIX = $(shell dirname "${AUTOMATION_PATH}")
STEMCELL_VERSION = $(shell echo "$${STEMBUILD_VERSION}")
LD_FLAGS = "-w -s -X github.com/cloudfoundry-incubator/stembuild/version.Version=${STEMCELL_VERSION}"

all : test build

build : out/stembuild

clean :
	rm -r version/version.go || true
	rm -r $(wildcard out/*) || true
	rm -r assets/stemcell_automation.go || true

format :
	go fmt ./...

integration : generate
	ginkgo -r -v -randomizeAllSpecs integration

generate: $(GOSRC) $(AUTOMATION_PATH)
	go get -u github.com/jteeuwen/go-bindata/...
	go-bindata -o assets/stemcell_automation.go -pkg assets -prefix $(AUTOMATION_PREFIX) $(AUTOMATION_PATH)

out/stembuild : generate $(GOSRC)
	go build -o $(COMMAND) -ldflags $(LD_FLAGS) .

test : units

units : format generate
	@ginkgo version
	ginkgo -r -v -randomizeAllSpecs -randomizeSuites -skipPackage integration,iaas_cli
	@echo "\nSWEET SUITE SUCCESS"

contract :
	ginkgo -r -v -randomizeAllSpecs -randomizeSuites iaas_cli

.PHONY : all build clean format generate
.PHONY : test units units-full integration integration-tests-full
