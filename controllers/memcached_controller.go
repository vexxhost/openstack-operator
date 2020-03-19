/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/alecthomas/units"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infrastructurev1alpha1 "opendev.org/vexxhost/openstack-operator/api/v1alpha1"
	"opendev.org/vexxhost/openstack-operator/version"
)

// MemcachedReconciler reconciles a Memcached object
type MemcachedReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.vexxhost.cloud,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.vexxhost.cloud,resources=memcacheds/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile does the reconcilication of Memcached instances
func (r *MemcachedReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("memcached", req.NamespacedName)

	var memcached infrastructurev1alpha1.Memcached
	if err := r.Get(ctx, req.NamespacedName, &memcached); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Calculate size per shared
	size := memcached.Spec.Megabytes / 2

	// Labels
	labels := map[string]string{
		"app.kubernetes.io/name":     "memcached",
		"app.kubernetes.io/instance": req.Name,
	}

	// Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("memcached-%s", req.Name),
			Labels:    labels,
		},
	}
	op, err := controllerutil.CreateOrUpdate(ctx, r, deployment, func() error {
		if deployment.ObjectMeta.CreationTimestamp.IsZero() {
			deployment.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: labels,
			}
		}

		deployment.Spec.Replicas = pointer.Int32Ptr(2)
		deployment.Spec.Template = corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: labels,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "memcached",
						Image: fmt.Sprintf("vexxhost/memcached:%s", version.Revision),
						Args:  []string{"-m", strconv.Itoa(size)},
						Ports: []v1.ContainerPort{
							{
								Name:          "memcached",
								ContainerPort: int32(11211),
							},
						},
						Resources: v1.ResourceRequirements{
							Limits: v1.ResourceList{
								v1.ResourceCPU:              *resource.NewMilliQuantity(1000, resource.DecimalSI),
								v1.ResourceMemory:           *resource.NewQuantity(int64(size)*int64(units.MiB)+int64(size)*102*int64(units.KiB), resource.BinarySI),
								v1.ResourceEphemeralStorage: *resource.NewQuantity(int64(units.MB)*1000, resource.DecimalSI),
							},
							Requests: v1.ResourceList{
								v1.ResourceCPU:              *resource.NewMilliQuantity(100, resource.DecimalSI),
								v1.ResourceMemory:           *resource.NewQuantity(int64(size)*int64(units.MiB), resource.BinarySI),
								v1.ResourceEphemeralStorage: *resource.NewQuantity(int64(units.MB)*500, resource.DecimalSI),
							},
						},
						StartupProbe: &v1.Probe{},
						ReadinessProbe: &v1.Probe{
							Handler: v1.Handler{
								TCPSocket: &v1.TCPSocketAction{
									Port: intstr.FromString("memcached"),
								},
							},
							PeriodSeconds: int32(10),
						},
						LivenessProbe: &v1.Probe{
							Handler: v1.Handler{
								TCPSocket: &v1.TCPSocketAction{
									Port: intstr.FromString("memcached"),
								},
							},
							InitialDelaySeconds: int32(15),
							PeriodSeconds:       int32(30),
						},
					},
					{
						Name:  "exporter",
						Image: fmt.Sprintf("vexxhost/memcached_exporter:%s", version.Revision),
						Ports: []v1.ContainerPort{
							{
								Name:          "metrics",
								ContainerPort: int32(9150),
							},
						},
						Resources: v1.ResourceRequirements{
							Limits: v1.ResourceList{
								v1.ResourceCPU:              *resource.NewMilliQuantity(1000, resource.DecimalSI),
								v1.ResourceMemory:           *resource.NewQuantity(int64(units.Mebibyte)*256, resource.BinarySI),
								v1.ResourceEphemeralStorage: *resource.NewQuantity(int64(units.MB)*1000, resource.DecimalSI),
							},
							Requests: v1.ResourceList{
								v1.ResourceCPU:              *resource.NewMilliQuantity(100, resource.DecimalSI),
								v1.ResourceMemory:           *resource.NewQuantity(int64(units.Mebibyte)*128, resource.BinarySI),
								v1.ResourceEphemeralStorage: *resource.NewQuantity(int64(units.MB)*500, resource.DecimalSI),
							},
						},
						StartupProbe: &v1.Probe{},
						ReadinessProbe: &v1.Probe{
							Handler: v1.Handler{
								HTTPGet: &v1.HTTPGetAction{
									Path: string("/metrics"),
									Port: intstr.FromString("metrics"),
								},
							},
							PeriodSeconds: int32(10),
						},
						LivenessProbe: &v1.Probe{
							Handler: v1.Handler{
								HTTPGet: &v1.HTTPGetAction{
									Path: string("/metrics"),
									Port: intstr.FromString("metrics"),
								},
							},
							InitialDelaySeconds: int32(15),
							PeriodSeconds:       int32(20),
						},
					},
				},
			},
		}

		return controllerutil.SetControllerReference(&memcached, deployment, r.Scheme)
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	log.WithValues("resource", "Deployment").WithValues("op", op).Info("Reconciled")

	// Pods
	pods := &corev1.PodList{}
	err = r.List(ctx, pods, client.InNamespace(req.Namespace), client.MatchingLabels(labels))
	if err != nil {
		return ctrl.Result{}, err
	}

	// Generate list of pod IP addresses
	servers := []string{}
	for _, pod := range pods.Items {
		// NOTE(mnaser): It's not possible that there is no pod IP assiged yet
		if len(pod.Status.PodIP) == 0 {
			continue
		}

		server := fmt.Sprintf("%s:11211", pod.Status.PodIP)
		servers = append(servers, server)
	}

	// If we don't have any servers, requeue.
	if len(servers) == 0 {
		return ctrl.Result{Requeue: true}, nil
	}

	// Mcrouter
	mcrouter := &infrastructurev1alpha1.Mcrouter{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("memcached-%s", req.Name),
			Labels:    labels,
		},
	}
	op, err = controllerutil.CreateOrUpdate(ctx, r, mcrouter, func() error {
		mcrouter.Spec.Route = "PoolRoute|default"
		mcrouter.Spec.Pools = map[string]infrastructurev1alpha1.McrouterPoolSpec{
			"default": infrastructurev1alpha1.McrouterPoolSpec{
				Servers: servers,
			},
		}

		return controllerutil.SetControllerReference(&memcached, mcrouter, r.Scheme)
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	log.WithValues("resource", "Mcrouter").WithValues("op", op).Info("Reconciled")

	return ctrl.Result{}, nil
}

// SetupWithManager initializes the controller with primary manager
func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha1.Memcached{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
