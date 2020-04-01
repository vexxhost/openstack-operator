package builders

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// StatefulSetBuilder defines the interface to build a StatefulSet
type StatefulSetBuilder struct {
	obj             *appsv1.StatefulSet
	podTemplateSpec *PodTemplateSpecBuilder
	pvcs            []*PersistentVolumeClaimBuilder
	scheme          *runtime.Scheme
	labels          map[string]string
	owner           metav1.Object
}

// StatefulSet returns a new StatefulSet builder
func StatefulSet(existing *appsv1.StatefulSet, owner metav1.Object, scheme *runtime.Scheme) *StatefulSetBuilder {
	return &StatefulSetBuilder{
		obj:    existing,
		labels: map[string]string{},
		owner:  owner,
		scheme: scheme,
	}
}

// Labels specifies labels for the StatefulSet
func (d *StatefulSetBuilder) Labels(labels map[string]string) *StatefulSetBuilder {
	d.labels = labels
	d.obj.ObjectMeta.Labels = d.labels
	return d
}

// Replicas defines the number of replicas
func (d *StatefulSetBuilder) Replicas(replicas int32) *StatefulSetBuilder {
	d.obj.Spec.Replicas = pointer.Int32Ptr(replicas)
	return d
}

// PodTemplateSpec defines a builder for the pod template spec
func (d *StatefulSetBuilder) PodTemplateSpec(podTemplateSpec *PodTemplateSpecBuilder) *StatefulSetBuilder {
	d.podTemplateSpec = podTemplateSpec
	return d
}

// PVCs defines a builder array for the PVC spec
func (d *StatefulSetBuilder) PVCs(pvcs ...*PersistentVolumeClaimBuilder) *StatefulSetBuilder {
	d.pvcs = pvcs
	return d
}

// Build creates a final StatefulSet objet
func (d *StatefulSetBuilder) Build() error {
	podTemplateSpec, err := d.podTemplateSpec.Labels(d.labels).Build()

	if err != nil {
		return err
	}

	d.obj.Spec.Template = podTemplateSpec

	for _, c := range d.pvcs {
		pvc, err := c.Build()
		if err != nil {
			return err
		}

		d.obj.Spec.VolumeClaimTemplates = append(d.obj.Spec.VolumeClaimTemplates, pvc)
	}

	return controllerutil.SetControllerReference(d.owner, d.obj, d.scheme)
}
