package construct

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-incubator/stembuild/assets"
)

//go:generate counterfeiter . zipUnarchiver
type zipUnarchiver interface {
	Unzip(fileArchive []byte, file string) ([]byte, error)
}

type WinRMManager struct {
	GuestManager GuestManager
	Unarchiver   zipUnarchiver
}

func (w *WinRMManager) Enable() error {
	failureString := "failed to enable WinRM: %s"
	saZip, err := assets.Asset(stemcellAutomationName)
	if err != nil {
		return fmt.Errorf(failureString, err)
	}

	bmZip, err := w.Unarchiver.Unzip(saZip, boshPsModules)
	if err != nil {
		return fmt.Errorf(failureString, err)
	}

	rawWinRM, err := w.Unarchiver.Unzip(bmZip, winRMPsScript)
	if err != nil {
		return fmt.Errorf(failureString, err)
	}

	rawWinRMwtCmd := append(rawWinRM, []byte("\nEnable-WinRM\n")...)

	base64WinRM := encodePowershellCommand(rawWinRMwtCmd)

	pid, err := w.GuestManager.StartProgramInGuest(context.Background(), powershell, fmt.Sprintf("-EncodedCommand %s", base64WinRM))
	if err != nil {
		return fmt.Errorf(failureString, err)
	}

	exitCode, err := w.GuestManager.ExitCodeForProgramInGuest(context.Background(), pid)
	if err != nil {
		return fmt.Errorf(failureString, err)
	}

	if exitCode != 0 {
		return fmt.Errorf(failureString, fmt.Sprintf("WinRM process on guest VM exited with code %d", exitCode))
	}

	return nil
}
