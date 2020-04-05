package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dnsv1 "opendev.org/vexxhost/openstack-operator/api/dns/v1"
	"opendev.org/vexxhost/openstack-operator/utils/baseutils"
	"opendev.org/vexxhost/openstack-operator/utils/openstackutils"
)

// ZoneReconciler reconciles a Zone object
type ZoneReconciler struct {
	client.Client
	Log             logr.Logger
	Scheme          *runtime.Scheme
	DesignateClient *openstackutils.DesignateClientBuilder
}

// +kubebuilder:rbac:groups=dns.openstack.org,resources=zones,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dns.openstack.org,resources=zones/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile does the reconcilication for create/update/delete Zone instances
func (r *ZoneReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	var (
		Zone          dnsv1.Zone
		designateName string
		credentials   corev1.Secret
		designate     dnsv1.Designate
	)
	ctx := context.Background()
	log := r.Log.WithValues("zone", req.NamespacedName)

	// Get Zone
	if err := r.Get(ctx, req.NamespacedName, &Zone); err != nil {
		log.Error(err, "unable to fetch Zone"+req.Name+":"+req.Namespace)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Find the corresponding designate Name
	if val, ok := Zone.Annotations[_designatingAnnotation]; ok {
		designateName = val
	} else if val, ok := Zone.Annotations[_defaultDesignatingAnnotation]; ok {
		designateName = val
	} else {
		err := errors.New("no designate annotation")
		log.Error(err, "No designate annotation."+req.Name+":"+req.Namespace)
		return ctrl.Result{}, err
	}

	// Get designate instance
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      designateName,
	}, &designate); err != nil {
		log.Error(err, "unable to fetch corresponding designate "+req.Name+":"+req.Namespace)
		return ctrl.Result{}, err
	}

	// 2 Get credentials
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      designate.Spec.Credentials,
	}, &credentials); err != nil {
		log.Error(err, "unable to fetch rc secret "+designate.Spec.Credentials+":"+req.Namespace)
		return ctrl.Result{}, err
	}
	credential, ok := credentials.Data["clouds.yaml"]
	if !ok {
		err := fmt.Errorf("rc secret syntax error ")
		log.Error(err, designate.Spec.Credentials+":"+designate.Spec.CloudName)
		return ctrl.Result{}, err
	}

	// 3 Get designate client
	if err := openstackutils.DesignateClient(r.DesignateClient, credential, designate.Spec.CloudName); err != nil {
		log.WithValues("resource", "designateClient").WithValues("op", "op").Info("ClientCreationFailed" + err.Error())
		return ctrl.Result{}, err
	}

	// Use Finalizer for the async deletion
	zoneFinalizeName := "zone.finalizers.dns.openstack.org"
	if Zone.ObjectMeta.DeletionTimestamp.IsZero() {
		if !(baseutils.ContainsString(Zone.ObjectMeta.Finalizers, zoneFinalizeName)) {
			Zone.ObjectMeta.Finalizers = append(Zone.ObjectMeta.Finalizers, zoneFinalizeName)
			if err := r.Update(ctx, &Zone); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if baseutils.ContainsString(Zone.ObjectMeta.Finalizers, zoneFinalizeName) {
			if err := r.DesignateClient.DeleteZone(Zone.Spec.Domain); err != nil {
				return ctrl.Result{}, err
			}

			log.Info("Zone deletion using finalizer")
			Zone.ObjectMeta.Finalizers = baseutils.RemoveString(Zone.ObjectMeta.Finalizers, zoneFinalizeName)
			if err := r.Update(ctx, &Zone); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Create or update
	if err := r.DesignateClient.CreateOrUpdateZone(Zone.Spec.Domain, Zone.Spec.TTL, Zone.Spec.Email); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager initializes the controller with primary manager
func (r *ZoneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dnsv1.Zone{}).
		Complete(r)
}
