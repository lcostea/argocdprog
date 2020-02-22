/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "lcostea.io/knapps/pkg/apis/knapps/v1alpha1"
	scheme "lcostea.io/knapps/pkg/generated/clientset/versioned/scheme"
)

// ArgocdprogsGetter has a method to return a ArgocdprogInterface.
// A group's client should implement this interface.
type ArgocdprogsGetter interface {
	Argocdprogs(namespace string) ArgocdprogInterface
}

// ArgocdprogInterface has methods to work with Argocdprog resources.
type ArgocdprogInterface interface {
	Create(*v1alpha1.Argocdprog) (*v1alpha1.Argocdprog, error)
	Update(*v1alpha1.Argocdprog) (*v1alpha1.Argocdprog, error)
	UpdateStatus(*v1alpha1.Argocdprog) (*v1alpha1.Argocdprog, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Argocdprog, error)
	List(opts v1.ListOptions) (*v1alpha1.ArgocdprogList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Argocdprog, err error)
	ArgocdprogExpansion
}

// argocdprogs implements ArgocdprogInterface
type argocdprogs struct {
	client rest.Interface
	ns     string
}

// newArgocdprogs returns a Argocdprogs
func newArgocdprogs(c *KnappsV1alpha1Client, namespace string) *argocdprogs {
	return &argocdprogs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the argocdprog, and returns the corresponding argocdprog object, and an error if there is any.
func (c *argocdprogs) Get(name string, options v1.GetOptions) (result *v1alpha1.Argocdprog, err error) {
	result = &v1alpha1.Argocdprog{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("argocdprogs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Argocdprogs that match those selectors.
func (c *argocdprogs) List(opts v1.ListOptions) (result *v1alpha1.ArgocdprogList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ArgocdprogList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("argocdprogs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested argocdprogs.
func (c *argocdprogs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("argocdprogs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a argocdprog and creates it.  Returns the server's representation of the argocdprog, and an error, if there is any.
func (c *argocdprogs) Create(argocdprog *v1alpha1.Argocdprog) (result *v1alpha1.Argocdprog, err error) {
	result = &v1alpha1.Argocdprog{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("argocdprogs").
		Body(argocdprog).
		Do().
		Into(result)
	return
}

// Update takes the representation of a argocdprog and updates it. Returns the server's representation of the argocdprog, and an error, if there is any.
func (c *argocdprogs) Update(argocdprog *v1alpha1.Argocdprog) (result *v1alpha1.Argocdprog, err error) {
	result = &v1alpha1.Argocdprog{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("argocdprogs").
		Name(argocdprog.Name).
		Body(argocdprog).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *argocdprogs) UpdateStatus(argocdprog *v1alpha1.Argocdprog) (result *v1alpha1.Argocdprog, err error) {
	result = &v1alpha1.Argocdprog{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("argocdprogs").
		Name(argocdprog.Name).
		SubResource("status").
		Body(argocdprog).
		Do().
		Into(result)
	return
}

// Delete takes name of the argocdprog and deletes it. Returns an error if one occurs.
func (c *argocdprogs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("argocdprogs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *argocdprogs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("argocdprogs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched argocdprog.
func (c *argocdprogs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Argocdprog, err error) {
	result = &v1alpha1.Argocdprog{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("argocdprogs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}