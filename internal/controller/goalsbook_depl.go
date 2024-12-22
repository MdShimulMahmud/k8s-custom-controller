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

func (r *GoalsbookReconciler) GoalsbookDeployment(ctx context.Context, req ctrl.Request) error {

	log := log.FromContext(ctx)
	goalsbook := &webappv1.Goalsbook{}

	err := r.Get(ctx, types.NamespacedName{
		Name:      req.Name,
		Namespace: req.Namespace,
	}, goalsbook)

	if err != nil {
		return err
	}
	deployment := &appsv1.Deployment{}

	err = r.Get(ctx, types.NamespacedName{
		Name:      req.Name,
		Namespace: req.Namespace,
	}, deployment)

	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating new Deployment", "Deployment", goalsbook.Name)

		deployment = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      goalsbook.Name,
				Namespace: req.Namespace,
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
			Spec: appsv1.DeploymentSpec{
				Replicas: goalsbook.Spec.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{goalsbook.Name: goalsbook.Name + "-depl"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{goalsbook.Name: goalsbook.Name + "-depl"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  goalsbook.Name + "-pod",
								Image: goalsbook.Spec.ImageName,
								Ports: []corev1.ContainerPort{{ContainerPort: goalsbook.Spec.ContainerPort}},
								Env: []corev1.EnvVar{
									{
										Name: "DATABASE_PASSWORD",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: goalsbook.Name + "-sec",
												},
												Key: "password",
											},
										},
									},
									{
										Name: "DATABASE_USERNAME",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: goalsbook.Name + "-sec",
												},
												Key: "username",
											},
										},
									},
									{
										Name: "DATABASE_URL",
										ValueFrom: &corev1.EnvVarSource{
											ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: goalsbook.Name + "-cmp",
												},
												Key: "database",
											},
										},
									},
								},
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      goalsbook.Name + "-pvc",
										MountPath: "/app/logs/",
									},
								},
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: goalsbook.Name + "-pvc",
								VolumeSource: corev1.VolumeSource{
									PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
										ClaimName: goalsbook.Name + "-claim",
									},
								},
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
		log.Info("Created new Deployment", "Deployment", goalsbook.Name)
	} else {
		log.Info("Deployment already exists")

		if err = r.Update(ctx, deployment); err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment:::", goalsbook.Name, "Err:::", err.Error())
			return err
		}
		log.Info("Updated Deployment", "Deployment", goalsbook.Name)

		return nil
	}

	return nil
}
