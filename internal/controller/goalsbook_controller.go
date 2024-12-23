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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	webappv1 "goals.dev/goalsbook/api/v1"
)

// GoalsbookReconciler reconciles a Goalsbook object
type GoalsbookReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.goals.dev,resources=goalsbooks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.goals.dev,resources=goalsbooks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=webapp.goals.dev,resources=goalsbooks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Goalsbook object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *GoalsbookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx)

	goalsbook := &webappv1.Goalsbook{}

	// Check CRD
	err := r.Get(ctx, types.NamespacedName{
		Name:      req.Name,
		Namespace: req.Namespace,
	}, goalsbook)

	if err != nil {
		if !errors.IsNotFound(err) {
			log.Info("CRD Doesn't Exists")
			return ctrl.Result{}, err
		}
	}

	if err := r.DbSecretReconcile(ctx, req, goalsbook); err != nil {
		log.Error(err, "Failed to reconcile Secrect")
		return ctrl.Result{}, err
	}

	if err := r.DbConfigMapReconcile(ctx, req, goalsbook); err != nil {
		log.Error(err, "Failed to reconcile ConfigMap")
		return ctrl.Result{}, err
	}

	if err := r.GoalsbookPVC(ctx, req, goalsbook); err != nil {
		log.Error(err, "Failed to reconcile PVC")
		return ctrl.Result{}, err
	}

	if err := r.GoalsbookDeployment(ctx, req); err != nil {
		log.Error(err, "Failed to reconcile Deployment")
		return ctrl.Result{}, err
	}

	if err := r.AddDeploymentFinalizer(ctx, req); err != nil {
		log.Error(err, "Failed to reconcile Deployment Finalizer")
		return ctrl.Result{}, err
	}

	if err := r.GoalsbookService(ctx, req, goalsbook); err != nil {
		log.Error(err, "Failed to reconcile Service")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GoalsbookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Goalsbook{}).
		Owns(&appsv1.Deployment{}).
		Named("goalsbook").
		Complete(r)
}
