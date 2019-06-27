package guest_manager_test

import (
	"context"
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/stembuild/iaas_cli/iaas_clients/guest_manager/guest_managerfakes"

	"github.com/cloudfoundry-incubator/stembuild/iaas_cli/iaas_clients/guest_manager"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware/govmomi/vim25/types"
)

var _ = Describe("GuestManager", func() {
	var (
		auth types.NamePasswordAuthentication
		ctx  context.Context
	)

	BeforeEach(func() {
		ctx = context.TODO()
		auth = types.NamePasswordAuthentication{}
	})

	Describe("StartProgramInGuest", func() {
		It("runs the command on the guest", func() {

			expectedPid := int64(600)

			procManager := guest_managerfakes.FakeProcManager{}
			procManager.StartProgramReturns(expectedPid, nil)

			guestManager := guest_manager.NewGuestManager(auth, &procManager)

			pid, err := guestManager.StartProgramInGuest(ctx, "mkdir", "C:\\dummy")
			Expect(err).NotTo(HaveOccurred())

			Expect(pid).To(Equal(expectedPid))

		})

		It("returns an error if StartProgram does", func() {

			procManager := guest_managerfakes.FakeProcManager{}
			procManager.StartProgramReturns(int64(0), errors.New("You aint nothin but a hound dog"))

			guestManager := guest_manager.NewGuestManager(auth, &procManager)

			_, err := guestManager.StartProgramInGuest(ctx, "mkdir", "C:\\dummy")
			Expect(err).To(MatchError("vcenter_client - could not run process: mkdir C:\\dummy on guest os, error: You aint nothin but a hound dog"))
		})
	})

	Describe("ExitCodeForProgramInGuest", func() {
		It("obtains the exit code when the given program, being run on the guest os, exits", func() {

			expectedTime := time.Now()
			expectedExitCode := int32(0)
			processInfo := types.GuestProcessInfo{
				ExitCode: expectedExitCode,
				EndTime:  &expectedTime,
			}

			procManager := guest_managerfakes.FakeProcManager{}
			procManager.ListProcessesReturns([]types.GuestProcessInfo{processInfo}, nil)

			guestManager := guest_manager.NewGuestManager(auth, &procManager)

			exitCode, err := guestManager.ExitCodeForProgramInGuest(ctx, 1000)
			Expect(err).NotTo(HaveOccurred())
			Expect(exitCode).To(Equal(expectedExitCode))
		})

		It("returns an error if ListProcesses does", func() {
			procManager := guest_managerfakes.FakeProcManager{}
			procManager.ListProcessesReturns(nil, errors.New("yo"))

			guestManager := guest_manager.NewGuestManager(auth, &procManager)

			_, err := guestManager.ExitCodeForProgramInGuest(ctx, 1000)
			Expect(err).To(MatchError("vcenter_client - could not observe program exiting: yo"))
		})

		It("returns an error if ListProcesses does not find pid", func() {
			procManager := guest_managerfakes.FakeProcManager{}
			procManager.ListProcessesReturns([]types.GuestProcessInfo{}, nil)

			guestManager := guest_manager.NewGuestManager(auth, &procManager)

			_, err := guestManager.ExitCodeForProgramInGuest(ctx, 1000)
			Expect(err).To(MatchError("vcenter_client - could not observe program exiting"))
		})
	})
})
