package admission

import (
	"context"
	"strconv"

	log "k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	//IBM does NOT provide a golang package we can use to easily unmarshal the request JSON.
	//But it's trivial to produce Go structures that can be used.
	"encoding/json"
	"os"
	"regexp"
)

// IntegrationServerValidator validates ACE IntegrationServer instances
type IntegrationServerValidator struct {
	Client client.Client
}

// Metadata struct for parsing
type MetaDataContent struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type Pod struct {
	Containers Container `json:"containers"`
}

type Container struct {
	PodRuntime PodRuntime `json:"runtime"`
}

type PodRuntime struct {
	Image string `json:"image"`
}

type License struct {
	Accept  bool   `json:"accept"`
	License string `json:"license"`
	Use     string `json:"use"`
}

type IntegrationServerSpec struct {
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`
	EnableMetrics bool              `json:"enableMetrics"`
	Version       string            `json:"version"`
	DisableRoutes bool              `json:"disableRoutes"`
	Pod           Pod               `json:"pod"`
	License       License           `json:"license"`
}

type IntegrationServer struct {
	MetaDataContent       MetaDataContent       `json:"metadata"`
	IntegrationServerSpec IntegrationServerSpec `json:"spec"`
}

func (v *IntegrationServerValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	raw := req.Object.Raw

	//server id codes may only be digits
	reServerId := regexp.MustCompile(`^\d+$`)

	//*************************************************************************
	//* Unmarshal the IntegrationServer JSON into Go structures
	//*************************************************************************
	is := IntegrationServer{}

	if err := json.Unmarshal(raw, &is); err != nil {
		msg := "Invalid request"
		log.Error(msg)
		return admission.Denied(msg)
	}

	log.Infof("Validating IntegrationServer %s in namespace %s.", is.MetaDataContent.Name, is.MetaDataContent.Namespace)

	//*************************************************************************
	// Validate the metadata content in the IntegrationServer YAML.
	//*************************************************************************
	//Validate required labels on IntegrationServer
	if len(is.MetaDataContent.Labels) == 0 {
		msg := "Required labels not defined in metadata."
		log.Error(msg)
		return admission.Denied(msg)
	} else {
		if is.MetaDataContent.Labels["ibm.com/serverid"] == "" {
			msg := "Server denied: Required label ibm.com/serverid is missing or blank in metadata."
			log.Error(msg)
			return admission.Denied(msg)
		} else if !reServerId.MatchString(is.MetaDataContent.Labels["ibm.com/serverid"]) {
			msg := "Server denied: Label ibm.com/serverid in metadata may only contain digits."
			log.Error(msg)
			return admission.Denied(msg)
		}
	}

	// Validate enableMetrics
	if !is.IntegrationServerSpec.EnableMetrics {
		msg := "Server denied: The enableMetrics property must be set to true."
		log.Error(msg)
		return admission.Denied(msg)
	}

    // License info:
	// https://www.ibm.com/docs/en/app-connect/container?topic=resources-licensing-reference-app-connect-operator
	// Validate license use. NonProduction for dev, int, qa, Production for prod.
	// Environment variable IS_PRODUCTION set to TRUE or FALSE is set on the pod.
	isProduction, err := strconv.ParseBool(os.Getenv("IS_PRODUCTION"))

	if err != nil {
		msg := "Environment variable IS_PRODUCTION value " + os.Getenv("IS_PRODUCTION") + " could not be converted to a boolean."
		log.Fatal(msg)
		return admission.Denied(msg)
	}

	reLicenseUse := regexp.MustCompile(`^\S+NonProduction$`)
	if isProduction {
		if reLicenseUse.MatchString(is.IntegrationServerSpec.License.Use) {
			msg := "License use " + is.IntegrationServerSpec.License.Use + " is not valid for a production environment."
			log.Error(msg)
			return admission.Denied(msg)
		}
	} else {
		if !reLicenseUse.MatchString(is.IntegrationServerSpec.License.Use) {
			msg := "License use " + is.IntegrationServerSpec.License.Use + " is not valid for a non-production environment."
			log.Error(msg)
			return admission.Denied(msg)
		}
	}

	msg := "IntegrationServer " + is.MetaDataContent.Name + " in namespace " + is.MetaDataContent.Namespace + " allowed."
	log.Infof(msg)
	return admission.Allowed(msg)
}
