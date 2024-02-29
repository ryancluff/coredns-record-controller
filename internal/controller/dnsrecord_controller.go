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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"

	pfsensev1 "github.com/ryancluff/pfsense-dns-controller/api/v1"
	pfsense_client "github.com/ryancluff/pfsense-dns-controller/internal/pfsense_client"
)

const (
	serviceField = ".spec.service"
	// ingressField      = ".spec.ingress"
	// ingressRouteField = ".spec.ingressRoute"
)

// DnsRecordReconciler reconciles a DnsRecord object
type DnsRecordReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	PfClient pfsense_client.PfsenseClient
	Name     string
	Log      logr.Logger
}

//+kubebuilder:rbac:groups=pfsense.rcluff.com,resources=dnsentries,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pfsense.rcluff.com,resources=dnsentries/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pfsense.rcluff.com,resources=dnsentries/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DnsRecord object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *DnsRecordReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	currentOverridesByName, err := r.PfClient.GetHostOverrides()
	if err != nil {
		log.Error(err, "unable to get host overrides")
		return ctrl.Result{}, err
	}

	var dnsEntries pfsensev1.DnsRecordList
	if err := r.List(ctx, &dnsEntries); err != nil {
		log.Error(err, "unable to fetch DnsRecordList")
		return ctrl.Result{}, err
	}

	var desiredHostOverridesByName = make(map[string]pfsense_client.HostOverride)
	var dnsEntriesByTag = make(map[string]pfsensev1.DnsRecord)
	for _, dnsRecord := range dnsEntries.Items {
		tag := fmt.Sprintf("%s/%s/%s", r.Name, dnsRecord.Namespace, dnsRecord.Name)
		dnsEntriesByTag[tag] = dnsRecord

		ips := dnsRecord.Spec.IPs
		if dnsRecord.Spec.Service != "" {
			var service corev1.Service
			var namespacedName = client.ObjectKey{
				Namespace: dnsRecord.Namespace,
				Name:      dnsRecord.Spec.Service,
			}
			if err := r.Get(ctx, namespacedName, &service); err != nil {
				log.Error(err, "unable to fetch Service", "service", service)
				return ctrl.Result{}, err
			}
			ips = append(ips, service.Spec.ClusterIP)
		}

		hostOverride := pfsense_client.HostOverride{
			Host:   dnsRecord.Spec.Host,
			Domain: dnsRecord.Spec.Domain,
			IPs:    ips,
			Tag:    tag,
		}

		if _, ok := desiredHostOverridesByName[hostOverride.Name()]; !ok {
			desiredHostOverridesByName[hostOverride.Name()] = hostOverride
		} else {
			log.Info("Duplicate DNS Entry", "hostOverride", hostOverride)
			dnsRecord.Status.State = pfsensev1.ErrorState
			dnsRecord.Status.IPs = []string{}
			if err := r.Status().Update(ctx, &dnsRecord); err != nil {
				log.Error(err, "unable to update DnsRecord status")
				return ctrl.Result{}, err
			}
		}
	}

	for _, hostOverride := range currentOverridesByName {
		if strings.HasPrefix(hostOverride.Tag, r.Name) {
			if dnsRecord, ok := dnsEntriesByTag[hostOverride.Tag]; ok {
				if dnsRecord.Status.State == pfsensev1.PendingState {
					dnsRecord.Status.State = pfsensev1.ReadyState
					dnsRecord.Status.IPs = hostOverride.IPs
					if err := r.Status().Update(ctx, &dnsRecord); err != nil {
						log.Error(err, "unable to update DnsRecord status")
						return ctrl.Result{}, err
					}
				}
			}
		}
	}

	newHostOverridesByName := make(map[string]pfsense_client.HostOverride)
	for name, hostOverride := range currentOverridesByName {
		if !strings.HasPrefix(hostOverride.Tag, r.Name) {
			newHostOverridesByName[name] = hostOverride
		}
	}
	for name, hostOverride := range desiredHostOverridesByName {
		dnsRecord := dnsEntriesByTag[hostOverride.Tag]
		if dnsRecord.Status.State == pfsensev1.ReadyState {
			newHostOverridesByName[name] = hostOverride
		}
	}
	for name, hostOverride := range desiredHostOverridesByName {
		dnsRecord := dnsEntriesByTag[hostOverride.Tag]
		if dnsRecord.Status.State != pfsensev1.ReadyState {
			if newHostOverride, ok := newHostOverridesByName[name]; !ok {
				newHostOverridesByName[name] = hostOverride
				dnsRecord.Status.State = pfsensev1.PendingState
				dnsRecord.Status.IPs = newHostOverride.IPs
				if err := r.Status().Update(ctx, &dnsRecord); err != nil {
					log.Error(err, "unable to update DnsRecord status")
					return ctrl.Result{}, err
				}
			} else if dnsRecord.Status.State != pfsensev1.ErrorState {
				log.Info("Duplicate DNS Entry", "hostOverride1", hostOverride, "hostOverride2", newHostOverride)
				dnsRecord.Status.State = pfsensev1.ErrorState
				dnsRecord.Status.IPs = []string{}
				if err := r.Status().Update(ctx, &dnsRecord); err != nil {
					log.Error(err, "unable to update DnsRecord status")
					return ctrl.Result{}, err
				}
			}
		}
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
		if _, ok := currentOverridesByName[name]; !ok {
			diff.Created = append(diff.Created, hostOverride)
		} else if !reflect.DeepEqual(hostOverride, currentOverridesByName[name]) {
			diff.Updated = append(diff.Updated, hostOverride)
		} else {
			unchanged = append(unchanged, hostOverride)
		}
	}
	for name, hostOverride := range currentOverridesByName {
		if _, ok := newHostOverridesByName[name]; !ok {
			diff.Deleted = append(diff.Deleted, hostOverride)
		}
	}

	if len(diff.Created) == 1 && len(diff.Updated) == 0 && len(diff.Deleted) == 0 {
		log.Info("Host Override change required", "created", diff.Created)
	}
	if len(diff.Created)+len(diff.Updated)+len(diff.Deleted) > 0 {
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
func (r *DnsRecordReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pfsensev1.DnsRecord{}).
		Complete(r)
}
