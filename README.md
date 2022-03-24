# k8scr

A `kubectl` plugin for pushing OCI images through the Kubernetes API server.

## Quickstart


<p align="center">
  <img src="docs/media/k8scr.svg">
</p>

1. Build `k8scr`

```
make build
```

2. Move to location in `PATH`

```
sudo mv ./build/k8scr /usr/local/bin/kubectl-k8scr
```

3. Deploy simple in-memory registry into cluster

```
kubectl apply -f distribution.yaml
```

> Optional: tail logs to observe results of next step with `kubectl logs k8scr -f`.

4. Push image to registry

```
kubectl k8scr push crossplane/crossplane:v1.2.1
```

## Usage

```
Usage: k8scr <command>

Push and pull images through the Kubernetes API server.

Flags:
  -h, --help                   Show context-sensitive help.
      --kubeconfig=STRING      Override default kubeconfig path.
  -n, --namespace="default"    Namespace of registry Pod.
  -r, --registry="k8scr"       Name of registry Pod.

Commands:
  push <image>

  pull <image>
```

## How Does This Work?

`k8scr` uses
[`go-containerregistry`](https://github.com/google/go-containerregistry) to push
and pull images, but passes in an
[`http.RoundTripper`](https://golang.org/pkg/net/http/#RoundTripper) that
reconstructs [OCI
distribution](https://github.com/opencontainers/distribution-spec/blob/main/spec.md)
compliant requests so that they pass through the Kubernetes API server `Pod`
proxy endpoint, before eventually calling the underlying transport constructed
from a user's `kubeconfig`. This allows for pushing and pulling directly to and
from an OCI image registry running in a Kubernetes cluster without having to
expose it publicly or privately. Any user with access to the cluster and
`pods/proxy` RBAC permissions for the registry `Pod` is able to push and pull.

## What Else Can It Do?

Pretty much any of the operations
[`go-containerregistry`](https://github.com/google/go-containerregistry)
supports could also be supported here as the transport is pluggable. I'll likely
move it upstream or offer it as a stand-alone library if there is enough
interest.
