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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webappv1 "goals.dev/goalsbook/api/v1"
)

// UserbookReconciler reconciles a Userbook object
type UserbookReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.goals.dev,resources=userbooks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.goals.dev,resources=userbooks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=webapp.goals.dev,resources=userbooks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Userbook object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *UserbookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	userbook := &webappv1.Userbook{}

	err := r.Get(ctx, types.NamespacedName{
		Name:      req.Name,
		Namespace: req.Namespace,
	}, userbook)

	if err != nil {
		log.Info("CRD not found", userbook)
		return ctrl.Result{}, err
	}

	if err = r.UserbookDeploymentReconcile(ctx, req, userbook); err != nil {
		log.Info("Error Found while creating Userbook deployment", userbook.Name)
		return ctrl.Result{}, err
	}

	if err = r.AddDeploymentFinalizer(ctx, req); err != nil {
		log.Info("Error Found while adding Userbook Finalizer ", userbook.Name)
		return ctrl.Result{}, err
	}

	if err = r.UserbookServiceReconcile(ctx, req, userbook); err != nil {
		log.Info("Error Found while creating Userbook service", userbook.Name)
		return ctrl.Result{}, err
	}

	if err := r.GoalsbookIngressReconcile(ctx, req, userbook); err != nil {
		log.Error(err, "Failed to reconcile Ingress")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserbookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Userbook{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		Named("userbook").
		Complete(r)
}
