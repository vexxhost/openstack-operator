// Copyright 2020 VEXXHOST, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	monitoringv1 "opendev.org/vexxhost/openstack-operator/api/monitoring/v1"
	infrastructurev1alpha1 "opendev.org/vexxhost/openstack-operator/api/v1alpha1"
	"opendev.org/vexxhost/openstack-operator/builders"
	"opendev.org/vexxhost/openstack-operator/utils/baseutils"
	"opendev.org/vexxhost/openstack-operator/utils/k8sutils"
)

// RabbitmqReconciler reconciles a Rabbitmq object
type RabbitmqReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

const (
	_rabbitmqDefaultUsernameCfgKey = "username"
	_rabbitmqDefaultPasswordCfgKey = "password"
	_rabbitmqBuiltinMetricPort     = 15692
	_rabbitmqPort                  = 5672
	_rabbitmqRunAsUser             = 999
	_rabbitmqRunAsGroup            = 999
)

// +kubebuilder:rbac:groups=infrastructure.vexxhost.cloud,resources=rabbitmqs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.vexxhost.cloud,resources=rabbitmqs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheusrules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=podmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets;services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile does the reconcilication of Rabbitmq instances
func (r *RabbitmqReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("rabbitmq", req.NamespacedName)

	var Rabbitmq infrastructurev1alpha1.Rabbitmq
	if err := r.Get(ctx, req.NamespacedName, &Rabbitmq); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Labels
	typeLabels := baseutils.MergeMapsWithoutOverwrite(map[string]string{
		"app.kubernetes.io/name":       "rabbitmq",
		"app.kubernetes.io/managed-by": "openstack-operator",
	}, Rabbitmq.Labels)

	labels := map[string]string{
		"app.kubernetes.io/name":       "rabbitmq",
		"app.kubernetes.io/managed-by": "openstack-operator",
		"app.kubernetes.io/instance":   req.Name,
	}

	// Deployment
	if res, err := r.ReconcileDeployment(ctx, req, &Rabbitmq, log, labels); err != nil || res != (ctrl.Result{}) {
		return res, err
	}

	// PodMonitor
	if res, err := r.ReconcilePodMonitor(ctx, req, &Rabbitmq, log, typeLabels); err != nil || res != (ctrl.Result{}) {
		return res, err
	}

	// Alertrule
	if res, err := r.ReconcilePrometheusRule(ctx, req, &Rabbitmq, log, typeLabels); err != nil || res != (ctrl.Result{}) {
		return res, err
	}

	// Service
	if res, err := r.ReconcileService(ctx, req, &Rabbitmq, log, labels); err != nil || res != (ctrl.Result{}) {
		return res, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager initializes the controller with primary manager
func (r *RabbitmqReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha1.Rabbitmq{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&monitoringv1.PodMonitor{}).
		Owns(&monitoringv1.PrometheusRule{}).
		Complete(r)
}

// ReconcileService reconciles the service
func (r *RabbitmqReconciler) ReconcileService(ctx context.Context, req ctrl.Request, rabbitmq *infrastructurev1alpha1.Rabbitmq, log logr.Logger, labels map[string]string) (ctrl.Result, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("rabbitmq-%s", req.Name),
		},
	}
	op, err := k8sutils.CreateOrUpdate(ctx, r, service, func() error {
		return builders.Service(service, rabbitmq, r.Scheme).
			Port("rabbitmq", 5672).
			Selector(labels).
			Build()
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	log.WithValues("resource", "Service").WithValues("op", op).Info("Reconciled")
	return ctrl.Result{}, nil
}

// ReconcilePodMonitor reconciles the podMonitor
func (r *RabbitmqReconciler) ReconcilePodMonitor(ctx context.Context, req ctrl.Request, rabbitmq *infrastructurev1alpha1.Rabbitmq, log logr.Logger, labels map[string]string) (ctrl.Result, error) {
	podMonitor := &monitoringv1.PodMonitor{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "monitoring.coreos.com/v1",
			Kind:       "PodMonitor",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      "rabbitmq-podmonitor",
		},
	}
	op, err := k8sutils.CreateOrUpdate(ctx, r, podMonitor, func() error {
		return builders.PodMonitor(podMonitor, rabbitmq, r.Scheme).
			Labels(labels).
			Selector(map[string]string{
				"app.kubernetes.io/name": "rabbitmq",
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
	log.WithValues("resource", "rabbitmq-podmonitor").WithValues("op", op).Info("Reconciled")
	return ctrl.Result{}, nil
}

// ReconcilePrometheusRule reconciles the prometheusRule
func (r *RabbitmqReconciler) ReconcilePrometheusRule(ctx context.Context, req ctrl.Request, rabbitmq *infrastructurev1alpha1.Rabbitmq, log logr.Logger, labels map[string]string) (ctrl.Result, error) {
	alertRule := &monitoringv1.PrometheusRule{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      "rabbitmq-alertrule",
		},
	}
	op, err := k8sutils.CreateOrUpdate(ctx, r, alertRule, func() error {
		return builders.PrometheusRule(alertRule, rabbitmq, r.Scheme).
			Labels(labels).
			RuleGroups(builders.RuleGroup().
				Name("rabbitmq-rule").
				Rules(
					builders.Rule().
						Alert("RabbitmqDown").
						Message("Rabbitmq node down.").
						Priority(1).
						Expr("rabbitmq_up == 0"),
					builders.Rule().
						Alert("RabbitmqTooManyConnections").
						Message("RabbitMQ instance has too many connections.").
						Priority(1).
						Expr("rabbitmq_connectionsTotal > 1000"),
					builders.Rule().
						Alert("RabbitmqTooManyMessagesInQueue").
						Message("Queue is filling up.").
						Priority(1).
						Expr("rabbitmq_queue_messages_ready > 1000"),
					builders.Rule().
						Alert("RabbitmqSlowQueueConsuming").
						Message("Queue messages are consumed slowly.").
						Priority(1).
						Expr("time() - rabbitmq_queue_head_message_timestamp > 60"),
				).
				Interval("1m")).
			Build()
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	log.WithValues("resource", "rabbitmq-alertrule").WithValues("op", op).Info("Reconciled")
	return ctrl.Result{}, nil
}

// ReconcileDeployment reconciles the deployment
func (r *RabbitmqReconciler) ReconcileDeployment(ctx context.Context, req ctrl.Request, rabbitmq *infrastructurev1alpha1.Rabbitmq, log logr.Logger, labels map[string]string) (ctrl.Result, error) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("rabbitmq-%s", req.Name),
		},
	}
	op, err := k8sutils.CreateOrUpdate(ctx, r, deployment, func() error {
		return builders.Deployment(deployment, rabbitmq, r.Scheme).
			Labels(labels).
			Replicas(1).
			PodTemplateSpec(
				builders.PodTemplateSpec().
					PodSpec(
						builders.PodSpec().
							NodeSelector(rabbitmq.Spec.NodeSelector).
							Tolerations(rabbitmq.Spec.Tolerations).
							Containers(
								builders.Container("rabbitmq", "vexxhost/rabbitmq:latest").
									EnvVarFromSecret("RABBITMQ_DEFAULT_USER", rabbitmq.Spec.AuthSecret, _rabbitmqDefaultUsernameCfgKey).
									EnvVarFromSecret("RABBITMQ_DEFAULT_PASS", rabbitmq.Spec.AuthSecret, _rabbitmqDefaultPasswordCfgKey).
									Port("rabbitmq", _rabbitmqPort).
									Port("metrics", _rabbitmqBuiltinMetricPort).
									PortProbe("rabbitmq", 15, 30).
									Resources(500, 512, 500, 2).
									SecurityContext(
										builders.SecurityContext().
											RunAsUser(_rabbitmqRunAsUser).
											RunAsGroup(_rabbitmqRunAsGroup),
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
	return ctrl.Result{}, nil
}
