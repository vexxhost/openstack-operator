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
	"sort"
	"strconv"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infrastructurev1alpha1 "opendev.org/vexxhost/openstack-operator/api/v1alpha1"

	monitoringv1 "opendev.org/vexxhost/openstack-operator/api/monitoring/v1"
	"opendev.org/vexxhost/openstack-operator/builders"
	"opendev.org/vexxhost/openstack-operator/utils"
)

// MemcachedReconciler reconciles a Memcached object
type MemcachedReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.vexxhost.cloud,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.vexxhost.cloud,resources=memcacheds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheusrules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=podmonitors,verbs=get;list;watch;create;update;patch;delete
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
		"app.kubernetes.io/name":       "memcached",
		"app.kubernetes.io/instance":   req.Name,
		"app.kubernetes.io/managed-by": "openstack-operator",
	}

	// Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("memcached-%s", req.Name),
			Labels:    labels,
		},
	}

	op, err := utils.CreateOrUpdate(ctx, r, deployment, func() error {
		return builders.Deployment(deployment, &memcached, r.Scheme).
			Labels(labels).
			Replicas(2).
			PodTemplateSpec(
				builders.PodTemplateSpec().
					Labels(labels).
					PodSpec(
						builders.PodSpec().
							NodeSelector(memcached.Spec.NodeSelector).
							Tolerations(memcached.Spec.Tolerations).
							Containers(
								builders.Container("memcached", "vexxhost/memcached:latest").
									Args("-m", strconv.Itoa(size)).
									Port("memcached", 11211).PortProbe("memcached", 10, 30).
									Resources(1000, int64(size), 500, 1.10).
									SecurityContext(
										builders.SecurityContext().
											RunAsUser(1001),
									),
								builders.Container("exporter", "vexxhost/memcached_exporter:latest").
									Port("metrics", 9150).HTTPProbe("metrics", "/metrics", 10, 30).
									Resources(500, 128, 500, 2).
									SecurityContext(
										builders.SecurityContext().
											RunAsUser(1001),
									),
							),
					),
			).
			Build()
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	log.WithValues("resource", "Deployment").WithValues("op", op).Info("Reconciled")

	// PodMonitor

	podMonitor := &monitoringv1.PodMonitor{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "monitoring.coreos.com/v1",
			Kind:       "PodMonitor",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("memcached-podmonitor"),
			Labels: map[string]string{
				"app.kubernetes.io/name":       "memcached",
				"app.kubernetes.io/managed-by": "openstack-operator",
			},
		},
	}

	op, err = utils.CreateOrUpdate(ctx, r, podMonitor, func() error {
		return builders.PodMonitor(podMonitor, &memcached, r.Scheme).
			Selector(map[string]string{
				"app.kubernetes.io/name": "memcached",
			}).
			PodMetricsEndpoints(
				builders.PodMetricsEndpoint().
					Port("metrics").
					Path("/metrics").
					Interval("15s"),
			).Build()

	})
	if err != nil {
		return ctrl.Result{}, err
	}
	log.WithValues("resource", "podmonitor").WithValues("op", op).Info("Reconciled")

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

	// Alertrule
	alertRule := &monitoringv1.PrometheusRule{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("memcached-alertrule"),
		},
	}
	op, err = utils.CreateOrUpdate(ctx, r, alertRule, func() error {

		return builders.PrometheusRule(alertRule, &memcached, r.Scheme).
			RuleGroups(builders.RuleGroup().
				Name("memcached-rule").
				Rules(

					builders.Rule().
						Alert("MemcachedConnectionLimit").
						Message("This memcached connection is over max.").
						Priority(1).
						Expr("memcached_current_connections/memcached_max_connections*100 >90"),
				).
				Interval("1m")).
			Build()
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	log.WithValues("resource", "memcached-alertrule").WithValues("op", op).Info("Reconciled")

	// Make sure that they're sorted so we're idempotent
	sort.Strings(servers)

	// Mcrouter
	mcrouter := &infrastructurev1alpha1.Mcrouter{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("memcached-%s", req.Name),
			Labels:    labels,
		},
	}
	op, err = controllerutil.CreateOrUpdate(ctx, r, mcrouter, func() error {
		mcrouter.Spec.NodeSelector = memcached.Spec.NodeSelector
		mcrouter.Spec.Tolerations = memcached.Spec.Tolerations

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
		Owns(&infrastructurev1alpha1.Mcrouter{}).
		Owns(&monitoringv1.PodMonitor{}).
		Complete(r)
}
