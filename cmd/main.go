/*
Copyright 2024.

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
	"fmt"
	"os"
	"strings"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	redis "github.com/redis/go-redis/v9"

	corednsv1 "github.com/ryancluff/coredns-record-controller/api/v1"
	"github.com/ryancluff/coredns-record-controller/internal/controller"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(corednsv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	envList := []string{
		// "REDIS_HOST",
		// "REDIS_PORT",
		// "REDIS_PASSWORD",
		"NAMESPACE"}

	env := make(map[string]string)
	missing := []string{}
	for _, key := range envList {
		var ok bool
		env[key], ok = os.LookupEnv(key)
		if !ok {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		err := fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
		setupLog.Error(err, "missing required environment variables")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsserver.Options{BindAddress: metricsAddr},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "9741f11e.rcluff.com",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "password1", // no password set
		DB:       0,           // use default DB
	})

	if err = (&controller.RecordAReconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		RedisClient: client,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "RecordA")
		os.Exit(1)
	}
	if err = (&controller.RecordAAAAReconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		RedisClient: client,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "RecordAAAA")
		os.Exit(1)
	}
	if err = (&controller.RecordCNAMEReconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		RedisClient: client,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "RecordCNAME")
		os.Exit(1)
	}
	if err = (&controller.RecordTXTReconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		RedisClient: client,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "RecordTXT")
		os.Exit(1)
	}
	if err = (&controller.RecordNSReconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		RedisClient: client,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "RecordNS")
		os.Exit(1)
	}
	if err = (&controller.RecordMXReconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		RedisClient: client,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "RecordMX")
		os.Exit(1)
	}
	if err = (&controller.RecordSRVReconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		RedisClient: client,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "RecordSRV")
		os.Exit(1)
	}
	if err = (&controller.RecordSOAReconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		RedisClient: client,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "RecordSOA")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
