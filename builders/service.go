package builders

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ServiceBuilder defines the interface to build a service
type ServiceBuilder struct {
	obj    *corev1.Service
	owner  metav1.Object
	scheme *runtime.Scheme
}

// Service returns a new service builder
func Service(existing *corev1.Service, owner metav1.Object, scheme *runtime.Scheme) *ServiceBuilder {
	existing.Spec.Ports = []corev1.ServicePort{}

	return &ServiceBuilder{
		obj:    existing,
		owner:  owner,
		scheme: scheme,
	}
}

// Port appends a port to the service
func (s *ServiceBuilder) Port(name string, port int32) *ServiceBuilder {
	s.obj.Spec.Ports = append(s.obj.Spec.Ports, corev1.ServicePort{
		Name:       name,
		Protocol:   v1.ProtocolTCP,
		Port:       port,
		TargetPort: intstr.FromString(name),
	})
	return s
}

// Selector defines the service selectors
func (s *ServiceBuilder) Selector(labels map[string]string) *ServiceBuilder {
	s.obj.Spec.Selector = labels
	return s
}

// Build returns a complete Service object
func (s *ServiceBuilder) Build() error {
	return controllerutil.SetControllerReference(s.owner, s.obj, s.scheme)
}
