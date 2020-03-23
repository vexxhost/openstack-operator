package builders

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

// SecurityContextBuilder defines the interface to build a securityContext
type SecurityContextBuilder struct {
	obj *corev1.SecurityContext
}

// SecurityContext returns a new SecurityContext builder
func SecurityContext() *SecurityContextBuilder {
	securityContext := &corev1.SecurityContext{}
	return &SecurityContextBuilder{
		obj: securityContext,
	}
}

// RunAsUser sets the RunAsUser inside this SecurityContext
func (sc *SecurityContextBuilder) RunAsUser(userID int64) *SecurityContextBuilder {
	sc.obj.RunAsUser = pointer.Int64Ptr(userID)
	return sc
}

// RunAsGroup sets the RunAsGroup inside this SecurityContext
func (sc *SecurityContextBuilder) RunAsGroup(groupID int64) *SecurityContextBuilder {
	sc.obj.RunAsGroup = pointer.Int64Ptr(groupID)
	return sc
}

// RunAsNonRoot sets the RunAsNonRoot inside this SecurityContext
func (sc *SecurityContextBuilder) RunAsNonRoot(flag bool) *SecurityContextBuilder {
	sc.obj.RunAsNonRoot = pointer.BoolPtr(flag)
	return sc
}

// Build returns a complete ConfigMap object
func (sc *SecurityContextBuilder) Build() (corev1.SecurityContext, error) {
	return *sc.obj, nil
}
