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
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	pfsensev1 "github.com/ryancluff/pfsense-dns-controller/api/v1"
	pfsense_client "github.com/ryancluff/pfsense-dns-controller/internal/pfsense_client"
)

// DnsEntryReconciler reconciles a DnsEntry object
type DnsEntryReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	PfClient pfsense_client.PfsenseClient
	Name     string
	Log      logr.Logger
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

	var dnsEntries pfsensev1.DnsEntryList
	if err := r.List(ctx, &dnsEntries); err != nil {
		log.Error(err, "unable to fetch DnsEntryList")
		err = nil
	}

	var dnsEntriesByName = make(map[string]pfsense_client.HostOverride)
	for _, dnsEntry := range dnsEntries.Items {
		tag := fmt.Sprintf("%s/%s/%s", r.Name, dnsEntry.Namespace, dnsEntry.Name)
		hostOverride := pfsense_client.HostOverride{
			Host:   dnsEntry.Spec.Host,
			Domain: dnsEntry.Spec.Domain,
			IP:     []string{dnsEntry.Spec.Ip},
			Tag:    tag,
		}
		dnsEntriesByName[hostOverride.Name()] = hostOverride
	}

	hostOverridesByName, err := r.PfClient.GetHostOverrides()
	if err != nil {
		log.Error(err, "unable to get host overrides")
		return ctrl.Result{}, err
	}

	newHostOverridesByName := make(map[string]pfsense_client.HostOverride)
	for name, hostOverride := range hostOverridesByName {
		if !strings.HasPrefix(hostOverride.Tag, r.Name) {
			newHostOverridesByName[name] = hostOverride
		}
	}
	for name, hostOverride := range dnsEntriesByName {
		newHostOverridesByName[name] = hostOverride
	}

	diff := struct {
		Created []pfsense_client.HostOverride
		Updated []pfsense_client.HostOverride
		Deleted []pfsense_client.HostOverride
	}{
		Created: []pfsense_client.HostOverride{},
		Updated: []pfsense_client.HostOverride{},
		Deleted: []pfsense_client.HostOverride{},
	}

	unchanged := []pfsense_client.HostOverride{}
	for name, hostOverride := range newHostOverridesByName {
		if _, ok := hostOverridesByName[name]; !ok {
			diff.Created = append(diff.Created, hostOverride)
		} else if !reflect.DeepEqual(hostOverride, hostOverridesByName[name]) {
			diff.Updated = append(diff.Updated, hostOverride)
		} else {
			unchanged = append(unchanged, hostOverride)
		}
	}
	for name, hostOverride := range hostOverridesByName {
		if _, ok := newHostOverridesByName[name]; !ok {
			diff.Deleted = append(diff.Deleted, hostOverride)
		}
	}

	if len(diff.Created)+len(diff.Updated)+len(diff.Deleted) == 1 {
		r.Log.Info("Host Override change required", "diff", diff)
		r.Log.V(1).Info("Host Override change required", "unchanged", unchanged)
		if err = r.PfClient.SetHostOverrides(newHostOverridesByName); err != nil {
			log.Error(err, "unable to flush host overrides")
			return ctrl.Result{}, err
		}
		log.V(1).Info("Host Override change complete")
	} else {
		log.V(1).Info("Host Override change not required")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DnsEntryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pfsensev1.DnsEntry{}).
		Complete(r)
}
