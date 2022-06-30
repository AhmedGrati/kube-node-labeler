# kube-node-labeler

Kube Node Labeler is a Kubernetes operator which aims to automate the management of Kubernetes nodes for administrators in terms of attributes like Labels, Taint, and Annotations.

To read more about Kubernetes operators you could read this amazing <a href="https://github.com/cncf/tag-app-delivery/blob/main/operator-wg/whitepaper/Operator-WhitePaper_v1-0.md">white paper</a> written by the CNCF (Cloud Native Computing Foundation).

## Use Case

When a node is down you lose all the attached Attributes and NodeSpec, so it is frustrating for Kubernetes administrators to re-configure this node. You could generalize it for a 1xx-nodes cluster :').

## Features

- [x] [Adding attributes](#adding-attributes)
- [x] [Overwrite attributes](#overwrite-attributes)

Attributes = Labels, Taints, Annotations

## Adding Attributes

This operator gives the ability to add attributes to one or more nodes that could be selected using RegEx and/or NodeSelectorTerms, and/or specific size of nodes (if it is 0 it would select all of them).

Since you could only have only unique keys in labels and annotations, adding an existing key would keep the same key with the same old value.

For taints, It would be added normally even though a key already exists.

I am open to change the behaviour of adding labels if you have any useful suggestions.

Example:

```yaml
apiVersion: kubebuilder.kube.node.labeler.io/v1alpha1
kind: NodeLabeler
metadata:
  annotations:
    annotation1: value1
  generation: 1
  labels:
    type: unit-test
  name: simple-node-labeler
  namespace: default
spec:
  size: 3
  nodeNamePatterns:
    - "[a-zA-Z]*minikube"
  merge:
    annotations:
      merge-annotation: "true"
    labels:
      merge-label: "true"
      test-label: "true"
    taints:
      - effect: NoSchedule
        key: key1
        value: value1
      - effect: NoExecute
        key: key2
        value: value2
  nodeSelectorTerms:
    - matchExpressions:
        - key: beta.kubernetes.io/arch
          operator: In
          values:
            - amd64
            - arch
```

## Overwrite Attributes

This operator helps also in overwriting existing labels, annotations, and taints.
For labels and annotations, even though the label or annotation does not exist, it would be added.

But for taints it is not the case. In fact, it would override only the existing keys.

As I've said I am open to any useful suggestions.

Example:

```yaml
apiVersion: kubebuilder.kube.node.labeler.io/v1alpha1
kind: NodeLabeler
metadata:
  annotations:
    annotation1: value1
  generation: 1
  labels:
    type: unit-test
  name: simple-node-labeler
  namespace: default
spec:
  nodeSelectorTerms:
    - matchExpressions:
        - key: beta.kubernetes.io/arch
          operator: In
          values:
            - amd64
            - arch
  overwrite:
    annotations:
      merge-annotation: "false"
      overwrite-annotation: "true"
    labels:
      merge-label: "false"
      overwrite-label: "true"
    taints:
      - effect: PreferNoSchedule
        key: key1
        value: value1
```

## Examples

### Minimal Manifests

```yaml
apiVersion: kubebuilder.kube.node.labeler.io/v1alpha1
kind: NodeLabeler
metadata:
  annotations:
    annotation1: value1
  generation: 1
  labels:
    type: unit-test
  name: simple-node-labeler
  namespace: default
spec:
  nodeSelectorTerms:
    - matchExpressions:
        - key: beta.kubernetes.io/arch
          operator: In
          values:
            - amd64
            - arch
  overwrite:
    annotations:
      merge-annotation: "false"
      overwrite-annotation: "true"
```

```yaml
apiVersion: kubebuilder.kube.node.labeler.io/v1alpha1
kind: NodeLabeler
metadata:
  annotations:
    annotation1: value1
  generation: 1
  labels:
    type: unit-test
  name: simple-node-labeler
  namespace: default
spec:
  size: 3
  nodeNamePatterns:
    - "[a-zA-Z]*minikube"
  nodeSelectorTerms:
    - matchExpressions:
        - key: beta.kubernetes.io/arch
          operator: In
          values:
            - amd64
            - arch
  merge:
    annotations:
      merge-annotation: "true"
      overwrite-annotation: "false"
  overwrite:
    annotations:
      merge-annotation: "false"
      overwrite-annotation: "true"
    labels:
      merge-label: "false"
      overwrite-label: "true"
    taints:
      - effect: PreferNoSchedule
        key: key1
        value: value1
```

### More Complex Manifests

```yaml
apiVersion: kubebuilder.kube.node.labeler.io/v1alpha1
kind: NodeLabeler
metadata:
  annotations:
    annotation1: value1
  generation: 1
  labels:
    type: unit-test
  name: simple-node-labeler
  namespace: default
spec:
  size: 3
  nodeNamePatterns:
    - "[a-zA-Z]*minikube"
  merge:
    annotations:
      merge-annotation: "true"
    labels:
      merge-label: "true"
      test-label: "true"
    taints:
      - effect: NoSchedule
        key: key1
        value: value1
      - effect: NoExecute
        key: key2
        value: value2
  nodeSelectorTerms:
    - matchExpressions:
        - key: beta.kubernetes.io/arch
          operator: In
          values:
            - amd64
            - arch
```

## Prerequisets

- Kubernetes cluster
- Kubebuilder
- Git

## Getting Started

You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) or [Minikube](https://minikube.sigs.k8s.io/docs/) to get a local cluster for testing, or run against a remote cluster.

**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running the Operator on the cluster

1- Clone the repository

```sh
git clone https://github.com/AhmedGrati/kube-node-labeler
```

2- Use Kubebuilder to add CRD, controller or webhooks

```sh
kubebuilder create api --group kubebuilder --version v1alpha1 --kind NodeLabeler
```

```sh
kubebuilder create webhook --group batch --version v1alpha1 --kind NodeLabeler --defaulting --programmatic-validation
```

3- Compile the repository

```sh
make generate
```

4- Install the CRDs on your cluster

```sh
make install
```

You should make sure that the CRDs are installed successfully by using this command

```sh
kubectl get crd
```

5- Start the controller. Before that, make sure that ports `8081` and `8090` are available.

```sh
make run
```

6- Add a new CRD object to the cluster.

```sh
kubectl apply -f config/samples/
```

### Uninstall CRDs

To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller

UnDeploy the controller to the cluster:

```sh
make undeploy
```

## Contributing

### How it works

This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster

### Test It Out

1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions

If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## Next Steps

Project still under development. Next features would be:

- Add E2E and integration testing
- Add Release part to the pipeline
- Add Helm integration
- Add watch on new nodes to label them
- Enhance error handling
- Add Prometheus metrics

## LICENSE

The source code for the site is licensed under the Apache2 License, which you can find in the LICENSE file.
