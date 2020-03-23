package builders

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

// PodSpecBuilder is an interface for building a PodSpec
type PodSpecBuilder struct {
	obj        *corev1.PodSpec
	containers []*ContainerBuilder
	volumes    []*VolumeBuilder
}

// PodSpec returns a builder object for a PodSpec
func PodSpec() *PodSpecBuilder {
	podSpec := &corev1.PodSpec{
		DNSPolicy:     corev1.DNSClusterFirst,
		RestartPolicy: corev1.RestartPolicyAlways,
		SchedulerName: "default-scheduler",
		// SecurityContext: &v1.PodSecurityContext{
		// 	RunAsNonRoot: pointer.BoolPtr(true),
		// },
		TerminationGracePeriodSeconds: pointer.Int64Ptr(10),
	}

	return &PodSpecBuilder{
		obj: podSpec,
	}
}

// Containers appends a container builder to the PodSpec
func (ps *PodSpecBuilder) Containers(c ...*ContainerBuilder) *PodSpecBuilder {
	ps.containers = c
	return ps
}

// Volumes appends a volume builder to the PodSpec
func (ps *PodSpecBuilder) Volumes(v ...*VolumeBuilder) *PodSpecBuilder {
	ps.volumes = v
	return ps
}

// NodeSelector defines a NodeSelector for PodSpec
func (ps *PodSpecBuilder) NodeSelector(selector map[string]string) *PodSpecBuilder {
	ps.obj.NodeSelector = selector
	return ps
}

// Tolerations defines tolerations for PodSpec
func (ps *PodSpecBuilder) Tolerations(tolerations []v1.Toleration) *PodSpecBuilder {
	ps.obj.Tolerations = tolerations
	return ps
}

// Build generates an object ensuring that all sub-objects work
func (ps *PodSpecBuilder) Build() (corev1.PodSpec, error) {
	for _, c := range ps.containers {
		container, err := c.Build()
		if err != nil {
			return corev1.PodSpec{}, err
		}

		ps.obj.Containers = append(ps.obj.Containers, container)
	}

	for _, v := range ps.volumes {
		volume := v.Build()
		ps.obj.Volumes = append(ps.obj.Volumes, volume)
	}

	return *ps.obj, nil
}
