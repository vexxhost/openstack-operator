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

// SecretBuilder defines the interface to build a Secret
type SecretBuilder struct {
	obj    *corev1.Secret
	owner  metav1.Object
	scheme *runtime.Scheme
}

// Secret returns a new secret builder
func Secret(existing *corev1.Secret, owner metav1.Object, scheme *runtime.Scheme) *SecretBuilder {
	existing.Data = map[string][]byte{}
	existing.StringData = map[string]string{}

	return &SecretBuilder{
		obj:    existing,
		owner:  owner,
		scheme: scheme,
	}
}

// Data sets a key inside this Secret
func (cm *SecretBuilder) Data(key, value string) *SecretBuilder {
	cm.obj.Data[key] = []byte(value)
	return cm
}

// StringData sets a key inside this Secret
func (cm *SecretBuilder) StringData(key, value string) *SecretBuilder {
	cm.obj.StringData[key] = value
	return cm
}

// SecretType sets the secret type
func (cm *SecretBuilder) SecretType(value string) *SecretBuilder {
	cm.obj.Type = corev1.SecretType(value)
	return cm
}

// Build returns a complete Secret object
func (cm *SecretBuilder) Build() error {
	return controllerutil.SetControllerReference(cm.owner, cm.obj, cm.scheme)
}
