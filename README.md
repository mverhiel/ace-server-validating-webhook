# ACE Admission Webhook

## Functionality

This project implements a custom [validating admission controller](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/) that validates the creation and update of ACE IntegrationServer instances. The admission controller is written in [Golang](https://go.dev/doc/) and is packaged as a [Helm](https://helm.sh/) chart.

It performs the following validations:

* It verifies that the required annotations are defined in the IntegrationServer metadata AND in the IntegrationServer spec. The annotations defined in IntegrationServer spec will be attached to all integration server pods the ACE operator creates. It also enforces the values of annotations using regular expressions. Required annotations:

    * `app.metlife.com/dpccode`

* It verifies that the required labels are defined in the IntegrationServer metadata AND in the IntegrationServer spec. The labels defined in IntegrationServer spec will be attached to all integration server pods the ACE operator creates. Required labels:

    * `app.metlife.com/eai`

        The EAI code can only consist of digits.

* It verifies that disableRoutes is set to true.
* It verifies that container images are pulled from value specifed in chart's containerRegistry value. Default value: `dtr.metlife.com`.
* It verifies that enableMetrics is set to true.
* If the OCP cluster is a non-production cluster, ensure that the license use specifies a non-production value. If the OCP cluster is a production environment, ensure that the license use specifies a production value. Chart value isProduction defines whether or not the OCP cluster is a production environment.

The admission webhook only considers AdmissionReview objects for custom projects that have the `app.metlife.com/kind: cp4i` label. See [helm/ace-validating-webhook/templates/validatingWebhookConfiguration.yaml](helm/ace-validating-webhook/templates/validatingWebhookConfiguration.yaml):

```yaml
  namespaceSelector:
    matchLabels:
        {{ include "ace-validating-webhook.webhookLabels" . }}
```

The validating webhook enforces MetLife standards for ACE integration servers and prevents the creation or update of integration servers that do not meet the standards.

## Building

The pipeline [build.yaml](build.yaml) compiles the Go executable and builds a container that is pushed to DTR. The pipelines needs to install the [Go tools](https://learn.microsoft.com/en-us/azure/devops/pipelines/tasks/reference/go-tool-v0?view=azure-pipelines) in order to build the Go source code. Since OCP is running on Linux, the pipeline must run on a Linux build agent.

## Installing

The admission WebHook is packaged as a Helm chart. The pipeline [deploy.yaml](deploy.yaml) uses Helm to deploy the Helm chart. The admission controller should be installed in every ACE OCP cluster.

The following chart [values](helm/ace-validating-webhook/values.yaml) may be set:

* isProduction

    A boolean that determines whether or not the OCP cluster is a production environment. Default: false.

* containerRegistry

    The host name of the container registry ACE container images must be pulled from. Default: dtr.metlife.com

## Testing

To test the webhook locally (after installing the webhook):

1. Open a PowerShell window.
2. Login to OCP cluster as a user that can create/delete projects.
3. Make [test](test) directory the current directory.
4. Run the test suite: `.\runtests.ps1 -isProductionParam true|false`. The isProductionParam parameter defaults to false to test non-production OCP clusters.

The test cases are automatically executed when the webhook is installed to an OCP cluster using the [deploy.yaml](deploy.yaml) pipeline.

## Developing

Install the following on your laptop to develop the code:

1. The [go toolkit](https://golang.google.cn/dl/). **You must have admin rights to run the installer.**

2. Install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go) in Visual Studio Code.

## TODOs

* Create tests cases for production clusters. The non-production tests can be copied and then change the license use value in each test case.

## References

* This project is based up [Kubernetes admission control with validating webhooks](https://developers.redhat.com/articles/2021/09/17/kubernetes-admission-control-validating-webhooks)
* [Dynamic Admission Control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)
* [Kubernetes admission controllers in 5 minutes](https://sysdig.com/blog/kubernetes-admission-controllers/)
* [Kubernetes admission controller Git](https://github.com/kubernetes-sigs/controller-runtime/tree/master/pkg/webhook/admission)
* [Kubernetes Validating Webhooks](https://medium.com/swlh/kubernetes-validating-webhook-implementation-60f3352b66a)
* [Sample Webhook](https://github.com/ChrisTheShark/sample-vwebhook)
* [Golang](https://go.dev/doc/)
* [Helm](https://helm.sh/)
