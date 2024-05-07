package main

import "runtime"

// is linux?
func isLinux() bool {
	return runtime.GOOS == "linux"
}
