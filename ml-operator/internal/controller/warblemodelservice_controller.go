/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	mlv1alpha1 "github.com/warble-platform/warble-oss-tech/ml-operator/api/v1alpha1"
)

// WarbleModelServiceReconciler reconciles a WarbleModelService object
type WarbleModelServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ml.warble.oss,resources=warblemodelservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ml.warble.oss,resources=warblemodelservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ml.warble.oss,resources=warblemodelservices/finalizers,verbs=update
// +kubebuilder:rbac:groups=ray.io,resources=rayclusters,verbs=get;list;watch;create;update;patch;delete

func (r *WarbleModelServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Fetch the WarbleModelService instance
	var modelService mlv1alpha1.WarbleModelService
	if err := r.Get(ctx, req.NamespacedName, &modelService); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Define the RayCluster dynamically
	rayCluster := &unstructured.Unstructured{}
	rayCluster.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "ray.io",
		Version: "v1alpha1",
		Kind:    "RayCluster",
	})
	rayCluster.SetName(modelService.Name + "-raycluster")
	rayCluster.SetNamespace(modelService.Namespace)

	// Construct the spec (Epic 2.1)
	spec := map[string]interface{}{
		"rayVersion": "2.9.0",
		"headGroupSpec": map[string]interface{}{
			"rayStartParams": map[string]interface{}{
				"dashboard-host": "0.0.0.0",
			},
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "ray-head",
							"image": "rayproject/ray:2.9.0",
						},
					},
				},
			},
		},
		"workerGroupSpecs": []interface{}{
			map[string]interface{}{
				"groupName":      "small-group",
				"replicas":       modelService.Spec.WorkerReplicas,
				"minReplicas":    1,
				"maxReplicas":    5,
				"rayStartParams": map[string]interface{}{},
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "ray-worker",
								"image": "rayproject/ray:2.9.0",
							},
						},
					},
				},
			},
		},
	}
	rayCluster.Object["spec"] = spec

	// Set owner reference
	if err := ctrl.SetControllerReference(&modelService, rayCluster, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Apply the RayCluster (Create or Update)
	found := &unstructured.Unstructured{}
	found.SetGroupVersionKind(rayCluster.GroupVersionKind())
	err := r.Get(ctx, client.ObjectKey{Name: rayCluster.GetName(), Namespace: rayCluster.GetNamespace()}, found)
	if err != nil && apierrors.IsNotFound(err) {
		log.Info("Creating a new RayCluster", "Namespace", rayCluster.GetNamespace(), "Name", rayCluster.GetName())
		err = r.Create(ctx, rayCluster)
		if err != nil {
			return ctrl.Result{}, err
		}
		// RayCluster created successfully
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Monitor KubeRay Resource Status (Epic 2.2)
	status, foundStatus, _ := unstructured.NestedString(found.Object, "status", "state")
	if foundStatus && modelService.Status.RayClusterStatus != status {
		modelService.Status.RayClusterStatus = status
		err := r.Status().Update(ctx, &modelService)
		if err != nil {
			log.Error(err, "Failed to update WarbleModelService status")
			return ctrl.Result{}, err
		}
		log.Info(fmt.Sprintf("Updated WarbleModelService RayClusterStatus to %s", status))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WarbleModelServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// We need to setup watch for unstructured RayClusters
	rayCluster := &unstructured.Unstructured{}
	rayCluster.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "ray.io",
		Version: "v1alpha1",
		Kind:    "RayCluster",
	})

	return ctrl.NewControllerManagedBy(mgr).
		For(&mlv1alpha1.WarbleModelService{}).
		Owns(rayCluster).
		Named("warblemodelservice").
		Complete(r)
}
