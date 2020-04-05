package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dnsv1 "opendev.org/vexxhost/openstack-operator/api/dns/v1"
	"opendev.org/vexxhost/openstack-operator/builders"
	"opendev.org/vexxhost/openstack-operator/utils/baseutils"
	"opendev.org/vexxhost/openstack-operator/utils/k8sutils"
	"opendev.org/vexxhost/openstack-operator/utils/openstackutils"
)

// DesignateReconciler reconciles a Designate object
type DesignateReconciler struct {
	client.Client
	Log             logr.Logger
	Scheme          *runtime.Scheme
	DesignateClient *openstackutils.DesignateClientBuilder
}

const (
	_autoReconcilePeriod          = 15 * time.Second
	_designatingAnnotation        = "dns.openstack.org/designate"
	_defaultDesignatingAnnotation = "dns.openstack.org/is-default-designate"
)

// +kubebuilder:rbac:groups=dns.openstack.org,resources=designates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dns.openstack.org,resources=designates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dns.openstack.org,resources=zones,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dns.openstack.org,resources=zones/status,verbs=get;update;patch

// Reconcile does the reconcilication of designate instances
func (r *DesignateReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	var (
		credentials corev1.Secret
		designate   dnsv1.Designate
	)
	ctx := context.Background()
	log := r.Log.WithValues("Designate", req.NamespacedName)
	labels := map[string]string{
		"app.kubernetes.io/name":       "designate",
		"app.kubernetes.io/managed-by": "openstack-operator",
	}

	// 1 Get designate instance
	if err := r.Get(ctx, req.NamespacedName, &designate); err != nil {
		log.Error(err, "unable to fetch designate "+req.Name+":"+req.Namespace)
		return ctrl.Result{}, client.IgnoreNotFound(err)
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

	// 4 Create zone CRs
	// 4-1 Get zone list from the designate
	desinateZoneSpeclist := map[string]dnsv1.ZoneSpec{}
	desinateZoneNameList := []string{}
	designateZones, err := r.DesignateClient.ListZone()
	if err != nil {
		log.WithValues("resource", "Zone").WithValues("op", "op").Info("Error: Get zone list in the designate" + err.Error())
		return ctrl.Result{}, err
	}
	for _, zone := range designateZones {
		desinateZoneNameList = append(desinateZoneNameList, zone.Name)
		desinateZoneSpeclist[zone.Name] = dnsv1.ZoneSpec{
			Domain: zone.Name,
			Email:  zone.Email,
			TTL:    zone.TTL,
		}
	}

	log.Info("Get Zone list in the Designate")
	log.Info("Zone list in the Designate" + fmt.Sprintf("%v", desinateZoneSpeclist))

	// 4-2 Get zone list in the cluster
	clusterZoneObjectMetalist := map[string]metav1.ObjectMeta{}
	clusterZoneNameList := []string{}
	clusterZones := &dnsv1.ZoneList{}

	if err := r.List(context.Background(), clusterZones); err != nil {
		log.WithValues("resource", "Zone").WithValues("op", "op").Info("Error: Get zone list in the cluster" + err.Error())
		return ctrl.Result{}, err
	}
	for _, zone := range clusterZones.Items {
		clusterZoneNameList = append(clusterZoneNameList, zone.Spec.Domain)
		clusterZoneObjectMetalist[zone.Spec.Domain] = metav1.ObjectMeta{
			Name:      zone.Name,
			Namespace: zone.Namespace,
		}
	}
	log.Info("Zone list in the cluster" + fmt.Sprintf("%v", clusterZoneNameList))

	clusterOnlyNameList, designateOnlyNameList := baseutils.CompareStrSlice(clusterZoneNameList, desinateZoneNameList)
	log.Info("Zone list in the only cluster" + fmt.Sprintf("%v", clusterOnlyNameList))
	log.Info("Zone list in the only designate" + fmt.Sprintf("%v", designateOnlyNameList))

	// 4-3 Create zone list (designateOnlyNameList) in the cluster
	log.Info("Create Zone list in the cluster")
	for _, zoneName := range designateOnlyNameList {

		Zone := &dnsv1.Zone{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      strings.ReplaceAll(zoneName[:len(zoneName)-1], ".", "-"),
			},
		}
		op, err := k8sutils.CreateOrUpdate(ctx, r, Zone, func() error {
			return builders.Zone(Zone, &designate, r.Scheme).
				Labels(labels).
				Annotation(_designatingAnnotation, req.Name).
				Domain(zoneName).
				TTL(desinateZoneSpeclist[zoneName].TTL).
				Email(desinateZoneSpeclist[zoneName].Email).
				Build()
		})
		if err != nil {
			return ctrl.Result{}, err
		}
		log.WithValues("resource", "Zone").WithValues("op", op).Info("Reconciled")
		// err = r.Create(context.Background(), Zone)
		// if err != nil {
		// 	log.WithValues("resource", "Zone").WithValues("op", "op").Info("ZoneCreationFailed on Cluster -" + zoneName + ":" + err.Error())
		// 	return ctrl.Result{}, err
		// }
	}

	// 4-4 Delete zone list (clusterOnlyNameList) in the cluster
	log.Info("Delete Zone list in the cluster")
	for _, zoneName := range clusterOnlyNameList {
		Zone := &dnsv1.Zone{
			ObjectMeta: clusterZoneObjectMetalist[zoneName],
		}
		err = r.Delete(context.Background(), Zone)
		if err != nil {
			log.WithValues("resource", "Zone").WithValues("op", "op").Info("ZoneCreationFailed on Cluster -" + zoneName + ":" + err.Error())
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{Requeue: true, RequeueAfter: _autoReconcilePeriod}, nil
}

// SetupWithManager initializes the controller with primary manager
func (r *DesignateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dnsv1.Designate{}).
		Owns(&dnsv1.Zone{}).
		Complete(r)
}
