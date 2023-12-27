package kubernetes

import (
	"context"
	"time"

	v1 "github.com/cloudogu/k8s-blueprint-operator/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type BlueprintInterface interface {
	// Create takes the representation of a blueprint and creates it.  Returns the server's representation of the blueprint, and an error, if there is any.
	Create(ctx context.Context, blueprint *v1.Blueprint, opts metav1.CreateOptions) (*v1.Blueprint, error)

	// Update takes the representation of a blueprint and updates it. Returns the server's representation of the blueprint, and an error, if there is any.
	Update(ctx context.Context, blueprint *v1.Blueprint, opts metav1.UpdateOptions) (*v1.Blueprint, error)

	// UpdateStatus was generated because the type contains a Status member.
	UpdateStatus(ctx context.Context, blueprint *v1.Blueprint, opts metav1.UpdateOptions) (*v1.Blueprint, error)

	// Delete takes name of the blueprint and deletes it. Returns an error if one occurs.
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error

	// DeleteCollection deletes a collection of objects.
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error

	// Get takes name of the blueprint, and returns the corresponding blueprint object, and an error if there is any.
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Blueprint, error)

	// List takes label and field selectors, and returns the list of Blueprints that match those selectors.
	List(ctx context.Context, opts metav1.ListOptions) (*v1.BlueprintList, error)

	// Watch returns a watch.Interface that watches the requested blueprints.
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)

	// Patch applies the patch and returns the patched blueprint.
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Blueprint, err error)
}

type blueprintClient struct {
	client rest.Interface
	ns     string
}

// Get takes name of the blueprint, and returns the corresponding blueprint object, and an error if there is any.
func (d *blueprintClient) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Blueprint, err error) {
	result = &v1.Blueprint{}
	err = d.client.Get().
		Namespace(d.ns).
		Resource("blueprints").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Blueprints that match those selectors.
func (d *blueprintClient) List(ctx context.Context, opts metav1.ListOptions) (result *v1.BlueprintList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.BlueprintList{}
	err = d.client.Get().
		Namespace(d.ns).
		Resource("blueprints").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested blueprints.
func (d *blueprintClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return d.client.Get().
		Namespace(d.ns).
		Resource("blueprints").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a blueprint and creates it.  Returns the server's representation of the blueprint, and an error, if there is any.
func (d *blueprintClient) Create(ctx context.Context, blueprint *v1.Blueprint, opts metav1.CreateOptions) (result *v1.Blueprint, err error) {
	result = &v1.Blueprint{}
	err = d.client.Post().
		Namespace(d.ns).
		Resource("blueprints").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(blueprint).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a blueprint and updates it. Returns the server's representation of the blueprint, and an error, if there is any.
func (d *blueprintClient) Update(ctx context.Context, blueprint *v1.Blueprint, opts metav1.UpdateOptions) (result *v1.Blueprint, err error) {
	result = &v1.Blueprint{}
	err = d.client.Put().
		Namespace(d.ns).
		Resource("blueprints").
		Name(blueprint.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(blueprint).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (d *blueprintClient) UpdateStatus(ctx context.Context, blueprint *v1.Blueprint, opts metav1.UpdateOptions) (result *v1.Blueprint, err error) {
	result = &v1.Blueprint{}
	err = d.client.Put().
		Namespace(d.ns).
		Resource("blueprints").
		Name(blueprint.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(blueprint).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the blueprint and deletes it. Returns an error if one occurs.
func (d *blueprintClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return d.client.Delete().
		Namespace(d.ns).
		Resource("blueprints").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (d *blueprintClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return d.client.Delete().
		Namespace(d.ns).
		Resource("blueprints").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched blueprint.
func (d *blueprintClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Blueprint, err error) {
	result = &v1.Blueprint{}
	err = d.client.Patch(pt).
		Namespace(d.ns).
		Resource("blueprints").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
