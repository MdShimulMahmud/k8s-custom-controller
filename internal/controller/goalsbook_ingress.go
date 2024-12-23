package controller

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webappv1 "goals.dev/goalsbook/api/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *UserbookReconciler) GoalsbookIngressReconcile(ctx context.Context, req ctrl.Request, userbook *webappv1.Userbook) error {
	log := log.FromContext(ctx)

	goalsbook := &webappv1.Goalsbook{}

	err := r.Get(ctx, types.NamespacedName{
		Name:      "backend-deployment",
		Namespace: req.Namespace,
	}, goalsbook)

	if err != nil {
		log.Info("CRD doesn't find", "Name", goalsbook.Name)
		return err
	}

	err = r.Get(ctx, types.NamespacedName{
		Name:      userbook.Name,
		Namespace: req.Namespace,
	}, userbook)

	if err != nil {
		log.Info("CRD doesn't find", "Name", userbook.Name)
		return err
	}

	service := &corev1.Service{}
	serviceName := goalsbook.Name + "-svc"

	if err := r.Get(ctx, types.NamespacedName{
		Name:      serviceName,
		Namespace: req.Namespace,
	}, service); err != nil {
		log.Info("Service not found", "ServiceName", serviceName)
		return err
	}

	serviceName = userbook.Name + "-svc"
	if err := r.Get(ctx, types.NamespacedName{
		Name:      serviceName,
		Namespace: req.Namespace,
	}, service); err != nil {
		log.Info("Service not found", "ServiceName", serviceName)
		return err
	}

	ingressName := userbook.Name + "-ingress"

	ingress := &networkingv1.Ingress{}
	err = r.Get(ctx, types.NamespacedName{
		Name:      ingressName,
		Namespace: req.Namespace,
	}, ingress)

	if err != nil && errors.IsNotFound(err) {
		log.Info("Ingress not found, creating new one", "IngressName", ingressName)

		ingress = &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ingressName,
				Namespace: req.Namespace,
				Labels: map[string]string{
					"app": userbook.Name,
				},
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
			Spec: networkingv1.IngressSpec{
				IngressClassName: ptr.To("nginx"),
				Rules: []networkingv1.IngressRule{
					{
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path:     "/",
										PathType: ptr.To(networkingv1.PathTypePrefix),
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: userbook.Name + "-svc",
												Port: networkingv1.ServiceBackendPort{Number: userbook.Spec.ContainerPort},
											},
										},
									},
									{
										Path:     "/goals",
										PathType: ptr.To(networkingv1.PathTypePrefix),
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: goalsbook.Name + "-svc",
												Port: networkingv1.ServiceBackendPort{Number: goalsbook.Spec.ContainerPort},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		// Create the Ingress resource
		if err := r.Create(ctx, ingress); err != nil {
			log.Error(err, "Failed to create Ingress", "IngressName", ingressName)
			return err
		}
		log.Info("Ingress created successfully", "IngressName", ingressName)
		return nil
	}
	log.Info("Ingress already exists, no changes required", "IngressName", ingressName)

	return nil
}
