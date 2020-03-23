package builders

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// DeploymentBuilder defines the interface to build a deployment
type DeploymentBuilder struct {
	obj             *appsv1.Deployment
	podTemplateSpec *PodTemplateSpecBuilder
	owner           metav1.Object
	scheme          *runtime.Scheme
	labels          map[string]string
}

// Deployment returns a new deployment builder
func Deployment(existing *appsv1.Deployment, owner metav1.Object, scheme *runtime.Scheme) *DeploymentBuilder {
	return &DeploymentBuilder{
		obj:    existing,
		labels: map[string]string{},
		owner:  owner,
		scheme: scheme,
	}
}

// Labels specifies labels for the deployment
func (d *DeploymentBuilder) Labels(labels map[string]string) *DeploymentBuilder {
	d.labels = labels
	d.obj.ObjectMeta.Labels = d.labels
	d.obj.Spec.Selector = &metav1.LabelSelector{MatchLabels: d.labels}
	return d
}

// Replicas defines the number of replicas
func (d *DeploymentBuilder) Replicas(replicas int32) *DeploymentBuilder {
	d.obj.Spec.Replicas = pointer.Int32Ptr(replicas)
	return d
}

// PodTemplateSpec defines a builder for the pod template spec
func (d *DeploymentBuilder) PodTemplateSpec(podTemplateSpec *PodTemplateSpecBuilder) *DeploymentBuilder {
	d.podTemplateSpec = podTemplateSpec
	return d
}

// Build creates a final deployment objet
func (d *DeploymentBuilder) Build() error {
	podTemplateSpec, err := d.podTemplateSpec.Labels(d.labels).Build()
	if err != nil {
		return err
	}

	d.obj.Spec.Template = podTemplateSpec
	return controllerutil.SetControllerReference(d.owner, d.obj, d.scheme)
}
