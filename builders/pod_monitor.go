package builders

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	monitoringv1 "opendev.org/vexxhost/openstack-operator/api/monitoring/v1"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// PodMonitorBuilder provides an interface to build podmonitors
type PodMonitorBuilder struct {
	obj                 *monitoringv1.PodMonitor
	podMetricsEndpoints []*PodMetricsEndpointBuilder
	owner               metav1.Object
	scheme              *runtime.Scheme
}

// PodMonitor returns a new podmonitor builder
func PodMonitor(existing *monitoringv1.PodMonitor, owner metav1.Object, scheme *runtime.Scheme) *PodMonitorBuilder {
	return &PodMonitorBuilder{
		obj:    existing,
		owner:  owner,
		scheme: scheme,
	}
}

func (pm *PodMonitorBuilder) Selector(matchLabels map[string]string) *PodMonitorBuilder {
	pm.obj.Spec.Selector = metav1.LabelSelector{
		MatchLabels: matchLabels,
	}
	return pm
}

func (pm *PodMonitorBuilder) PodTargetLabels(podTargetLabels []string) *PodMonitorBuilder {
	pm.obj.Spec.PodTargetLabels = podTargetLabels
	return pm
}
func (pm *PodMonitorBuilder) JobLabel(jobLabel string) *PodMonitorBuilder {
	pm.obj.Spec.JobLabel = jobLabel
	return pm
}
func (pm *PodMonitorBuilder) NamespaceSelector(any bool, matchNames []string) *PodMonitorBuilder {
	pm.obj.Spec.NamespaceSelector = monitoringv1.NamespaceSelector{
		Any:        any,
		MatchNames: matchNames,
	}

	return pm
}
func (pm *PodMonitorBuilder) SampleLimit(sampleLimit uint64) *PodMonitorBuilder {
	pm.obj.Spec.SampleLimit = sampleLimit
	return pm
}

func (pm *PodMonitorBuilder) PodMetricsEndpoints(pme ...*PodMetricsEndpointBuilder) *PodMonitorBuilder {
	pm.podMetricsEndpoints = pme
	return pm
}

// Build returns the object after making certain assertions
func (pm *PodMonitorBuilder) Build() error {
	pm.obj.Spec.PodMetricsEndpoints = []monitoringv1.PodMetricsEndpoint{}
	for _, pmeBuilder := range pm.podMetricsEndpoints {
		podMetricsEndpoint, err := pmeBuilder.Build()
		if err != nil {
			return err
		}

		pm.obj.Spec.PodMetricsEndpoints = append(pm.obj.Spec.PodMetricsEndpoints, podMetricsEndpoint)
	}

	if !pm.isOwnedByOthers() {
		return controllerutil.SetControllerReference(pm.owner, pm.obj, pm.scheme)
	}
	return nil
}

// isOwnedByOthers checks if this podMonitor has been possessed by an another object already.
func (pm *PodMonitorBuilder) isOwnedByOthers() bool {
	ownerName := pm.owner.GetName()

	existingRefs := pm.obj.GetOwnerReferences()
	for _, r := range existingRefs {
		if r.Name == ownerName {
			return false
		} else if r.Controller != nil && *r.Controller {
			return true
		}
	}
	return false
}
