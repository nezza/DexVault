// +build linux

package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"syscall"
)

func Mlock() {
	fmt.Println("Enabling mlock protection.")
	// TODO Check if successful
	unix.Mlockall(syscall.MCL_CURRENT | syscall.MCL_FUTURE)
}
