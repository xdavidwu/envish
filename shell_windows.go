//go:build windows

package main

import (
	"path/filepath"

	"golang.org/x/sys/windows"
)

func defaultShell() string {
	dir, err := windows.GetSystemDirectory()
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, "cmd.exe")
}
