package construct

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/cloudfoundry-incubator/stembuild/version"
)

//go:generate counterfeiter . OSValidatorMessenger
type OSValidatorMessenger interface {
	OSVersionFileCreationFailed(errorMessage string)
	ExitCodeRetrievalFailed(errorMessage string)
	DownloadFileFailed(errorMessage string)
}

//go:generate counterfeiter . VersionGetter
type VersionGetter interface {
	GetVersion() string
}

type OSVersionValidator struct {
	GuestManager GuestManager
	Messenger    OSValidatorMessenger
}

func (v *OSVersionValidator) Validate(stembuildVersion string) error {
	pid, err := v.GuestManager.StartProgramInGuest(
		context.Background(),
		powershell,
		"[System.Environment]::OSVersion.Version.Build > C:\\Windows\\Temp\\version.log",
	)
	if err != nil {
		v.Messenger.OSVersionFileCreationFailed(err.Error())
		return nil
	}

	exitCode, err := v.GuestManager.ExitCodeForProgramInGuest(context.Background(), pid)
	if err != nil {
		v.Messenger.ExitCodeRetrievalFailed(err.Error())
		return nil
	}
	if exitCode != 0 {
		v.Messenger.OSVersionFileCreationFailed(fmt.Sprintf("OS version file creation failed with non-zero exit code: %d", exitCode))
		return nil
	}

	fileReader, _, err := v.GuestManager.DownloadFileInGuest(context.Background(), "C:\\Windows\\Temp\\version.log")
	if err != nil {
		v.Messenger.DownloadFileFailed(err.Error())
		return nil
	}
	buf, err := ioutil.ReadAll(fileReader)
	if err != nil {
		v.Messenger.DownloadFileFailed(err.Error())
		return nil
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	guestOSVersion := version.GetOSVersionFromBuildNumber(reg.ReplaceAllString(string(buf), ""))

	if !strings.Contains(stembuildVersion, guestOSVersion) {
		return fmt.Errorf("OS version of stembuild and guest OS VM do not match. Guest OS Version:'%s', Stembuild Version:'%s'", guestOSVersion, stembuildVersion)
	}

	return nil
}
