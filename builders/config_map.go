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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ConfigMapBuilder defines the interface to build a ConfigMap
type ConfigMapBuilder struct {
	obj    *corev1.ConfigMap
	owner  metav1.Object
	scheme *runtime.Scheme
}

// ConfigMap returns a new service builder
func ConfigMap(existing *corev1.ConfigMap, owner metav1.Object, scheme *runtime.Scheme) *ConfigMapBuilder {
	existing.Data = map[string]string{}

	return &ConfigMapBuilder{
		obj:    existing,
		owner:  owner,
		scheme: scheme,
	}
}

// Data sets a key inside this ConfigMap
func (cm *ConfigMapBuilder) Data(key, value string) *ConfigMapBuilder {
	cm.obj.Data[key] = value
	return cm
}

// Build returns a complete ConfigMap object
func (cm *ConfigMapBuilder) Build() error {
	return controllerutil.SetControllerReference(cm.owner, cm.obj, cm.scheme)
}
