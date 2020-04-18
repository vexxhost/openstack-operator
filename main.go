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

package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	dnsv1 "opendev.org/vexxhost/openstack-operator/api/dns/v1"
	"opendev.org/vexxhost/openstack-operator/controllers"
	"opendev.org/vexxhost/openstack-operator/utils/openstackutils"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = dnsv1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Create manager
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "6a915454.vexxhost.cloud",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Get Designate client
	designateClientBuilder := new(openstackutils.DesignateClientBuilder)
	designateClientBuilder.SetAuthFailed()

	// Setup controllers with manager
	setupZoneReconciler(mgr, designateClientBuilder)
	setupDesignateReconciler(mgr, designateClientBuilder)

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// setupZoneReconciler setups the Zone controller with manager
func setupZoneReconciler(mgr ctrl.Manager, designateClientBuilder *openstackutils.DesignateClientBuilder) {
	if err := (&controllers.ZoneReconciler{
		Client:          mgr.GetClient(),
		Log:             ctrl.Log.WithName("controllers").WithName("Zone"),
		Scheme:          mgr.GetScheme(),
		DesignateClient: designateClientBuilder,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Zone")
		os.Exit(1)
	}
}

// setupDesignateReconciler setups the Designate controller with manager
func setupDesignateReconciler(mgr ctrl.Manager, designateClientBuilder *openstackutils.DesignateClientBuilder) {
	if err := (&controllers.DesignateReconciler{
		Client:          mgr.GetClient(),
		Log:             ctrl.Log.WithName("controllers").WithName("Zone"),
		Scheme:          mgr.GetScheme(),
		DesignateClient: designateClientBuilder,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Designate")
		os.Exit(1)
	}
}
