package builders

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	dnsv1 "opendev.org/vexxhost/openstack-operator/api/dns/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ZoneBuilder defines the interface to build a Zone
type ZoneBuilder struct {
	obj    *dnsv1.Zone
	owner  metav1.Object
	scheme *runtime.Scheme
}

// Zone returns a new service builder
func Zone(existing *dnsv1.Zone, owner metav1.Object, scheme *runtime.Scheme) *ZoneBuilder {
	existing.Annotations = map[string]string{}
	return &ZoneBuilder{
		obj:    existing,
		owner:  owner,
		scheme: scheme,
	}
}

// Annotation sets one set annotation
func (z *ZoneBuilder) Annotation(key, value string) *ZoneBuilder {
	z.obj.Annotations[key] = value
	return z
}

// Labels specifies labels for the Zone
func (z *ZoneBuilder) Labels(labels map[string]string) *ZoneBuilder {
	z.obj.ObjectMeta.Labels = labels
	return z
}

// Domain sets Domain for the Zone
func (z *ZoneBuilder) Domain(domain string) *ZoneBuilder {
	z.obj.Spec.Domain = domain
	return z
}

// TTL sets TTL for the Zone
func (z *ZoneBuilder) TTL(ttl int) *ZoneBuilder {
	z.obj.Spec.TTL = ttl
	return z
}

// Email sets TTL for the Email
func (z *ZoneBuilder) Email(email string) *ZoneBuilder {
	z.obj.Spec.Email = email
	return z
}

// Build returns a complete Zone object
func (z *ZoneBuilder) Build() error {
	return controllerutil.SetControllerReference(z.owner, z.obj, z.scheme)
}
