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

// RecordCNAMESpec defines the desired state of RecordCNAME
type RecordCNAMESpec struct {
	Zone string `json:"zone,omitempty"`
	Host string `json:"host"`
	TTL  int    `json:"ttl,omitempty"`
}

// RecordCNAMEStatus defines the observed state of RecordCNAME
type RecordCNAMEStatus struct {
	State State `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RecordCNAME is the Schema for the recordcnames API
type RecordCNAME struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RecordCNAMESpec   `json:"spec,omitempty"`
	Status RecordCNAMEStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RecordCNAMEList contains a list of RecordCNAME
type RecordCNAMEList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RecordCNAME `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RecordCNAME{}, &RecordCNAMEList{})
}
