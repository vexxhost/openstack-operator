// Copyright 2020 VEXXHOST, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package builders

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	infrastructurev1alpha1 "opendev.org/vexxhost/openstack-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// McrouterBuilder defines the interface to build a Mcrouter
type McrouterBuilder struct {
	obj       *infrastructurev1alpha1.Mcrouter
	poolSpecs map[string]*McrouterPoolSpecBuilder
	owner     metav1.Object
	scheme    *runtime.Scheme
}

// Mcrouter returns a new mcrouter builder
func Mcrouter(existing *infrastructurev1alpha1.Mcrouter, owner metav1.Object, scheme *runtime.Scheme) *McrouterBuilder {
	if existing.Spec.Pools == nil {
		existing.Spec.Pools = map[string]infrastructurev1alpha1.McrouterPoolSpec{}
	}
	return &McrouterBuilder{
		obj:       existing,
		poolSpecs: map[string]*McrouterPoolSpecBuilder{},
		owner:     owner,
		scheme:    scheme,
	}
}

// Labels specifies labels for the Mcrouter
func (z *McrouterBuilder) Labels(labels map[string]string) *McrouterBuilder {
	z.obj.ObjectMeta.Labels = labels
	return z
}

// NodeSelector defines a NodeSelector for Mcrouter
func (z *McrouterBuilder) NodeSelector(selector map[string]string) *McrouterBuilder {
	z.obj.Spec.NodeSelector = selector
	return z
}

// Tolerations defines tolerations for Mcrouter
func (z *McrouterBuilder) Tolerations(tolerations []v1.Toleration) *McrouterBuilder {
	z.obj.Spec.Tolerations = tolerations
	return z
}

// Route defines route for Mcrouter
func (z *McrouterBuilder) Route(route string) *McrouterBuilder {
	z.obj.Spec.Route = route
	return z
}

// Pool defines one set pool for Mcrouter
func (z *McrouterBuilder) Pool(poolName string, poolSpec *McrouterPoolSpecBuilder) *McrouterBuilder {
	z.poolSpecs[poolName] = poolSpec
	return z
}

// Build returns a complete Mcrouter object
func (z *McrouterBuilder) Build() error {
	for key, value := range z.poolSpecs {
		pool, err := value.Build()
		if err != nil {
			return err
		}
		z.obj.Spec.Pools[key] = pool
	}
	return controllerutil.SetControllerReference(z.owner, z.obj, z.scheme)
}
