# kube-node-labeler

Kube Node Labeler is a Kubernetes operator which aims to automate the management of Kubernetes nodes for administrators in terms of attributes like Labels, Taint, and Annotations.

To read more about Kubernetes operators you could read this amazing <a href="https://github.com/cncf/tag-app-delivery/blob/main/operator-wg/whitepaper/Operator-WhitePaper_v1-0.md">white paper</a> written by the CNCF(Cloud Native Computing Foundation).
## Note

Project still under development.

## Use Case
When a node is down you lose all the attached Attributes and NodeSpec, so it is frustrating for Kubernetes administrators to re-configure this node. You could generalize it for a 1xx-nodes cluster :').

## Features

- [x] Node Selection
- [x] Adding Attributes
- [x] Overwrite Attributes
- [ ] Perform Actions on Specific Number of Nodes
- [ ] Add Regex Pattern for Nodes Selection

Attributes = Labels, Taints, Annotations

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) or [Minikube](https://minikube.sigs.k8s.io/docs/) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Git and clone the repository

```sh
git clone https://github.com/AhmedGrati/kube-node-labeler
```
2- Install kubebuilder and use it to add CRD, controller or webhooks
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

## LICENSE
The source code for the site is licensed under the Apache2 License, which you can find in the LICENSE file.