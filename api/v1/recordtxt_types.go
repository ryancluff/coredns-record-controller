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

// RecordTXTSpec defines the desired state of RecordTXT
type RecordTXTSpec struct {
	Zone string `json:"zone,omitempty"`
	Text string `json:"text,omitempty"`
	TTL  int    `json:"ttl,omitempty"`
}

// RecordTXTStatus defines the observed state of RecordTXT
type RecordTXTStatus struct {
	State State `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RecordTXT is the Schema for the recordtxts API
type RecordTXT struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RecordTXTSpec   `json:"spec,omitempty"`
	Status RecordTXTStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RecordTXTList contains a list of RecordTXT
type RecordTXTList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RecordTXT `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RecordTXT{}, &RecordTXTList{})
}
