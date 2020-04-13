package controllers

import (
	"context"
	"encoding/json"
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

// McrouterReconciler reconciles a Mcrouter object
type McrouterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var (
	_mcrouterBaseLabel = map[string]string{
		"app.kubernetes.io/name":       "mcrouter",
		"app.kubernetes.io/managed-by": "openstack-operator",
	}
)

// +kubebuilder:rbac:groups=infrastructure.vexxhost.cloud,resources=mcrouters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.vexxhost.cloud,resources=mcrouters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheusrules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=podmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps;services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile does the reconcilication of Mcrouter instances
func (r *McrouterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	var mcrouter infrastructurev1alpha1.Mcrouter
	if err := r.Get(ctx, req.NamespacedName, &mcrouter); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// ConfigMap
	if err := r.reconcileConfigmap(req, &mcrouter); err != nil {
		return ctrl.Result{}, err
	}

	// Deployment
	if err := r.reconcileDeployment(req, &mcrouter); err != nil {
		return ctrl.Result{}, err
	}

	// PodMonitor
	if err := r.reconcilePodMonitor(req, &mcrouter); err != nil {
		return ctrl.Result{}, err
	}

	// Alertrule
	if err := r.reconcilePrometheusRule(req, &mcrouter); err != nil {
		return ctrl.Result{}, err
	}

	// Service
	if err := r.reconcileService(req, &mcrouter); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager initializes the controller with primary manager
func (r *McrouterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha1.Mcrouter{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&monitoringv1.PodMonitor{}).
		Owns(&monitoringv1.PrometheusRule{}).
		Complete(r)
}

func (r *McrouterReconciler) reconcileConfigmap(req ctrl.Request, mcrouter *infrastructurev1alpha1.Mcrouter) error {
	ctx := context.Background()
	log := r.Log.WithValues("mcrouter", req.NamespacedName)

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("mcrouter-%s", req.Name),
		},
	}
	op, err := k8sutils.CreateOrUpdate(ctx, r, configMap, func() error {
		b, err := json.Marshal(mcrouter.Spec)
		if err != nil {
			return err
		}
		return builders.ConfigMap(configMap, mcrouter, r.Scheme).
			Data("config.json", string(b)).
			Build()
	})
	if err != nil {
		log.WithValues("resource", "ConfigMap").WithValues("op", op).Error(err, "Reconciled Failed")
	} else {
		log.WithValues("resource", "ConfigMap").WithValues("op", op).Info("Reconciled")
	}
	return err
}

func (r *McrouterReconciler) reconcileDeployment(req ctrl.Request, mcrouter *infrastructurev1alpha1.Mcrouter) error {
	ctx := context.Background()
	log := r.Log.WithValues("mcrouter", req.NamespacedName)

	labels := baseutils.MergeMaps(_mcrouterBaseLabel, map[string]string{
		"app.kubernetes.io/instance": req.Name,
	})
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("mcrouter-%s", req.Name),
		},
	}
	op, err := k8sutils.CreateOrUpdate(ctx, r, deployment, func() error {
		return builders.Deployment(deployment, mcrouter, r.Scheme).
			Labels(labels).
			Replicas(2).
			PodTemplateSpec(
				builders.PodTemplateSpec().
					PodSpec(
						builders.PodSpec().
							NodeSelector(mcrouter.Spec.NodeSelector).
							Tolerations(mcrouter.Spec.Tolerations).
							Containers(
								builders.Container("mcrouter", "vexxhost/mcrouter:latest").
									Args("-p", "11211", "-f", "/data/config.json").
									Port("mcrouter", 11211).PortProbe("mcrouter", 10, 30).
									Resources(500, 128, 500, 2).
									Volume("config", "/data").
									SecurityContext(
										builders.SecurityContext().
											RunAsUser(999).
											RunAsGroup(999),
									),
								builders.Container("exporter", "vexxhost/mcrouter-exporter:latest").
									Args("-mcrouter.address", "localhost:11211").
									Port("metrics", 9442).HTTPProbe("metrics", "/metrics", 10, 30).
									Resources(500, 128, 500, 2).
									SecurityContext(
										builders.SecurityContext().
											RunAsUser(1001),
									),
							).
							Volumes(
								builders.Volume("config").FromConfigMap(fmt.Sprintf("mcrouter-%s", req.Name)),
							),
					),
			).
			Build()
	})

	if err != nil {
		log.WithValues("resource", "Deployment").WithValues("op", op).Error(err, "Reconciled Failed")
	} else {
		log.WithValues("resource", "Deployment").WithValues("op", op).Info("Reconciled")
	}
	return err
}

func (r *McrouterReconciler) reconcilePodMonitor(req ctrl.Request, mcrouter *infrastructurev1alpha1.Mcrouter) error {
	ctx := context.Background()
	log := r.Log.WithValues("mcrouter", req.NamespacedName)
	labels := baseutils.MergeMapsWithoutOverwrite(_mcrouterBaseLabel, mcrouter.Labels)
	podMonitor := &monitoringv1.PodMonitor{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "monitoring.coreos.com/v1",
			Kind:       "PodMonitor",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      "mcrouter-podmonitor",
		},
	}

	op, err := k8sutils.CreateOrUpdate(ctx, r, podMonitor, func() error {
		return builders.PodMonitor(podMonitor, mcrouter, r.Scheme).
			Labels(labels).
			Selector(map[string]string{
				"app.kubernetes.io/name": "mcrouter",
			}).
			PodMetricsEndpoints(
				builders.PodMetricsEndpoint().
					Port("metrics").
					Path("/metrics").
					Interval("15s"),
			).Build()

	})
	if err != nil {
		log.WithValues("resource", "PodMonitor").WithValues("op", op).Error(err, "Reconciled Failed")
	} else {
		log.WithValues("resource", "PodMonitor").WithValues("op", op).Info("Reconciled")
	}
	return err
}

func (r *McrouterReconciler) reconcilePrometheusRule(req ctrl.Request, mcrouter *infrastructurev1alpha1.Mcrouter) error {
	ctx := context.Background()
	log := r.Log.WithValues("mcrouter", req.NamespacedName)
	labels := baseutils.MergeMapsWithoutOverwrite(_mcrouterBaseLabel, mcrouter.Labels)

	alertRule := &monitoringv1.PrometheusRule{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      "mcrouter-alertrule",
		},
	}

	op, err := k8sutils.CreateOrUpdate(ctx, r, alertRule, func() error {

		return builders.PrometheusRule(alertRule, mcrouter, r.Scheme).
			Labels(labels).
			RuleGroups(builders.RuleGroup().
				Name("mcrouter-rule").
				Rules(
					builders.Rule().
						Alert("McrouterBackendDown").
						Message("Backend Memcached servers are down.").
						Priority(1).
						Expr("mcrouter_servers{state='down'}!=0"),
					builders.Rule().
						Alert("McrouterBackendTimeout").
						Message("Backend Memcached servers are timeout.").
						Priority(1).
						Expr("mcrouter_server_memcached_timeout_count>0"),
				).
				Interval("1m")).
			Build()
	})
	if err != nil {
		log.WithValues("resource", "AlertRule").WithValues("op", op).Error(err, "Reconciled Failed")
	} else {
		log.WithValues("resource", "AlertRule").WithValues("op", op).Info("Reconciled")
	}
	return err
}

func (r *McrouterReconciler) reconcileService(req ctrl.Request, mcrouter *infrastructurev1alpha1.Mcrouter) error {
	ctx := context.Background()
	log := r.Log.WithValues("mcrouter", req.NamespacedName)
	labels := baseutils.MergeMapsWithoutOverwrite(_mcrouterBaseLabel, mcrouter.Labels)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("mcrouter-%s", req.Name),
		},
	}
	op, err := k8sutils.CreateOrUpdate(ctx, r, service, func() error {
		return builders.Service(service, mcrouter, r.Scheme).
			Port("mcrouter", 11211).
			Selector(labels).
			Build()
	})
	if err != nil {
		log.WithValues("resource", "Service").WithValues("op", op).Error(err, "Reconciled Failed")
	} else {
		log.WithValues("resource", "Service").WithValues("op", op).Info("Reconciled")
	}
	return err
}
