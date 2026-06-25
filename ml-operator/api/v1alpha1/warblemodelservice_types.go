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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WarbleModelServiceSpec defines the desired state of WarbleModelService
type WarbleModelServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html

	// ModelURI is the location of the model weights (e.g., s3://my-bucket/model/)
	// +required
	ModelURI string `json:"modelURI"`

	// Replicas is the number of KServe predictor replicas to run.
	// +optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// WorkerReplicas is the number of Ray worker nodes to provision in the KubeRay cluster.
	// +optional
	// +kubebuilder:default=1
	WorkerReplicas int32 `json:"workerReplicas,omitempty"`

	// GPUEnabled specifies if the model requires GPU resources.
	// +optional
	// +kubebuilder:default=false
	GPUEnabled bool `json:"gpuEnabled,omitempty"`
}

// WarbleModelServiceStatus defines the observed state of WarbleModelService.
type WarbleModelServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// RayClusterStatus reflects the current state of the underlying RayCluster resource.
	// +optional
	RayClusterStatus string `json:"rayClusterStatus,omitempty"`

	// InferenceServiceStatus reflects the current state of the underlying KServe InferenceService resource.
	// +optional
	InferenceServiceStatus string `json:"inferenceServiceStatus,omitempty"`

	// Conditions represent the latest available observations of an object's state
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// WarbleModelService is the Schema for the warblemodelservices API
type WarbleModelService struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of WarbleModelService
	// +required
	Spec WarbleModelServiceSpec `json:"spec"`

	// status defines the observed state of WarbleModelService
	// +optional
	Status WarbleModelServiceStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// WarbleModelServiceList contains a list of WarbleModelService
type WarbleModelServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []WarbleModelService `json:"items"`
}

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(SchemeGroupVersion, &WarbleModelService{}, &WarbleModelServiceList{})
		return nil
	})
}
