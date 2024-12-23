package controller

import (
	"context"
	"fmt"

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

func (r *GoalsbookReconciler) GoalsbookService(ctx context.Context, req ctrl.Request, goalsbook *webappv1.Goalsbook) error {
	log := log.FromContext(ctx)

	// Check Deployment
	deployment := &appsv1.Deployment{}

	err := r.Get(ctx, types.NamespacedName{
		Name:      req.Name,
		Namespace: req.Namespace,
	}, deployment)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Deployment Not Found! Need to create Deployment")
			return err
		}
	}

	service := &corev1.Service{}
	serviceName := goalsbook.Name + "-svc"

	// Check if the Service already exists
	err = r.Get(ctx, types.NamespacedName{
		Name:      serviceName,
		Namespace: goalsbook.Namespace,
	}, service)

	if err != nil && errors.IsNotFound(err) {
		log.Info("Service not found. Creating new Service", "Service", serviceName)
		fmt.Println("--------------------------------------------")
		fmt.Printf("Container Port: %v\n", goalsbook.Spec.ContainerPort)

		// Define a new Service
		service = &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: goalsbook.Namespace,
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: goalsbook.APIVersion,
						Kind:       goalsbook.Kind,
						Name:       goalsbook.Name,
						UID:        goalsbook.UID,
						Controller: ptr.To(true),
					},
				},
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{goalsbook.Name: goalsbook.Name + "-depl"},
				Ports: []corev1.ServicePort{
					{
						Protocol: corev1.ProtocolTCP,
						Port:     goalsbook.Spec.ContainerPort,

						TargetPort: intstr.FromInt32(goalsbook.Spec.ContainerPort),
					},
				},
				Type: corev1.ServiceTypeClusterIP,
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

	// log.Info("Service has found, Need to be updated", "serviceName", serviceName)

	// if err := r.Status().Update(ctx, service); err != nil {
	// 	log.Error(err, "Error encountered while updating secret: ", serviceName)
	// 	return err
	// }

	return nil
}
