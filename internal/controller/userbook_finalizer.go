package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	// corev1 "k8s.io/api/core/v1"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	webappv1 "goals.dev/goalsbook/api/v1"
	controllerutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const UserbookFinalizer = "goals.dev/frontend-protection"

func (r *UserbookReconciler) AddDeploymentFinalizer(ctx context.Context, req ctrl.Request) error {
	log := log.FromContext(ctx)
	userbook := &webappv1.Userbook{}

	if err := r.Get(ctx, types.NamespacedName{
		Name:      req.Name,
		Namespace: req.Namespace,
	}, userbook); err != nil {
		return err
	}

	if userbook.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(userbook, UserbookFinalizer) {
			log.Info("Adding finalizer to userbook", "userbook", userbook.Name)
			controllerutil.AddFinalizer(userbook, UserbookFinalizer)
			if err := r.Update(ctx, userbook); err != nil {
				log.Error(err, "Failed to update userbook after adding finalizer", "userbook", userbook.Name)
				return err
			}
		}
	} else {
		// deletion request
		if controllerutil.ContainsFinalizer(userbook, UserbookFinalizer) {
			controllerutil.RemoveFinalizer(userbook, UserbookFinalizer)
			log.Info("Removing finalizer from userbook", "userbook", userbook.Name)
			if err := r.Update(ctx, userbook); err != nil {
				log.Error(err, "Failed to update userbook", "userbook", userbook.Name)
				return err
			}
		}
	}

	return nil
}
