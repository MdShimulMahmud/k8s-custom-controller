package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	// corev1 "k8s.io/api/core/v1"

	"k8s.io/utils/ptr"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	webappv1 "goals.dev/goalsbook/api/v1"
)

func (r *UserbookReconciler) UserbookDeploymentReconcile(ctx context.Context, req ctrl.Request, userbook *webappv1.Userbook) error {

	log := log.FromContext(ctx)

	deployment := &appsv1.Deployment{}

	err := r.Get(ctx, types.NamespacedName{
		Name:      userbook.Name,
		Namespace: req.Namespace,
	}, deployment)

	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating new Deployment", "Deployment", userbook.Name)

		deployment = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      userbook.Name,
				Namespace: req.Namespace,
				OwnerReferences: []metav1.OwnerReference{

					{
						APIVersion: userbook.APIVersion,
						Kind:       userbook.Kind,
						Name:       userbook.Name,
						UID:        userbook.UID,
						Controller: ptr.To(true),
					},
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: userbook.Spec.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{userbook.Name: userbook.Name + "-depl"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{userbook.Name: userbook.Name + "-depl"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  userbook.Name + "-pod",
								Image: userbook.Spec.ImageName,
								Ports: []corev1.ContainerPort{{ContainerPort: userbook.Spec.ContainerPort}},
							},
						},
					},
				},
			},
		}
		if err = r.Create(ctx, deployment); err != nil {
			log.Error(err, "Failed to create Deployment", "Deployment:::", err.Error())
			return err
		}
		log.Info("Created new Deployment", "Deployment", userbook.Name)
		return nil
	} else {
		log.Info("Deployment already exists")

		if err = r.Status().Update(ctx, deployment); err != nil {
			log.Error(err, "Failed to update deployments Status", "Deployment:::", userbook.Name, "Err:::", err.Error())
			return err
		}
		log.Info("Updated Deployment", "Deployment", userbook.Name)

		return nil
	}

}
