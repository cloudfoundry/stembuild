package main

// #cgo CFLAGS: -std=c11 -Wall -Wextra -Wpedantic -Wshadow -Wno-unused-parameter -O2
// #cgo LDFLAGS: -L/usr/local//Cellar/librsync/2.0.0_1/include -L/usr/local/Cellar/librsync/2.0.0_1/lib -lrsync
//
// #include <stdlib.h>
// #include <stdio.h>
// #include <string.h>
//
// #include <librsync.h>
//
// static rs_result go_patch(const char *basis_name, const char *delta_name,
//                           const char *new_name) {
//
// 	FILE *basis_file = fopen(basis_name, "rb");
// 	FILE *delta_file = fopen(delta_name, "rb");
// 	FILE *new_file = fopen(new_name, "wb");
//
// 	rs_stats_t stats;
// 	rs_result result = rs_patch_file(basis_file, delta_file, new_file, &stats);
//
// 	fclose(new_file);
// 	fclose(delta_file);
// 	fclose(basis_file);
//
// 	return result;
// }
import "C"

import (
	"errors"
	"fmt"
)

var NoTargetSumError = errors.New("Checksum request but missing target hash.")
var HashNoMatchError = errors.New("Final data hash does not match.")

type RsResult int

const (
	RS_DONE           RsResult = 0
	RS_BLOCKED        RsResult = 1
	RS_RUNNING        RsResult = 2
	RS_TEST_SKIPPED   RsResult = 77
	RS_IO_ERROR       RsResult = 100
	RS_SYNTAX_ERROR   RsResult = 101
	RS_MEM_ERROR      RsResult = 102
	RS_INPUT_ENDED    RsResult = 103
	RS_BAD_MAGIC      RsResult = 104
	RS_UNIMPLEMENTED  RsResult = 105
	RS_CORRUPT        RsResult = 106
	RS_INTERNAL_ERROR RsResult = 107
	RS_PARAM_ERROR    RsResult = 108
)

func (r RsResult) String() string {
	switch r {
	case RS_DONE:
		return "RS_DONE"
	case RS_BLOCKED:
		return "RS_BLOCKED"
	case RS_RUNNING:
		return "RS_RUNNING"
	case RS_TEST_SKIPPED:
		return "RS_TEST_SKIPPED"
	case RS_IO_ERROR:
		return "RS_IO_ERROR"
	case RS_SYNTAX_ERROR:
		return "RS_SYNTAX_ERROR"
	case RS_MEM_ERROR:
		return "RS_MEM_ERROR"
	case RS_INPUT_ENDED:
		return "RS_INPUT_ENDED"
	case RS_BAD_MAGIC:
		return "RS_BAD_MAGIC"
	case RS_UNIMPLEMENTED:
		return "RS_UNIMPLEMENTED"
	case RS_CORRUPT:
		return "RS_CORRUPT"
	case RS_INTERNAL_ERROR:
		return "RS_INTERNAL_ERROR"
	case RS_PARAM_ERROR:
		return "RS_PARAM_ERROR"
	}
	return "INVALID"
}

func Patch(basis, delta, newfile string) error {
	cbasis := C.CString(basis)
	cdelta := C.CString(delta)
	cnew := C.CString(newfile)

	result := RsResult(C.go_patch(cbasis, cdelta, cnew))
	if result == 0 {
		return nil
	}
	return fmt.Errorf("patch (basis: %s delta: %s newfile: %s: error - %s",
		basis, delta, newfile, result.String())
}
