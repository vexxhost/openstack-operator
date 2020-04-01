package builders

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

// VolumeBuilder provides an interface to build volumes
type VolumeBuilder struct {
	obj *corev1.Volume
}

// Volume returns a new volume builder
func Volume(name string) *VolumeBuilder {
	volume := &corev1.Volume{
		Name: name,
	}

	return &VolumeBuilder{
		obj: volume,
	}
}

// FromConfigMap sets the source of the volume from a ConfigMap
func (v *VolumeBuilder) FromConfigMap(name string) *VolumeBuilder {
	v.obj.VolumeSource = corev1.VolumeSource{
		ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: v1.LocalObjectReference{Name: name},
			DefaultMode:          pointer.Int32Ptr(420),
		},
	}
	return v
}

// FromSecret sets the source of the volume from a Secret
func (v *VolumeBuilder) FromSecret(name string) *VolumeBuilder {
	v.obj.VolumeSource = corev1.VolumeSource{
		Secret: &corev1.SecretVolumeSource{
			SecretName:  name,
			DefaultMode: pointer.Int32Ptr(420),
		},
	}
	return v
}

// FromPersistentVolumeClaim sets the source of the volume from a PVC
func (v *VolumeBuilder) FromPersistentVolumeClaim(name string) *VolumeBuilder {
	v.obj.VolumeSource = corev1.VolumeSource{
		PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
			ClaimName: name,
			ReadOnly:  false,
		},
	}
	return v
}

// Build returns the object after checking assertions
func (v *VolumeBuilder) Build() corev1.Volume {
	return *v.obj
}
