# ACE Admission Webhook

## Functionality

This project implements a custom [validating admission controller](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/) that validates the creation and update of ACE IntegrationServer instances. The admission controller is written in [Golang](https://go.dev/doc/) and is packaged as a [Helm](https://helm.sh/) chart.

It performs validation of anACE integration server.

## Building

Compile the go code as follows:

export CGO_ENABLED=0

go build -o bin/webhook main.go

Build the image. Example uses podman.

./buildImage.sh
 

## Installing

The admission WebHook is packaged as a Helm chart. The pipeline [deploy.yaml](deploy.yaml) uses Helm to deploy the Helm chart. The admission controller should be installed in every ACE OCP cluster.

The following chart [values](helm/ace-server-validating-webhook/values.yaml) may be set:

* isProduction

    A boolean that determines whether or not the OCP cluster is a production environment. Default: false.

* containerRegistry

    The host name of the container registry ACE container images must be pulled from. Default: dtr.metlife.com

Install command:

 helm install ace-server-validating-webhook ace-server-validating-webhook

## References

* This project is based up [Kubernetes admission control with validating webhooks](https://developers.redhat.com/articles/2021/09/17/kubernetes-admission-control-validating-webhooks)
* [Dynamic Admission Control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)
* [Kubernetes admission controllers in 5 minutes](https://sysdig.com/blog/kubernetes-admission-controllers/)
* [Kubernetes admission controller Git](https://github.com/kubernetes-sigs/controller-runtime/tree/master/pkg/webhook/admission)
* [Kubernetes Validating Webhooks](https://medium.com/swlh/kubernetes-validating-webhook-implementation-60f3352b66a)
* [Sample Webhook](https://github.com/ChrisTheShark/sample-vwebhook)
* [Golang](https://go.dev/doc/)
* [Helm](https://helm.sh/)
