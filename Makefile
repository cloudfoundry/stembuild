GOSRC = $(shell find . -name "*.go" ! -name "*test.go" ! -name "*fake*" ! -path "./integration/*")
FAKE_STEMCELL_AUTOMATION_PATH = integration/construct/assets/StemcellAutomation.zip
FAKE_STEMCELL_AUTOMATION_PREFIX = $(shell dirname "${FAKE_STEMCELL_AUTOMATION_PATH}")
STEMCELL_VERSION = $(shell echo "$${STEMBUILD_VERSION}")
LD_FLAGS = "-w -s -X github.com/cloudfoundry/stembuild/version.Version=${STEMCELL_VERSION}"

# These are the sources for StemcellAutomation.zip
STEMCELL_AUTOMATION_PS1 := $(shell ls stemcell-automation/*ps1 | grep -iv Test)
BOSH_AGENT_REPO ?= ${HOME}/go/src/github.com/cloudfoundry/bosh-agent
LGPO_URL = 'https://download.microsoft.com/download/8/5/C/85C25433-A1B0-4FFA-9429-7E023E7DA8D8/LGPO.zip'
BOSH_GCS_URL = 'https://s3.amazonaws.com/bosh-gcscli/bosh-gcscli-0.0.6-windows-amd64.exe'
# Ignore things under cis-merge* directory because the paths contain spaces and make doesn't like
# that
PSMODULES_SOURCES = $(shell find ./modules | grep -v .git | grep -vi "test" | grep -v cis-merge)
BOSH_AGENT_SOURCES = $(shell find $(BOSH_AGENT_REPO) | egrep -v ".git|test.go|fake|.md")

ifeq ($(OS),Windows_NT)
	COMMAND = out/stembuild.exe
	CP = cp
else
	UNAME = $(shell uname -s)
	COMMAND = out/stembuild

	ifeq ($(UNAME),Darwin)
		CP = cp -p
	else ifeq ($(UNAME),Linux)
		CP = cp --preserve=mode,ownership
	endif
endif

all : test build

build : out/stembuild

build-integration : generate-fake-stemcell-automation $(GOSRC)
	go build -o $(COMMAND) -ldflags $(LD_FLAGS) .

clean :
	rm -rf version/version.go assets/stemcell_automation.go assets/StemcellAutomation.zip assets/local/* out/*

format :
	go fmt ./...

integration : generate-fake-stemcell-automation
	go run github.com/onsi/ginkgo/v2/ginkgo -r -vv --randomize-all --keep-going --flake-attempts 2 --timeout 3h integration

integration/construct : generate-fake-stemcell-automation
	go run github.com/onsi/ginkgo/v2/ginkgo -r -v --randomize-all --keep-going --flake-attempts 2 --timeout 3h integration/construct

integration-badger : generate-fake-stemcell-automation
	go run github.com/onsi/ginkgo/v2/ginkgo -r -v --randomize-all --until-it-fails --timeout 3h integration

generate-fake-stemcell-automation: $(GOSRC) $(FAKE_STEMCELL_AUTOMATION_PATH)
	go run github.com/go-bindata/go-bindata/go-bindata -o assets/stemcell_automation.go -pkg assets -prefix $(FAKE_STEMCELL_AUTOMATION_PREFIX) $(FAKE_STEMCELL_AUTOMATION_PATH)

generate: assets/StemcellAutomation.zip
	go run github.com/go-bindata/go-bindata/go-bindata -o assets/stemcell_automation.go -pkg assets -prefix assets assets/StemcellAutomation.zip

out/stembuild : generate $(GOSRC)
	go build -o $(COMMAND) -ldflags $(LD_FLAGS) .

out/stembuild.exe : generate $(GOSRC)
	GOOS=windows go build -o out/stembuild.exe -ldflags $(LD_FLAGS) .

test : units

units : format generate-fake-stemcell-automation
	@go run github.com/onsi/ginkgo/v2/ginkgo version
	go run github.com/onsi/ginkgo/v2/ginkgo -r --randomize-all --randomize-suites --keep-going --skip-package integration,iaas_cli
	@echo "\nSWEET SUITE SUCCESS"

contract :
	go run github.com/onsi/ginkgo/v2/ginkgo -r --randomize-all --randomize-suites --keep-going --flake-attempts 2 iaas_cli

.PHONY : all build build-integration clean format generate generate-fake-stemcell-automation
.PHONY : test units units-full integration integration-tests-full

# ===============================================================================
# The following to create the StemcellAutomation.zip that's packaged in stembuild
# ===============================================================================

assets/local/bosh-agent.exe: $(BOSH_AGENT_SOURCES)
	cd $(BOSH_AGENT_REPO) && \
		GOOS=windows GOARCH=amd64 bin/build && \
		cd -
	mv $(BOSH_AGENT_REPO)/out/bosh-agent assets/local/bosh-agent.exe

assets/local/bosh-blobstore-dav.exe: $(BOSH_AGENT_REPO)/integration/windows/fixtures/bosh-blobstore-dav.exe
	@echo "### Creating/Updating assets/local/bosh-blobstore-dav.exe"
	$(CP) $(BOSH_AGENT_REPO)/integration/windows/fixtures/bosh-blobstore-dav.exe assets/local

assets/local/bosh-blobstore-gcs.exe:
	@echo "### Creating assets/local/bosh-blobstore-gcs.exe"
	curl -o assets/local/bosh-blobstore-gcs.exe -L $(BOSH_GCS_URL)

assets/local/bosh-blobstore-s3.exe: $(BOSH_AGENT_REPO)/integration/windows/fixtures/bosh-blobstore-s3.exe
	@echo "### Creating/Updating assets/local/bosh-blobstore-s3.exe"
	$(CP) $(BOSH_AGENT_REPO)/integration/windows/fixtures/bosh-blobstore-s3.exe assets/local

assets/local/bosh-psmodules.zip: $(PSMODULES_SOURCES)
	@echo "### Creating/Updating assets/local/bosh-psmodules.zip"
	 cd modules && zip -r ../bosh-psmodules.zip . && cd ..
	 mv bosh-psmodules.zip assets/local/bosh-psmodules.zip

assets/local/job-service-wrapper.exe: $(BOSH_AGENT_REPO)/integration/windows/fixtures/job-service-wrapper.exe
	@echo "### Creating/Updating assets/local/job-service-wrapper.exe"
	$(CP) $(BOSH_AGENT_REPO)/integration/windows/fixtures/job-service-wrapper.exe assets/local

assets/local/pipe.exe: $(BOSH_AGENT_SOURCES)
	cd $(BOSH_AGENT_REPO) && \
		GOOS=windows GOARCH=amd64 bin/build && \
		cd -
	mv $(BOSH_AGENT_REPO)/out/bosh-agent-pipe assets/local/pipe.exe

assets/local/service_wrapper.exe: $(BOSH_AGENT_REPO)/integration/windows/fixtures/service_wrapper.exe
	@echo "### Creating/Updating assets/local/service_wrapper.exe"
	$(CP) $(BOSH_AGENT_REPO)/integration/windows/fixtures/service_wrapper.exe assets/local

assets/local/service_wrapper.xml: $(BOSH_AGENT_REPO)/integration/windows/fixtures/service_wrapper.xml
	@echo "### Creating/Updating assets/local/service_wrapper.xml"
	$(CP) $(BOSH_AGENT_REPO)/integration/windows/fixtures/service_wrapper.xml assets/local

assets/local/tar.exe: $(BOSH_AGENT_REPO)/integration/windows/fixtures/tar.exe
	@echo "### Creating/Updating assets/local/tar.exe"
	$(CP) $(BOSH_AGENT_REPO)/integration/windows/fixtures/tar.exe assets/local

assets/local/agent.zip: assets/local/bosh-agent.exe assets/local/pipe.exe assets/local/service_wrapper.xml assets/local/service_wrapper.exe assets/local/bosh-blobstore-dav.exe assets/local/bosh-blobstore-gcs.exe assets/local/bosh-blobstore-s3.exe assets/local/job-service-wrapper.exe assets/local/tar.exe
	@echo "### Creating/Updating assets/local/agent.zip"
	mkdir -p assets/temp/deps
	$(CP) assets/local/service_wrapper.exe \
		assets/local/service_wrapper.xml \
		assets/local/bosh-agent.exe \
		assets/temp
	$(CP) assets/local/bosh-blobstore-dav.exe \
		assets/local/bosh-blobstore-gcs.exe \
		assets/local/bosh-blobstore-s3.exe \
		assets/local/job-service-wrapper.exe \
		assets/local/pipe.exe \
		assets/local/tar.exe \
		assets/temp/deps
	cd assets/temp && zip -r ../local/agent.zip * && cd -
	rm -rf assets/temp

assets/local/LGPO.zip:
	@echo "### Creating assets/local/LGPO.zip"
	curl -o assets/local/LGPO.zip -L $(LGPO_URL)

assets/local/OpenSSH-Win64.zip: $(BOSH_AGENT_REPO)/integration/windows/fixtures/OpenSSH-Win64.zip
	@echo "### Creating/Updating assets/local/OpenSSH-Win64.zip"
	$(CP) $(BOSH_AGENT_REPO)/integration/windows/fixtures/OpenSSH-Win64.zip assets/local

assets/local/deps.json: assets/local/agent.zip assets/local/bosh-psmodules.zip assets/local/LGPO.zip assets/local/OpenSSH-Win64.zip
	@echo "### Creating/Updating assets/local/deps.json"
	@#Note: The order of the following matters, change the script before changing these
	stemcell-automation/generate-dep-json.bash \
		assets/local/OpenSSH-Win64.zip \
		assets/local/bosh-psmodules.zip \
		assets/local/agent.zip \
		assets/local/LGPO.zip \
		> assets/local/deps.json

assets/StemcellAutomation.zip: $(STEMCELL_AUTOMATION_PS1) assets/local/OpenSSH-Win64.zip assets/local/bosh-psmodules.zip assets/local/deps.json assets/local/agent.zip
	@echo "### Creating/Updating assets/StemcellAutomation.zip"
	mkdir -p assets/temp
	cp -a $(STEMCELL_AUTOMATION_PS1) \
		assets/local/OpenSSH-Win64.zip \
		assets/local/bosh-psmodules.zip \
		assets/local/deps.json \
		assets/local/agent.zip \
		assets/temp
	cd assets/temp && zip ../StemcellAutomation.zip * && cd -
	rm -rf assets/temp
