package kubernetes

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/kubernetes/blueprintcr/v1"
)

// Interface extends the kubernetes.Interface to add functionality for handling the custom resources of this operator.
type Interface interface {
	kubernetes.Interface
	// EcosystemV1Alpha1 returns a getter for the custom resources of this operator.
	EcosystemV1Alpha1() V1Alpha1Interface
}

// V1Alpha1Interface provides access for Kubernetes resources.
type V1Alpha1Interface interface {
	blueprintGetter
}

// blueprintGetter provides a Kubernetes client for accessing Kubernetes Cloudogu EcoSystem Blueprint CRs.
type blueprintGetter interface {
	// Blueprints returns a client for blueprints in the given namespace.
	Blueprints(namespace string) BlueprintInterface
}

// NewClientSet creates a new instance of the client set for this operator.
func NewClientSet(config *rest.Config, clientSet *kubernetes.Clientset) (*ClientSet, error) {
	if config == nil {
		return nil, fmt.Errorf("could not instantiate ecosystem client: rest config must not be nil")
	}

	if clientSet == nil {
		return nil, fmt.Errorf("could not instantiate ecosystem client: client set must not be nil")
	}

	blueprintClient, err := newForConfig(config)
	if err != nil {
		return nil, err
	}

	return &ClientSet{
		Interface:         clientSet,
		ecosystemV1Alpha1: blueprintClient,
	}, nil
}

// ClientSet extends the kubernetes.Interface to add functionality for handling the custom resources of this operator.
type ClientSet struct {
	kubernetes.Interface
	ecosystemV1Alpha1 V1Alpha1Interface
}

// EcosystemV1Alpha1 returns a getter for the custom resources of this operator.
func (cs *ClientSet) EcosystemV1Alpha1() V1Alpha1Interface {
	return cs.ecosystemV1Alpha1
}

// newForConfig creates a new V1Alpha1Client for a given rest.Config.
func newForConfig(c *rest.Config) (*V1Alpha1Client, error) {
	config := *c
	gv := schema.GroupVersion{Group: v1.GroupVersion.Group, Version: v1.GroupVersion.Version}
	config.ContentConfig.GroupVersion = &gv
	config.APIPath = "/apis"

	s := scheme.Scheme
	err := v1.AddToScheme(s)
	if err != nil {
		return nil, err
	}

	metav1.AddToGroupVersion(s, gv)
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &V1Alpha1Client{restClient: client}, nil
}

type restInterface interface {
	rest.Interface
}

// V1Alpha1Client is a getter for the custom resources of this operator.
type V1Alpha1Client struct {
	restClient restInterface
}

// Blueprints returns a client for Blueprints in the given namespace.
func (brc *V1Alpha1Client) Blueprints(namespace string) BlueprintInterface {
	return &blueprintClient{
		client: brc.restClient,
		ns:     namespace,
	}
}
