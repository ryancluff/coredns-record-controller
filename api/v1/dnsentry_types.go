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

// spec
type DnsEntrySpec struct {
	// The host part of the DNS entry.
	Host string `json:"host"`

	// The domain part of the DNS entry.
	Domain string `json:"domain"`

	// The IP address of the DNS entry.
	Ip string `json:"ip,omitempty"`

	// The list of aliases for the DNS entry.
	Aliases []Alias `json:"aliases,omitempty"`
}

// spec.alias
type Alias struct {
	// The host part of the DNS alias.
	Host string `json:"host"`

	// The domain part of the DNS alias.
	Domain string `json:"domain"`
}

// status
type DnsEntryStatus struct {
	// Specifies the state of the DnsEntry.
	// Valid values are:
	// - "Pending" (default): the controller has not processed the request yet;
	// - "Ready": the controller has created the DnsEntry;
	// - "Error": the controller encountered an error reconciling the DnsEntry and will not retry.
	State State `json:"state,omitempty"`
}

// +kubebuilder:validation:Enum=Pending;Ready;Error

// status.state
type State string

const (
	PendingState State = "Pending"
	ReadyState   State = "Ready"
	ErrorState   State = "Error"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DnsEntry is the Schema for the dnsentries API
type DnsEntry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DnsEntrySpec   `json:"spec,omitempty"`
	Status DnsEntryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DnsEntryList contains a list of DnsEntry
type DnsEntryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DnsEntry `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DnsEntry{}, &DnsEntryList{})
}
