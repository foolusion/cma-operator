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

// FakeSDSPackageManagers implements SDSPackageManagerInterface
type FakeSDSPackageManagers struct {
	Fake *FakeCmaV1alpha1
	ns   string
}

var sdspackagemanagersResource = schema.GroupVersionResource{Group: "cma.sds.samsung.com", Version: "v1alpha1", Resource: "sdspackagemanagers"}

var sdspackagemanagersKind = schema.GroupVersionKind{Group: "cma.sds.samsung.com", Version: "v1alpha1", Kind: "SDSPackageManager"}

// Get takes name of the sDSPackageManager, and returns the corresponding sDSPackageManager object, and an error if there is any.
func (c *FakeSDSPackageManagers) Get(name string, options v1.GetOptions) (result *v1alpha1.SDSPackageManager, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(sdspackagemanagersResource, c.ns, name), &v1alpha1.SDSPackageManager{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SDSPackageManager), err
}

// List takes label and field selectors, and returns the list of SDSPackageManagers that match those selectors.
func (c *FakeSDSPackageManagers) List(opts v1.ListOptions) (result *v1alpha1.SDSPackageManagerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(sdspackagemanagersResource, sdspackagemanagersKind, c.ns, opts), &v1alpha1.SDSPackageManagerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SDSPackageManagerList{ListMeta: obj.(*v1alpha1.SDSPackageManagerList).ListMeta}
	for _, item := range obj.(*v1alpha1.SDSPackageManagerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested sDSPackageManagers.
func (c *FakeSDSPackageManagers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(sdspackagemanagersResource, c.ns, opts))

}

// Create takes the representation of a sDSPackageManager and creates it.  Returns the server's representation of the sDSPackageManager, and an error, if there is any.
func (c *FakeSDSPackageManagers) Create(sDSPackageManager *v1alpha1.SDSPackageManager) (result *v1alpha1.SDSPackageManager, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(sdspackagemanagersResource, c.ns, sDSPackageManager), &v1alpha1.SDSPackageManager{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SDSPackageManager), err
}

// Update takes the representation of a sDSPackageManager and updates it. Returns the server's representation of the sDSPackageManager, and an error, if there is any.
func (c *FakeSDSPackageManagers) Update(sDSPackageManager *v1alpha1.SDSPackageManager) (result *v1alpha1.SDSPackageManager, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(sdspackagemanagersResource, c.ns, sDSPackageManager), &v1alpha1.SDSPackageManager{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SDSPackageManager), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSDSPackageManagers) UpdateStatus(sDSPackageManager *v1alpha1.SDSPackageManager) (*v1alpha1.SDSPackageManager, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(sdspackagemanagersResource, "status", c.ns, sDSPackageManager), &v1alpha1.SDSPackageManager{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SDSPackageManager), err
}

// Delete takes name of the sDSPackageManager and deletes it. Returns an error if one occurs.
func (c *FakeSDSPackageManagers) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(sdspackagemanagersResource, c.ns, name), &v1alpha1.SDSPackageManager{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSDSPackageManagers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(sdspackagemanagersResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.SDSPackageManagerList{})
	return err
}

// Patch applies the patch and returns the patched sDSPackageManager.
func (c *FakeSDSPackageManagers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SDSPackageManager, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(sdspackagemanagersResource, c.ns, name, data, subresources...), &v1alpha1.SDSPackageManager{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SDSPackageManager), err
}
