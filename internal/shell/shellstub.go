//go:build !darwin && !freebsd && !linux && !windows
// +build !darwin,!freebsd,!linux,!windows

package shell

import "log"

func OpenFolder(folder string) error {
	log.Println(folder)
	return nil
}
