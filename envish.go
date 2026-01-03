package main

// #include <unistd.h>
// #include <pwd.h>
import "C"

import (
	"fmt"
	"os"
	"os/exec"

	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func currentUserShell() string {
	return C.GoString(C.getpwuid(C.getuid()).pw_shell)
}

func main() {
	dir, err := envtest.SetupEnvtestDefaultBinaryAssetsDirectory()
	if err != nil {
		panic(err)
	}

	env := envtest.Environment{
		DownloadBinaryAssets:  true,
		BinaryAssetsDirectory: dir,
	}
	_, err = env.Start()
	if err != nil {
		panic(err)
	}

	defer env.Stop()
	kcfg, err := os.CreateTemp("", "envish-kubeconfig-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(kcfg.Name())

	if _, err = kcfg.Write(env.KubeConfig); err != nil {
		panic(err)
	}

	cmd := exec.Command(currentUserShell())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", "KUBECONFIG", kcfg.Name()))
	cmd.Run()
}
