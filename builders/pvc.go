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
	"github.com/alecthomas/units"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PersistentVolumeClaimBuilder defines the interface to build a PVC
type PersistentVolumeClaimBuilder struct {
	obj *corev1.PersistentVolumeClaim
}

// PVC returns a new PVC builder
func PersistentVolumeClaim(existing *corev1.PersistentVolumeClaim) *PersistentVolumeClaimBuilder {

	return &PersistentVolumeClaimBuilder{
		obj: existing,
	}
}

func (pvc *PersistentVolumeClaimBuilder) ReadWriteOnce() *PersistentVolumeClaimBuilder {
	pvc.obj.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	return pvc
}

// Resources defines the resource configuration for the PV
func (pvc *PersistentVolumeClaimBuilder) Resources(storage int64) *PersistentVolumeClaimBuilder {
	storage = storage * int64(units.Megabyte)
	pvc.obj.Spec.Resources = v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceStorage: *resource.NewQuantity(storage, resource.DecimalSI),
		},
	}
	return pvc
}

func (pvc *PersistentVolumeClaimBuilder) ReadOnlyMany() *PersistentVolumeClaimBuilder {
	pvc.obj.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{"ReadOnlyMany"}
	return pvc
}

func (pvc *PersistentVolumeClaimBuilder) ReadWriteMany() *PersistentVolumeClaimBuilder {
	pvc.obj.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{"ReadWriteMany"}
	return pvc
}

func (pvc *PersistentVolumeClaimBuilder) Selector(selector metav1.LabelSelector) *PersistentVolumeClaimBuilder {
	pvc.obj.Spec.Selector = &selector
	return pvc
}

func (pvc *PersistentVolumeClaimBuilder) VolumeName(name string) *PersistentVolumeClaimBuilder {
	pvc.obj.Spec.VolumeName = name
	return pvc
}

func (pvc *PersistentVolumeClaimBuilder) StorageClassName(name string) *PersistentVolumeClaimBuilder {
	pvc.obj.Spec.StorageClassName = &name
	return pvc
}

func (pvc *PersistentVolumeClaimBuilder) Block() *PersistentVolumeClaimBuilder {
	*pvc.obj.Spec.VolumeMode = corev1.PersistentVolumeBlock
	return pvc
}

func (pvc *PersistentVolumeClaimBuilder) Filesystem() *PersistentVolumeClaimBuilder {
	*pvc.obj.Spec.VolumeMode = corev1.PersistentVolumeFilesystem
	return pvc
}

// Build returns a complete PVC object
func (pvc *PersistentVolumeClaimBuilder) Build() (corev1.PersistentVolumeClaim, error) {
	return *pvc.obj, nil
}
