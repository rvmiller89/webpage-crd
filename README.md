# webpage-crd

This is a sample [CustomResourceDefinition](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) (CRD) and [Controller](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) to extend the API running in Kubernetes clusters. Blog post about it [here](https://rvmiller.com/2020/07/04/tutorial-writing-a-kubernetes-crd-and-controller-with-kubebuilder/).

## Usage

Download and install [Kubebuilder](https://book.kubebuilder.io/quick-start.html), [Kustomize](https://kubernetes-sigs.github.io/kustomize/installation/), and [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/) to run your Kubernetes cluster locally. I prefer Kind over Minikube since it starts up faster, but you’re welcome to use any tool to deploy your cluster.

Then compile, install, and run locally:

```sh
$ make
go build -o bin/manager main.go
$ make install
kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/webpages.sandbox.rvmiller.com created
$ make run
2020-07-04T22:21:21.748-0400    INFO    setup   starting manager
...
```

Deploy an example custom resource:

`webpage.yaml`
```yaml
apiVersion: sandbox.rvmiller.com/v1beta1
kind: WebPage
metadata:
  name: sample-web-page
spec:
  html: |
    <html>
      <head>
        <title>WebPage CRD</title>
      </head>
      <body>
        <h2>This page served from a Kubernetes CRD!</h2>
      </body>
    </html>
```

```sh
$ kubectl apply -f webpage.yaml 
webpage.sandbox.rvmiller.com/sample-web-page created
$ kubectl get configmaps
NAME                     DATA   AGE
sample-web-page-config   1      10m
$ kubectl get pods
NAME                    READY   STATUS    RESTARTS   AGE
sample-web-page-nginx   1/1     Running   0          10m
$ kubectl port-forward sample-web-page-nginx 7070:80
Forwarding from 127.0.0.1:7070 -> 80
Forwarding from [::1]:7070 -> 80
```

Then you can see nginx serving the webpage:

<img src="https://rvmiller.com/wp-content/uploads/2020/07/crd.png" />
