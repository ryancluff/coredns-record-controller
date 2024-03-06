/*
Copyright 2024.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RecordNSSpec defines the desired state of RecordNS
type RecordNSSpec struct {
	Zone string `json:"zone"`
	Host string `json:"host,omitempty"`
	TTL  int    `json:"ttl,omitempty"`
}

// RecordNSStatus defines the observed state of RecordNS
type RecordNSStatus struct {
	State State `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RecordNS is the Schema for the recordns API
type RecordNS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RecordNSSpec   `json:"spec,omitempty"`
	Status RecordNSStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RecordNSList contains a list of RecordNS
type RecordNSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RecordNS `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RecordNS{}, &RecordNSList{})
}
