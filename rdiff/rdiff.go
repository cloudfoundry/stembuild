package rdiff

// #cgo CFLAGS: -Wall -O2 -DNDEBUG -Wno-unused-function
//
// #include <stdlib.h>   // free
// #include <stdio.h>    // fopen, fclose
// #include <time.h>     // difftime
//
// #include "librsync.h" // rs_patch_file, rs_result, rs_strerror
//
import "C"

import (
	"errors"
	"fmt"
	"os"
	"unsafe"
)

type rsResult int

const (
	rs_DONE           rsResult = 0
	rs_BLOCKED        rsResult = 1
	rs_RUNNING        rsResult = 2
	rs_TEST_SKIPPED   rsResult = 77
	rs_IO_ERROR       rsResult = 100
	rs_SYNTAX_ERROR   rsResult = 101
	rs_MEM_ERROR      rsResult = 102
	rs_INPUT_ENDED    rsResult = 103
	rs_BAD_MAGIC      rsResult = 104
	rs_UNIMPLEMENTED  rsResult = 105
	rs_CORRUPT        rsResult = 106
	rs_INTERNAL_ERROR rsResult = 107
	rs_PARAM_ERROR    rsResult = 108
)

var rsResultStr = map[rsResult]string{
	rs_DONE:           "RS_DONE",
	rs_BLOCKED:        "RS_BLOCKED",
	rs_RUNNING:        "RS_RUNNING",
	rs_TEST_SKIPPED:   "RS_TEST_SKIPPED",
	rs_IO_ERROR:       "RS_IO_ERROR",
	rs_SYNTAX_ERROR:   "RS_SYNTAX_ERROR",
	rs_MEM_ERROR:      "RS_MEM_ERROR",
	rs_INPUT_ENDED:    "RS_INPUT_ENDED",
	rs_BAD_MAGIC:      "RS_BAD_MAGIC",
	rs_UNIMPLEMENTED:  "RS_UNIMPLEMENTED",
	rs_CORRUPT:        "RS_CORRUPT",
	rs_INTERNAL_ERROR: "RS_INTERNAL_ERROR",
	rs_PARAM_ERROR:    "RS_PARAM_ERROR",
}

func (r rsResult) String() string {
	if s, ok := rsResultStr[r]; ok {
		return s
	}
	return fmt.Sprintf("INVALID RS_RESULT: %d", r)
}

func rsError(r rsResult) error {
	if r == rs_DONE {
		return nil
	}
	if s, ok := rsResultStr[r]; ok {
		msg := C.GoString(C.rs_strerror(C.rs_result(r)))
		return errors.New(s + ": " + msg)
	}
	return fmt.Errorf("Invalid RS_RESULT: %d", r)
}

func fclose(f *C.FILE) {
	if f != nil {
		C.fclose(f)
	}
}

func fopen(filename, mode string) (*C.FILE, error) {
	s := C.CString(filename)
	defer C.free(unsafe.Pointer(s))

	m := C.CString(mode)
	defer C.free(unsafe.Pointer(m))

	// Do not check 'err != nil', instead check if f is nil.
	// CGO sets err to the value of errno, which may contain
	// an error from a previous call.
	f, err := C.fopen(s, m)
	if f == nil {
		return nil, &os.PathError{
			Op:   "fopen",
			Path: filename,
			Err:  err,
		}
	}
	return f, nil
}

func Patch(basis, delta, newfile string) error {
	fbasis, err := fopen(basis, "rb")
	if err != nil {
		return err
	}
	defer fclose(fbasis)

	fdelta, err := fopen(delta, "rb")
	if err != nil {
		return err
	}
	defer fclose(fdelta)

	// the 'x' flag is not supported on Windows
	// so manually check if the file exists
	if _, err := os.Stat(newfile); err == nil {
		return &os.PathError{
			Op:   "fopen",
			Path: newfile,
			Err:  os.ErrExist,
		}
	}
	fnew, err := fopen(newfile, "wb")
	if err != nil {
		return err
	}
	defer fclose(fnew)

	result := rsResult(C.rs_patch_file(fbasis, fdelta, fnew, nil))
	if result != rs_DONE {
		return rsError(result)
	}
	return nil
}
