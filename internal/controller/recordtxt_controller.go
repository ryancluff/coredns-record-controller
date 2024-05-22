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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	redis "github.com/redis/go-redis/v9"
	corednsv1 "github.com/ryancluff/coredns-record-controller/api/v1"
)

// RecordTXTReconciler reconciles a RecordTXT object
type RecordTXTReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	RedisClient *redis.Client
	Config      map[string]string
}

type TXT struct {
	Text string `json:"text,omitempty"`
	TTL  int    `json:"ttl,omitempty"`
}

//+kubebuilder:rbac:groups=coredns.rcluff.com,resources=recordtxts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coredns.rcluff.com,resources=recordtxts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=coredns.rcluff.com,resources=recordtxts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RecordTXT object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *RecordTXTReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RecordTXTReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corednsv1.RecordTXT{}).
		Complete(r)
}
