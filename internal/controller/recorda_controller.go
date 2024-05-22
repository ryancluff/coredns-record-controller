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
	"encoding/json"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	redis "github.com/redis/go-redis/v9"
	corednsv1 "github.com/ryancluff/coredns-record-controller/api/v1"
)

// RecordAReconciler reconciles a RecordA object
type RecordAReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	RedisClient *redis.Client
	Config      map[string]string
}

type A struct {
	IP4 string `json:"ip4,omitempty"`
	TTL int    `json:"ttl,omitempty"`
}

//+kubebuilder:rbac:groups=coredns.rcluff.com,resources=recorda,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coredns.rcluff.com,resources=recorda/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=coredns.rcluff.com,resources=recorda/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RecordA object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *RecordAReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var record corednsv1.RecordA
	if err := r.Get(ctx, req.NamespacedName, &record); err != nil {
		log.Error(err, "unable to fetch RecordA")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	a := A{
		IP4: record.Spec.IP4,
		TTL: record.Spec.TTL,
	}

	var zone string
	if record.Spec.Zone != "" {
		zone = record.Spec.Zone
	} else {
		zone = r.Config["DEFAULT_ZONE"]
	}

	var desiredZoneFile map[string][]map[string]string
	currentZoneFile, err := r.RedisClient.HGet(ctx, zone, record.Spec.Hostname).Result()
	if err == redis.Nil {
		desiredZoneFile = initializeZoneFile()
		desiredZoneFile["A"] = append(desiredZoneFile["A"], a)
	} else if err != nil {
		log.Error(err, "unable to get zone contents")
		return ctrl.Result{}, err
	}

	if existingRecordsStr, ok := zoneFile[record.Spec.Hostname]; ok {
		if redisErr := r.RedisClient.HSet(ctx, zone, record.Spec.Service, record.Spec.IP4).Err(); err != nil {
			log.Error(redisErr, "unable to set zone contents")
			return ctrl.Result{}, redisErr
		}

		var existingRecords A
		if err := json.Unmarshal([]byte(existingRecordsStr), &a); err != nil {
			log.Error(err, "unable to unmarshal existing records")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RecordAReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corednsv1.RecordA{}).
		Complete(r)
}
