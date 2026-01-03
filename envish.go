package main

// #include <unistd.h>
// #include <pwd.h>
import "C"

import (
	"fmt"
	"os"
	"os/exec"

	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func currentUserShell() string {
	return C.GoString(C.getpwuid(C.getuid()).pw_shell)
}

func envish() int {
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
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", clientcmd.RecommendedConfigPathEnvVar, kcfg.Name()))

	status := 0
	err = cmd.Run()
	if err != nil {
		eerr, ok := err.(*exec.ExitError)
		if !ok {
			panic(err)
		}
		status = eerr.ExitCode()
	}
	return status
}

func main() {
	os.Exit(envish())
}
