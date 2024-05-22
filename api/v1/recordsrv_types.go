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

// RecordSRVSpec defines the desired state of RecordSRV
type RecordSRVSpec struct {
	Zone     string `json:"zone,omitempty"`
	Hostname string `json:"hostname,omitempty"`

	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Weight   int    `json:"weight,omitempty"`
	TTL      int    `json:"ttl,omitempty"`
}

// RecordSRVStatus defines the observed state of RecordSRV
type RecordSRVStatus struct {
	State State `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RecordSRV is the Schema for the recordsrvs API
type RecordSRV struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RecordSRVSpec   `json:"spec,omitempty"`
	Status RecordSRVStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RecordSRVList contains a list of RecordSRV
type RecordSRVList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RecordSRV `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RecordSRV{}, &RecordSRVList{})
}
