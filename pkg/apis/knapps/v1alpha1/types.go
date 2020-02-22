/*
Copyright 2017 The Kubernetes Authors.

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
)

const (
	SyncPending = "PENDING"
	SyncRunning = "RUNNING"
	SyncError   = "ERROR"
	SyncDone    = "DONE"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Argocdprog is a specification for a Argocdprog resource
type Argocdprog struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ArgocdprogSpec   `json:"spec"`
	Status ArgocdprogStatus `json:"status"`
}

// ArgocdprogSpec is the spec for a Argocdprog resource
type ArgocdprogSpec struct {
	ApiServer string `json:"apiServer"`
	Schedule  string `json:"schedule"`
	SyncApp   string `json:"syncApp"`
}

// ArgocdprogStatus is the status for a Argocdprog resource
type ArgocdprogStatus struct {
	ScheduleStatus string `json:"scheduleStatus"`
	SyncStatus     string `json:"syncStatus"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ArgocdprogList is a list of Argocdprog resources
type ArgocdprogList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Argocdprog `json:"items"`
}
