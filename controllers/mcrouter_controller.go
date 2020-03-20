package controllers

import (
	"context"
	"encoding/json"
	"fmt"

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
			Labels:    labels,
		},
	}
	op, err := controllerutil.CreateOrUpdate(ctx, r, configMap, func() error {
		b, err := json.Marshal(mcrouter.Spec)
		if err != nil {
			return err
		}

		configMap.Data = map[string]string{
			"config.json": string(b),
		}

		return controllerutil.SetControllerReference(&mcrouter, configMap, r.Scheme)
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
			Labels:    labels,
		},
	}
	op, err = controllerutil.CreateOrUpdate(ctx, r, deployment, func() error {
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
						Name:  "mcrouter",
						Image: fmt.Sprintf("vexxhost/mcrouter:%s", version.Revision),
						Args:  []string{"-p", "11211", "-f", "/data/config.json"},
						Ports: []v1.ContainerPort{
							{
								Name:          "mcrouter",
								ContainerPort: int32(11211),
							},
						},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "config",
								MountPath: "/data",
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
								TCPSocket: &v1.TCPSocketAction{
									Port: intstr.FromString("mcrouter"),
								},
							},
							PeriodSeconds: int32(10),
						},
						LivenessProbe: &v1.Probe{
							Handler: v1.Handler{
								TCPSocket: &v1.TCPSocketAction{
									Port: intstr.FromString("mcrouter"),
								},
							},
							InitialDelaySeconds: int32(15),
							PeriodSeconds:       int32(30),
						},
					},
					{
						Name:  "exporter",
						Image: fmt.Sprintf("vexxhost/mcrouter_exporter:%s", version.Revision),
						Args:  []string{"-mcrouter.address", "localhost:11211"},
						Ports: []v1.ContainerPort{
							{
								Name:          "metrics",
								ContainerPort: int32(9442),
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
							InitialDelaySeconds: int32(5),
							PeriodSeconds:       int32(10),
						},
						LivenessProbe: &v1.Probe{
							Handler: v1.Handler{
								HTTPGet: &v1.HTTPGetAction{
									Path: string("/metrics"),
									Port: intstr.FromString("metrics"),
								},
							},
							InitialDelaySeconds: int32(15),
							PeriodSeconds:       int32(30),
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{Name: configMap.GetName()},
							},
						},
					},
				},
				NodeSelector: mcrouter.Spec.NodeSelector,
				Tolerations:  mcrouter.Spec.Tolerations,
			},
		}

		return controllerutil.SetControllerReference(&mcrouter, deployment, r.Scheme)
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
			Labels:    labels,
		},
	}
	op, err = controllerutil.CreateOrUpdate(ctx, r, service, func() error {
		service.Spec.Type = corev1.ServiceTypeClusterIP
		service.Spec.Ports = []v1.ServicePort{
			{
				Name:       "mcrouter",
				Port:       int32(11211),
				TargetPort: intstr.FromString("mcrouter"),
			},
		}
		service.Spec.Selector = labels

		return controllerutil.SetControllerReference(&mcrouter, service, r.Scheme)
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
