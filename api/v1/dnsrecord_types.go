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
type DnsRecordSpec struct {
	// The host part of the DNS entry.
	Host string `json:"host"`

	// The domain part of the DNS entry.
	Domain string `json:"domain"`

	// The list of aliases for the DNS entry.
	Aliases []Alias `json:"aliases,omitempty"`

	// The IP address of the DNS entry.
	IPs []string `json:"ips,omitempty"`

	// The name of the load balancer to use for dynamic DNS.
	Service string `json:"service,omitempty"`

	// // The name of the kubernetes ingress to use for dynamic DNS.
	// Ingress string `json:"ingress,omitempty"`

	// // The name of the Traefik ingressRoute to use for dynamic DNS.
	// IngressRoute string `json:"ingressRoute,omitempty"`
}

// spec.alias
type Alias struct {
	// The host part of the DNS alias.
	Host string `json:"host"`

	// The domain part of the DNS alias.
	Domain string `json:"domain"`
}

// status
type DnsRecordStatus struct {
	// Specifies the state of the DnsRecord.
	// Valid values are:
	// - "Pending" (default): the controller has not processed the request yet;
	// - "Ready": the controller has created the DnsRecord;
	// - "Error": the controller encountered an error reconciling the DnsRecord and will not retry.
	State State `json:"state,omitempty"`

	// The current IP address of the DNS entry.
	IPs []string `json:"ip,omitempty"`
}

// +kubebuilder:validation:Enum=Pending;Ready;Error

// status.state
type State string

const (
	PendingState State = "Pending"
	ReadyState   State = "Ready"
	ErrorState   State = "Error"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DnsRecord is the Schema for the dnsentries API
type DnsRecord struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DnsRecordSpec   `json:"spec,omitempty"`
	Status DnsRecordStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DnsRecordList contains a list of DnsRecord
type DnsRecordList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DnsRecord `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DnsRecord{}, &DnsRecordList{})
}
