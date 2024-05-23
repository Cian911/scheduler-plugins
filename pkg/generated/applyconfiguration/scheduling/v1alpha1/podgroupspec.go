/*
Copyright The Kubernetes Authors.

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
)

// PodGroupSpecApplyConfiguration represents an declarative configuration of the PodGroupSpec type for use
// with apply.
type PodGroupSpecApplyConfiguration struct {
	MinMember              *int32           `json:"minMember,omitempty"`
	MinResources           *v1.ResourceList `json:"minResources,omitempty"`
	ScheduleTimeoutSeconds *int32           `json:"scheduleTimeoutSeconds,omitempty"`
}

// PodGroupSpecApplyConfiguration constructs an declarative configuration of the PodGroupSpec type for use with
// apply.
func PodGroupSpec() *PodGroupSpecApplyConfiguration {
	return &PodGroupSpecApplyConfiguration{}
}

// WithMinMember sets the MinMember field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MinMember field is set to the value of the last call.
func (b *PodGroupSpecApplyConfiguration) WithMinMember(value int32) *PodGroupSpecApplyConfiguration {
	b.MinMember = &value
	return b
}

// WithMinResources sets the MinResources field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MinResources field is set to the value of the last call.
func (b *PodGroupSpecApplyConfiguration) WithMinResources(value v1.ResourceList) *PodGroupSpecApplyConfiguration {
	b.MinResources = &value
	return b
}

// WithScheduleTimeoutSeconds sets the ScheduleTimeoutSeconds field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ScheduleTimeoutSeconds field is set to the value of the last call.
func (b *PodGroupSpecApplyConfiguration) WithScheduleTimeoutSeconds(value int32) *PodGroupSpecApplyConfiguration {
	b.ScheduleTimeoutSeconds = &value
	return b
}