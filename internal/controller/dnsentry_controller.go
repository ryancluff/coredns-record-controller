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

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	pfsensev1 "github.com/ryancluff/pfsense-dns-controller/api/v1"
	pfsense_client "github.com/ryancluff/pfsense-dns-controller/internal/pfsense_client"
)

// DnsEntryReconciler reconciles a DnsEntry object
type DnsEntryReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	pfsense_client.PfsenseClient
}

//+kubebuilder:rbac:groups=pfsense.rcluff.com,resources=dnsentries,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pfsense.rcluff.com,resources=dnsentries/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pfsense.rcluff.com,resources=dnsentries/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DnsEntry object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *DnsEntryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var dnsEntry pfsensev1.DnsEntry
	if err := r.Get(ctx, req.NamespacedName, &dnsEntry); err != nil {
		log.Error(err, "unable to fetch DnsEntry")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	response, err := r.PfsenseClient.Call("GET", "/api/v1/services/unbound/host_override", nil)
	if err != nil {
		log.Error(err, "unable to get host overrides")
		return ctrl.Result{}, err
	}
	if response["status"] != "ok" {
		err = fmt.Errorf("response status: %.0f %s", response["code"].(float64), response["status"].(string))
		log.Error(err, "unable to get host overrides")
		return ctrl.Result{}, err
	}

	var id int = -1
	for i, value := range response["data"].([]interface{}) {
		host := value.(map[string]interface{})["host"].(string)
		domain := value.(map[string]interface{})["domain"].(string)

		if host == dnsEntry.Spec.Host && domain == dnsEntry.Spec.Domain {
			id = i
		}
	}

	hostOverride := map[string]interface{}{
		"host":   dnsEntry.Spec.Host,
		"domain": dnsEntry.Spec.Domain,
		"ip":     []string{dnsEntry.Spec.Ip},
		"descr":  "managed by pfsense-dns-controller/" + dnsEntry.Name,
	}

	var method string
	if id != -1 {
		ip_current := response["data"].([]interface{})[id].(map[string]interface{})["ip"].(string)
		if ip_current == dnsEntry.Spec.Ip {
			return ctrl.Result{}, nil
		}

		hostOverride["id"] = id
		method = "PUT"
	} else {
		method = "POST"
	}

	response, err = r.PfsenseClient.Call(method, "/api/v1/services/unbound/host_override", hostOverride)

	if err != nil {
		log.Error(err, "unable to apply host overrides")
		return ctrl.Result{}, err
	}
	if response["status"] != "ok" {
		err = fmt.Errorf("response status: %.0f %s", response["code"].(float64), response["status"].(string))
		log.Error(err, "unable to apply host overrides")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DnsEntryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pfsensev1.DnsEntry{}).
		Complete(r)
}
