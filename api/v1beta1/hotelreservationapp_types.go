/*
Copyright 2022.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HotelReservationAppSpec defines the desired state of HotelReservationApp
type HotelReservationAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	LogicNodeIp   string `json:"logicNodeIp"`
	LogicNodeName string `json:"logicNodeName"`

	DataNodeName string `json:"dataNodeName"`
	DataNodeIp   string `json:"dataNodeIp"`
	// The mirror image corresponding to the business service, including the dockerregistryprefix
	DockerRegistryPrefix string `json:"dockerRegistryPrefix"`
}

// HotelReservationAppStatus defines the observed state of HotelReservationApp
type HotelReservationAppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes []string `json:"nodes"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HotelReservationApp is the Schema for the hotelreservationapps API
type HotelReservationApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HotelReservationAppSpec   `json:"spec,omitempty"`
	Status HotelReservationAppStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HotelReservationAppList contains a list of HotelReservationApp
type HotelReservationAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HotelReservationApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HotelReservationApp{}, &HotelReservationAppList{})
}
