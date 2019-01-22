/*
Copyright 2019 Samsung SDS Cloud Native Computing Team.

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

package fake

import (
	v1alpha1 "github.com/samsung-cnct/cma-operator/pkg/apis/cma/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSDSAppBundles implements SDSAppBundleInterface
type FakeSDSAppBundles struct {
	Fake *FakeCmaV1alpha1
	ns   string
}

var sdsappbundlesResource = schema.GroupVersionResource{Group: "cma.sds.samsung.com", Version: "v1alpha1", Resource: "sdsappbundles"}

var sdsappbundlesKind = schema.GroupVersionKind{Group: "cma.sds.samsung.com", Version: "v1alpha1", Kind: "SDSAppBundle"}

// Get takes name of the sDSAppBundle, and returns the corresponding sDSAppBundle object, and an error if there is any.
func (c *FakeSDSAppBundles) Get(name string, options v1.GetOptions) (result *v1alpha1.SDSAppBundle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(sdsappbundlesResource, c.ns, name), &v1alpha1.SDSAppBundle{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SDSAppBundle), err
}

// List takes label and field selectors, and returns the list of SDSAppBundles that match those selectors.
func (c *FakeSDSAppBundles) List(opts v1.ListOptions) (result *v1alpha1.SDSAppBundleList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(sdsappbundlesResource, sdsappbundlesKind, c.ns, opts), &v1alpha1.SDSAppBundleList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SDSAppBundleList{ListMeta: obj.(*v1alpha1.SDSAppBundleList).ListMeta}
	for _, item := range obj.(*v1alpha1.SDSAppBundleList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested sDSAppBundles.
func (c *FakeSDSAppBundles) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(sdsappbundlesResource, c.ns, opts))

}

// Create takes the representation of a sDSAppBundle and creates it.  Returns the server's representation of the sDSAppBundle, and an error, if there is any.
func (c *FakeSDSAppBundles) Create(sDSAppBundle *v1alpha1.SDSAppBundle) (result *v1alpha1.SDSAppBundle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(sdsappbundlesResource, c.ns, sDSAppBundle), &v1alpha1.SDSAppBundle{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SDSAppBundle), err
}

// Update takes the representation of a sDSAppBundle and updates it. Returns the server's representation of the sDSAppBundle, and an error, if there is any.
func (c *FakeSDSAppBundles) Update(sDSAppBundle *v1alpha1.SDSAppBundle) (result *v1alpha1.SDSAppBundle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(sdsappbundlesResource, c.ns, sDSAppBundle), &v1alpha1.SDSAppBundle{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SDSAppBundle), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSDSAppBundles) UpdateStatus(sDSAppBundle *v1alpha1.SDSAppBundle) (*v1alpha1.SDSAppBundle, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(sdsappbundlesResource, "status", c.ns, sDSAppBundle), &v1alpha1.SDSAppBundle{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SDSAppBundle), err
}

// Delete takes name of the sDSAppBundle and deletes it. Returns an error if one occurs.
func (c *FakeSDSAppBundles) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(sdsappbundlesResource, c.ns, name), &v1alpha1.SDSAppBundle{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSDSAppBundles) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(sdsappbundlesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.SDSAppBundleList{})
	return err
}

// Patch applies the patch and returns the patched sDSAppBundle.
func (c *FakeSDSAppBundles) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SDSAppBundle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(sdsappbundlesResource, c.ns, name, data, subresources...), &v1alpha1.SDSAppBundle{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SDSAppBundle), err
}