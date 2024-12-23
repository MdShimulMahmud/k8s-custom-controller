package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	webappv1 "goals.dev/goalsbook/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *UserbookReconciler) UserbookServiceReconcile(ctx context.Context, req ctrl.Request, userbook *webappv1.Userbook) error {
	log := log.FromContext(ctx)

	// Check Deployment
	deployment := &appsv1.Deployment{}

	err := r.Get(ctx, types.NamespacedName{
		Name:      req.Name,
		Namespace: userbook.Namespace,
	}, deployment)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Deployment Not Found! Need to create Deployment")
			return err
		}
	}

	service := &corev1.Service{}
	serviceName := userbook.Name + "-svc"

	// Check if the Service already exists
	err = r.Get(ctx, types.NamespacedName{
		Name:      serviceName,
		Namespace: userbook.Namespace,
	}, service)

	if err != nil && errors.IsNotFound(err) {
		log.Info("Service not found. Creating new Service", "Service", serviceName)

		service = &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: userbook.Namespace,
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
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{userbook.Name: userbook.Name + "-depl"},
				Ports: []corev1.ServicePort{
					{
						Protocol: corev1.ProtocolTCP,
						Port:     userbook.Spec.ContainerPort,

						TargetPort: intstr.FromInt32(userbook.Spec.ContainerPort),
					},
				},
				Type: corev1.ServiceTypeLoadBalancer,
			},
		}

		// Create the Service
		if err = r.Create(ctx, service); err != nil {
			log.Error(err, "Failed to create Service", "Service", serviceName)
			return err
		}
		log.Info("Created new Service", "Service", serviceName)
		return nil
	}

	log.Info("Service already exists", "Service", serviceName)

	if err := r.Status().Update(ctx, service); err != nil {
		log.Error(err, "Error encountered while updating secret: ", serviceName)
		return err
	}

	return nil
}
