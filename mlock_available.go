// +build linux

package main

import (
	"fmt"
	"syscall"
	"golang.org/x/sys/unix"
)

func Mlock() {
	fmt.Println("Enabling mlock protection.")
	// TODO Check if successful
	unix.Mlockall(syscall.MCL_CURRENT | syscall.MCL_FUTURE)
}