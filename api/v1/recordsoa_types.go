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

// RecordSOASpec defines the desired state of RecordSOA
type RecordSOASpec struct {
	Zone    string `json:"zone"`
	TTL     int    `json:"ttl,omitempty"`
	MBOX    string `json:"mbox,omitempty"`
	NS      string `json:"ns,omitempty"`
	Refresh int    `json:"refresh,omitempty"`
	Retry   int    `json:"retry,omitempty"`
	Expire  int    `json:"expire,omitempty"`
}

// RecordSOAStatus defines the observed state of RecordSOA
type RecordSOAStatus struct {
	State State `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RecordSOA is the Schema for the recordsoa API
type RecordSOA struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RecordSOASpec   `json:"spec,omitempty"`
	Status RecordSOAStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RecordSOAList contains a list of RecordSOA
type RecordSOAList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RecordSOA `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RecordSOA{}, &RecordSOAList{})
}
