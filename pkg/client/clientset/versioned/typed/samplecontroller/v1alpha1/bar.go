/*
Copyright 2018 The Kubernetes Authors.

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
package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "k8s.io/sample-controller/pkg/apis/samplecontroller/v1alpha1"
	scheme "k8s.io/sample-controller/pkg/client/clientset/versioned/scheme"
)

// BarsGetter has a method to return a BarInterface.
// A group's client should implement this interface.
type BarsGetter interface {
	Bars(namespace string) BarInterface
}

// BarInterface has methods to work with Bar resources.
type BarInterface interface {
	Create(*v1alpha1.Bar) (*v1alpha1.Bar, error)
	Update(*v1alpha1.Bar) (*v1alpha1.Bar, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Bar, error)
	List(opts v1.ListOptions) (*v1alpha1.BarList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Bar, err error)
	BarExpansion
}

// bars implements BarInterface
type bars struct {
	client rest.Interface
	ns     string
}

// newBars returns a Bars
func newBars(c *SamplecontrollerV1alpha1Client, namespace string) *bars {
	return &bars{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the bar, and returns the corresponding bar object, and an error if there is any.
func (c *bars) Get(name string, options v1.GetOptions) (result *v1alpha1.Bar, err error) {
	result = &v1alpha1.Bar{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("bars").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Bars that match those selectors.
func (c *bars) List(opts v1.ListOptions) (result *v1alpha1.BarList, err error) {
	result = &v1alpha1.BarList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("bars").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested bars.
func (c *bars) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("bars").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a bar and creates it.  Returns the server's representation of the bar, and an error, if there is any.
func (c *bars) Create(bar *v1alpha1.Bar) (result *v1alpha1.Bar, err error) {
	result = &v1alpha1.Bar{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("bars").
		Body(bar).
		Do().
		Into(result)
	return
}

// Update takes the representation of a bar and updates it. Returns the server's representation of the bar, and an error, if there is any.
func (c *bars) Update(bar *v1alpha1.Bar) (result *v1alpha1.Bar, err error) {
	result = &v1alpha1.Bar{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("bars").
		Name(bar.Name).
		Body(bar).
		Do().
		Into(result)
	return
}

// Delete takes name of the bar and deletes it. Returns an error if one occurs.
func (c *bars) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("bars").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *bars) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("bars").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched bar.
func (c *bars) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Bar, err error) {
	result = &v1alpha1.Bar{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("bars").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}