//go:build cgo

package main

// #include <unistd.h>
// #include <pwd.h>
import "C"

func currentUserShell() string {
	passwd := C.getpwuid(C.getuid())
	if passwd == nil {
		return ""
	}
	return C.GoString(passwd.pw_shell)
}
