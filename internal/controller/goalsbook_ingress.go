package controller

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webappv1 "goals.dev/goalsbook/api/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *GoalsbookReconciler) GoalsbookIngressReconcile(ctx context.Context, req ctrl.Request) error {
	log := log.FromContext(ctx)

	goalsbook := &webappv1.Goalsbook{}

	err := r.Get(ctx, types.NamespacedName{
		Name:      req.Name,
		Namespace: req.Namespace,
	}, goalsbook)

	if err != nil {
		return err
	}
	ingressName := goalsbook.Name + "-ingress"

	ingress := &networkingv1.Ingress{}
	err = r.Get(ctx, types.NamespacedName{
		Name:      ingressName,
		Namespace: goalsbook.Namespace,
	}, ingress)

	if err != nil && errors.IsNotFound(err) {
		log.Info("Ingress not found, creating new one", "IngressName", ingressName)

		ingress = &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ingressName,
				Namespace: req.Namespace,
				Labels: map[string]string{
					"app": goalsbook.Name,
				},
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
												Name: goalsbook.Name + "-svc",
												Port: networkingv1.ServiceBackendPort{Number: goalsbook.Spec.ContainerPort},
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
	} else if err != nil {
		log.Error(err, "Failed to fetch Ingress")
		return err
	}
	log.Info("Ingress already exists, no changes required", "IngressName", ingressName)

	// log.Info("Ingress has found, Need to be updated", "ingressName", ingressName)

	// if err := r.Status().Update(ctx, ingress); err != nil {
	// 	log.Error(err, "Error encountered while updating secret: ", ingressName)
	// 	return err
	// }

	return nil
}
