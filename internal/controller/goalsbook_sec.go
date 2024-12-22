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

func (r *GoalsbookReconciler) DbSecretReconcile(ctx context.Context, req ctrl.Request, goalsbook *webappv1.Goalsbook) error {

	log := log.FromContext(ctx)

	secret := &corev1.Secret{}

	secretName := req.Name + "-sec"

	// _, err := GetSecret(secretName, goalsbook.Namespace)

	// if err != nil {
	// 	if errors.IsNotFound(err) {
	// 		log.Info("Secrect Not found, need to be created", secretName)
	// 		return nil
	// 	}
	// 	return err
	// }

	err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: goalsbook.Namespace}, secret)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Secret not found, creating a new one", "SecretName", secretName)

			// Create the Secret
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretName,
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
				Type: corev1.SecretTypeOpaque,
				Data: map[string][]byte{
					"password": []byte("akif"),
					"username": []byte("akif"),
				},
			}

			if err := r.Create(ctx, secret); err != nil {
				log.Error(err, "Failed to create Secret", "SecretName", secretName)
				return err
			}

			log.Info("Secret created successfully", "SecretName", secretName)
			return nil
		}
		log.Error(err, "Failed to fetch Secret")
		return err
	}

	return nil
}
