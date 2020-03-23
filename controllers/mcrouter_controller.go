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

	infrastructurev1alpha1 "opendev.org/vexxhost/openstack-operator/api/v1alpha1"
	"opendev.org/vexxhost/openstack-operator/builders"
	"opendev.org/vexxhost/openstack-operator/utils"
)

// McrouterReconciler reconciles a Mcrouter object
type McrouterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.vexxhost.cloud,resources=mcrouters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.vexxhost.cloud,resources=mcrouters/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=core,resources=configmaps;services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile does the reconcilication of Mcrouter instances
func (r *McrouterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("mcrouter", req.NamespacedName)

	var mcrouter infrastructurev1alpha1.Mcrouter
	if err := r.Get(ctx, req.NamespacedName, &mcrouter); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Labels
	labels := map[string]string{
		"app.kubernetes.io/name":     "mcrouter",
		"app.kubernetes.io/instance": req.Name,
	}

	// ConfigMap
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("mcrouter-%s", req.Name),
		},
	}
	op, err := utils.CreateOrUpdate(ctx, r, configMap, func() error {
		b, err := json.Marshal(mcrouter.Spec)
		if err != nil {
			return err
		}

		return builders.ConfigMap(configMap, &mcrouter, r.Scheme).
			Data("config.json", string(b)).
			Build()
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	log.WithValues("resource", "ConfigMap").WithValues("op", op).Info("Reconciled")

	// Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("mcrouter-%s", req.Name),
		},
	}
	op, err = utils.CreateOrUpdate(ctx, r, deployment, func() error {
		return builders.Deployment(deployment, &mcrouter, r.Scheme).
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
								builders.Container("exporter", "vexxhost/mcrouter_exporter:latest").
									Args("-mcrouter.address", "localhost:11211").
									Port("metrics", 9442).HTTPProbe("metrics", "/metrics", 10, 30).
									Resources(500, 128, 500, 2).
									SecurityContext(
										builders.SecurityContext().
											RunAsUser(1001),
									),
							).
							Volumes(
								builders.Volume("config").FromConfigMap(configMap.GetName()),
							),
					),
			).
			Build()
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	log.WithValues("resource", "Deployment").WithValues("op", op).Info("Reconciled")

	// Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
			Name:      fmt.Sprintf("mcrouter-%s", req.Name),
		},
	}
	op, err = utils.CreateOrUpdate(ctx, r, service, func() error {
		return builders.Service(service, &mcrouter, r.Scheme).
			Port("mcrouter", 11211).
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
func (r *McrouterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha1.Mcrouter{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
