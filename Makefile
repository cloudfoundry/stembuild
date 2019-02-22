LD_FLAGS = "-w -s"
GOSRC = $(shell find . -name "*.go" ! -name "*test.go" ! -name "*fake*" ! -path "./integration/*")
COMMAND = out/stembuild

all : test build

build : out/stembuild

clean :
	rm -r version/version.go
	rm -r $(wildcard out/*)

format :
	go fmt ./...

integration : build
	PATH="$(PWD)/out:$(PATH)" ginkgo -r -v -randomizeAllSpecs integration

out/stembuild : $(GOSRC)
	go generate
	go build -o $(COMMAND) -ldflags $(LD_FLAGS) .

test : units

units : format build
	@ginkgo version
	PATH="$(PWD)/out:$(PATH)" ginkgo -r -v -randomizeAllSpecs -randomizeSuites -skipPackage integration,iaas_cli
	@echo "\nSWEET SUITE SUCCESS"

contract :
	ginkgo -r -v -randomizeAllSpecs -randomizeSuites iaas_cli

.PHONY : all build clean format
.PHONY : test units units-full integration integration-tests-full
