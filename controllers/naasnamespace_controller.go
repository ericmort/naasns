package controllers

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1alpha1 "naasns/api/v1alpha1"
)

type NaasNamespaceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *NaasNamespaceReconciler) isTenantNamespace(ctx context.Context, namespace string) (bool, error) {
	ns := &corev1.Namespace{}
	err := r.Get(ctx, client.ObjectKey{Name: namespace}, ns)
	if err != nil {
		return false, err
	}

	tenantAnnotation, exists := ns.Annotations["naas/tenant"]
	return exists && tenantAnnotation == "true", nil
}

func (r *NaasNamespaceReconciler) Handler(ctx context.Context, namespaceName, naasNamespaceName string) error {
	log := ctrl.LoggerFrom(ctx)

	// Retrieve the NaasNamespace instance
	nassNamespace := &corev1alpha1.NaasNamespace{}
	err := r.Get(ctx, client.ObjectKey{Namespace: namespaceName, Name: naasNamespaceName}, nassNamespace)

	// Define the desired Namespace object
	targetNamespaceName := fmt.Sprintf("%s-%s", namespaceName, naasNamespaceName)
	targetNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: targetNamespaceName,
		},
	}

	// Check if the Namespace already exists
	existingNamespace := &corev1.Namespace{}
	nsErr := r.Get(ctx, client.ObjectKey{Name: targetNamespace.Name}, existingNamespace)

	if err != nil && errors.IsNotFound(err) {
		// NassNamespace doesn't exist
		if nsErr == nil {
			// Namespace exists, delete it
			log.Info("Deleting namespace", "Namespace", targetNamespaceName)
			if err := r.Delete(ctx, existingNamespace); err != nil {
				log.Error(err, "Failed to delete namespace", "Namespace", targetNamespaceName)
				return err
			}
		} else if !errors.IsNotFound(nsErr) {
			// Error getting the namespace
			log.Error(nsErr, "Failed to get Namespace")
			return nsErr
		}
	} else if err == nil {
		// NaasNamespace exists
		if nsErr != nil && errors.IsNotFound(nsErr) {
			// Namespace doesn't exist, create it
			log.Info("Creating new namespace", "Namespace", targetNamespaceName)
			if err := r.Create(ctx, targetNamespace); err != nil {
				log.Error(err, "Failed to create namespace", "Namespace", targetNamespaceName)
				return err
			}
		} else if nsErr != nil {
			// Error getting the namespace
			log.Error(nsErr, "Failed to get Namespace")
			return nsErr
		}
	} else {
		// Error getting the NassNamespace
		log.Error(err, "Failed to get NassNamespace")
		return err
	}

	return nil
}

func (r *NaasNamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// Check if the NaasNamespace is created in a tenant namespace
	isTenant, err := r.isTenantNamespace(ctx, req.Namespace)
	if err != nil {
		log.Error(err, "Failed to get tenant namespace information")
		return ctrl.Result{}, err
	}

	if !isTenant {
		log.Info("Ignoring NaasNamespace resource, as it is not created in a tenant namespace", "Namespace", req.Namespace)
		return ctrl.Result{}, nil
	}

	if err := r.Handler(ctx, req.Namespace, req.Name); err != nil {
		log.Error(err, "Failed to process NaasNamespace")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *NaasNamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.NaasNamespace{}).
		Complete(r)
}
