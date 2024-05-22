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

// RecordAAAASpec defines the desired state of RecordAAAA
type RecordAAAASpec struct {
	Zone     string `json:"zone,omitempty"`
	Hostname string `json:"hostname,omitempty"`

	IP6     string `json:"ip6,omitempty"`
	Service string `json:"service,omitempty"`
	TTL     int    `json:"ttl,omitempty"`
}

// RecordAAAAStatus defines the observed state of RecordAAAA
type RecordAAAAStatus struct {
	State State  `json:"state,omitempty"`
	IP    string `json:"ip,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RecordAAAA is the Schema for the recordaaaas API
type RecordAAAA struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RecordAAAASpec   `json:"spec,omitempty"`
	Status RecordAAAAStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RecordAAAAList contains a list of RecordAAAA
type RecordAAAAList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RecordAAAA `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RecordAAAA{}, &RecordAAAAList{})
}
