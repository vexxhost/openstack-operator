package builders

import (
	infrastructurev1alpha1 "opendev.org/vexxhost/openstack-operator/api/v1alpha1"
)

// McrouterPoolSpecBuilder defines the interface to build a McrouterPoolSpec
type McrouterPoolSpecBuilder struct {
	obj *infrastructurev1alpha1.McrouterPoolSpec
}

// McrouterPoolSpec returns a new mcrouterPoolSpec builder
func McrouterPoolSpec() *McrouterPoolSpecBuilder {
	poolSpec := &infrastructurev1alpha1.McrouterPoolSpec{
		Servers: []string{},
	}
	return &McrouterPoolSpecBuilder{
		obj: poolSpec,
	}
}

// Servers specifies servers for the McrouterPoolSpec
func (ps *McrouterPoolSpecBuilder) Servers(servers []string) *McrouterPoolSpecBuilder {
	ps.obj.Servers = servers
	return ps
}

// Build returns a complete McrouterPoolSpec object
func (ps *McrouterPoolSpecBuilder) Build() (infrastructurev1alpha1.McrouterPoolSpec, error) {
	return *ps.obj, nil
}
