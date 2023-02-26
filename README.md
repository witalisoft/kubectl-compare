# kubectl-compare

Simple tool to compare two Kubernetes objects and output them as a diff. It's using a more readable Kubernetes object
manifest by passing it through [kubectl-neat](https://github.com/itaysk/kubectl-neat).

## Usage

```console
$ kubectl-compare compare [k8s_object] [object1] [object2]
```

## Example

```console
$ kubectl-compare compare pods --namespace test pod1 pod2
  apiVersion: "v1"
  kind: "Pod"
  metadata:
    annotations:
      kubernetes.io/psp: "eks.privileged"
    labels:
      app: "test"
-   name: "pod1"
+   name: "pod2"
...
```

## TODO

- expose it as a kubectl plugin
- add Windows support
