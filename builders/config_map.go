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
