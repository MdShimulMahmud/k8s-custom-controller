package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	// corev1 "k8s.io/api/core/v1"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	webappv1 "goals.dev/goalsbook/api/v1"
)

func (r *GoalsbookReconciler) DbConfigMapReconcile(ctx context.Context, req ctrl.Request, goalsbook *webappv1.Goalsbook) error {
	log := log.FromContext(ctx)

	configMapName := goalsbook.Name + "-cmp"
	configMap := &corev1.ConfigMap{}

	err := r.Get(ctx, types.NamespacedName{Name: configMapName, Namespace: goalsbook.Namespace}, configMap)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("ConfigMap not found, creating a new one", "ConfigMapName", configMapName)

			configMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      configMapName,
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
				Data: map[string]string{
					"database": "cluster1.uktxj.mongodb.net/?retryWrites=true&w=majority&appName=Cluster1", // Example database URL
				},
			}

			// Create the ConfigMap
			if err := r.Create(ctx, configMap); err != nil {
				log.Error(err, "Failed to create ConfigMap", "ConfigMapName", configMapName)
				return err
			}

			log.Info("ConfigMap created successfully", "ConfigMapName", configMapName)
			return nil
		}
		log.Error(err, "Failed to fetch ConfigMap", "ConfigMapName", configMapName)
		return err
	}

	log.Info("ConfigMap already exists", "ConfigMapName", configMapName)

	// log.Info("ConfigMap has found, Need to be updated", "configMapName", configMapName)

	// if err := r.Status().Update(ctx, configMap); err != nil {
	// 	log.Error(err, "Error encountered while updating secret: ", configMapName)
	// 	return err
	// }
	return nil
}
