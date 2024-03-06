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

// RecordMXSpec defines the desired state of RecordMX
type RecordMXSpec struct {
	Zone     string `json:"zone"`
	Host     string `json:"host,omitempty"`
	Priority int    `json:"priority,omitempty"`
	TTL      int    `json:"ttl,omitempty"`
}

// RecordMXStatus defines the observed state of RecordMX
type RecordMXStatus struct {
	State State `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RecordMX is the Schema for the recordmxes API
type RecordMX struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RecordMXSpec   `json:"spec,omitempty"`
	Status RecordMXStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RecordMXList contains a list of RecordMX
type RecordMXList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RecordMX `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RecordMX{}, &RecordMXList{})
}
