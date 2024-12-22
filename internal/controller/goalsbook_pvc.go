package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	webappv1 "goals.dev/goalsbook/api/v1"
)

func (r *GoalsbookReconciler) GoalsbookPVC(ctx context.Context, req ctrl.Request, goalsbook *webappv1.Goalsbook) error {
	log := log.FromContext(ctx)

	pvc := &corev1.PersistentVolumeClaim{}
	pvcName := goalsbook.Name + "-claim"

	// Check if the PVC already exists
	err := r.Get(ctx, types.NamespacedName{
		Name:      pvcName,
		Namespace: goalsbook.Namespace,
	}, pvc)

	if err != nil && errors.IsNotFound(err) {
		log.Info("PVC not found. Creating new PVC", "PVC", pvcName)

		// Define a new PVC
		pvc := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pvcName,
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
			Spec: corev1.PersistentVolumeClaimSpec{
				StorageClassName: ptr.To("standard"),
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse("100Mi"),
					},
				},
			},
		}

		// Create the PVC
		if err = r.Create(ctx, pvc); err != nil {
			log.Error(err, "Failed to create PVC", "PVC", pvcName)
			return err
		}
		log.Info("Created new PVC", "PVC", pvcName)
		return nil
	}

	log.Info("PVC already exists", "PVC", pvcName)

	return nil
}
