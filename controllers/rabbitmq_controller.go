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
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("rabbitmq-%s", req.Name),
		},
	}
	op, err := k8sutils.CreateOrUpdate(ctx, r, deployment, func() error {
		return builders.Deployment(deployment, &Rabbitmq, r.Scheme).
			Labels(labels).
			Replicas(1).
			PodTemplateSpec(
				builders.PodTemplateSpec().
					PodSpec(
						builders.PodSpec().
							NodeSelector(Rabbitmq.Spec.NodeSelector).
							Tolerations(Rabbitmq.Spec.Tolerations).
							Containers(
								builders.Container("rabbitmq", "vexxhost/rabbitmq:latest").
									EnvVarFromSecret("RABBITMQ_DEFAULT_USER", Rabbitmq.Spec.AuthSecret, _rabbitmqDefaultUsernameCfgKey).
									EnvVarFromSecret("RABBITMQ_DEFAULT_PASS", Rabbitmq.Spec.AuthSecret, _rabbitmqDefaultPasswordCfgKey).
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

	// PodMonitor
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

	op, err = k8sutils.CreateOrUpdate(ctx, r, podMonitor, func() error {
		return builders.PodMonitor(podMonitor, &Rabbitmq, r.Scheme).
			Labels(typeLabels).
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

	// Alertrule
	alertRule := &monitoringv1.PrometheusRule{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      "rabbitmq-alertrule",
		},
	}
	op, err = k8sutils.CreateOrUpdate(ctx, r, alertRule, func() error {
		return builders.PrometheusRule(alertRule, &Rabbitmq, r.Scheme).
			Labels(typeLabels).
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

	// Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("rabbitmq-%s", req.Name),
		},
	}
	op, err = k8sutils.CreateOrUpdate(ctx, r, service, func() error {
		return builders.Service(service, &Rabbitmq, r.Scheme).
			Port("epmd", 4369).
			Port("amqp", 5671).
			Port("distport", 25672).
			Selector(labels).
			Build()
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	log.WithValues("resource", "Service").WithValues("op", op).Info("Reconciled")

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
