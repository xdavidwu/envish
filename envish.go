package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"

	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

const (
	help = `Usage: %s [OPTION]... [--] [COMMAND [ARG]...]

Run COMMAND or a shell with envtest configured

`
)

var (
	version        *string
	envtestBinPath *bool
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), help, os.Args[0])
	flag.PrintDefaults()
}

func findShell() string {
	sh := os.Getenv("SHELL")
	if sh == "" {
		sh = currentUserShell()
		if sh == "" {
			sh = "/bin/sh"
		}
	}
	return sh
}

func envish() int {
	dir, err := envtest.SetupEnvtestDefaultBinaryAssetsDirectory()
	if err != nil {
		panic(err)
	}

	env := envtest.Environment{
		DownloadBinaryAssets:        true,
		DownloadBinaryAssetsVersion: *version,
		BinaryAssetsDirectory:       dir,
	}
	_, err = env.Start()
	if err != nil {
		panic(err)
	}
	defer env.Stop()

	if *envtestBinPath {
		envtestBin := path.Dir(env.ControlPlane.KubectlPath)
		os.Setenv("PATH", envtestBin+":"+os.Getenv("PATH"))
	}

	kcfg, err := os.CreateTemp("", "envish-kubeconfig-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(kcfg.Name())

	if _, err = kcfg.Write(env.KubeConfig); err != nil {
		panic(err)
	}

	args := flag.Args()
	if len(args) == 0 {
		args = append(args, findShell())
	}
	cmd := exec.Command(args[0], args[1:]...)
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
	// XXX controller-runtime registers -kubeconfig
	if len(os.Args) == 0 {
		flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	} else {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}
	version = flag.String("envtest-version", "", "envtest binaries version, defaults to latest stable")
	envtestBinPath = flag.Bool("envtest-bin-path", true, "Add envtest binaries directory to front of $PATH")
	flag.CommandLine.Usage = usage
	flag.Parse()

	os.Exit(envish())
}
