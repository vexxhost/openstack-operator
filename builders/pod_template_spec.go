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
)

// PodTemplateSpecBuilder is an interface for building a PodTemplateSpecBuilder
type PodTemplateSpecBuilder struct {
	obj     *corev1.PodTemplateSpec
	podSpec *PodSpecBuilder
}

// PodTemplateSpec returns a builder object for a PodTemplateSpec
func PodTemplateSpec() *PodTemplateSpecBuilder {
	podTemplateSpec := &corev1.PodTemplateSpec{}

	return &PodTemplateSpecBuilder{
		obj: podTemplateSpec,
	}
}

// Labels sets up the labels for a PodTemplateSpec
func (pts *PodTemplateSpecBuilder) Labels(labels map[string]string) *PodTemplateSpecBuilder {
	pts.obj.ObjectMeta = metav1.ObjectMeta{
		Labels: labels,
	}
	return pts
}

// PodSpec points this builder to PodSpec builder
func (pts *PodTemplateSpecBuilder) PodSpec(podSpec *PodSpecBuilder) *PodTemplateSpecBuilder {
	pts.podSpec = podSpec
	return pts
}

// Build generates an object ensuring that all sub-objects work
func (pts *PodTemplateSpecBuilder) Build() (corev1.PodTemplateSpec, error) {
	podSpec, err := pts.podSpec.Build()
	if err != nil {
		return corev1.PodTemplateSpec{}, err
	}

	pts.obj.Spec = podSpec
	return *pts.obj, nil
}
