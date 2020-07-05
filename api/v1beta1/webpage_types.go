/*


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

// WebPageSpec defines the desired state of WebPage
type WebPageSpec struct {
	// Html field stores the static web page contents
	// +kubebuilder:validation:MinLength=1
	Html string `json:"html"`
}

// WebPageStatus defines the observed state of WebPage
type WebPageStatus struct {
	// Stores the last time the job was successfully scheduled.
	// +optional
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// WebPage is the Schema for the webpages API
type WebPage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebPageSpec   `json:"spec,omitempty"`
	Status WebPageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WebPageList contains a list of WebPage
type WebPageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebPage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WebPage{}, &WebPageList{})
}
