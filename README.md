# envish

Run commands with temporary kube-apiserver + etcd via [envtest](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/envtest)

Also update $PATH to include envtest binaries (e.g. `kubectl`)

## Usage

```
Usage: envish [OPTION]... [--] [COMMAND [ARG]...]

Run COMMAND or a shell with envtest configured

  -envtest-bin-path
    	Add envtest binaries directory to front of $PATH (default true)
  -envtest-version string
    	envtest binaries version, defaults to latest stable
```

## Examples

```sh-session
$ # switching versions
$ # $PATH is automatically set so that kubectl from envtest is used
$ envish kubectl version
Client Version: v1.35.0
Kustomize Version: v5.7.1
Server Version: v1.35.0
$ envish -envtest-version v1.34.1 -- kubectl version
Client Version: v1.34.1
Kustomize Version: v5.7.1
Server Version: v1.34.1
$ envish -envtest-version v1.33.0 -- kubectl version
Client Version: v1.33.0
Kustomize Version: v5.6.0
Server Version: v1.33.0
$
$ # shell
$ env PS1='% ' envish
% # $KUBECONFIG is configured to connect to envtest
% # beside kubectl, other tools respecting it also work
% echo $KUBECONFIG
/tmp/envish-kubeconfig-106267163
% flux check --pre
► checking prerequisites
✗ flux 2.7.3 <2.7.5 (new CLI version is available, please upgrade)
✔ Kubernetes 1.35.0 >=1.32.0-0
✔ prerequisites checks passed
% exit
$
```

## FAQs

- What versions of Kubernetes are available?

See [setup-envtest](https://github.com/kubernetes-sigs/controller-runtime/tree/main/tools/setup-envtest)
