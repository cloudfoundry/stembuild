//go:build tools
// +build tools

package main

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
	_ "github.com/onsi/ginkgo/v2/ginkgo"
)
