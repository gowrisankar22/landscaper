// Copyright 2020 Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// ConditionStatus is the status of a condition.
type ConditionStatus string

// ConditionType is a string alias.
type ConditionType string

// ErrorCode is a string alias.
type ErrorCode string

const (
	// ErrorUnauthorized indicates that the last error occurred due to invalid credentials.
	ErrorUnauthorized ErrorCode = "ERR_UNAUTHORIZED"
	// ErrorCleanupResources indicates that the last error occurred due to resources are stuck in deletion.
	ErrorCleanupResources ErrorCode = "ERR_CLEANUP"
	// ErrorConfigurationProblem indicates that the last error occurred due a configuration problem.
	ErrorConfigurationProblem ErrorCode = "ERR_CONFIGURATION_PROBLEM"
)

// Condition holds the information about the state of a resource.
type Condition struct {
	// Type of the Shoot condition.
	Type ConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=ConditionType"`
	// Status of the condition, one of True, False, Unknown.
	Status ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime" protobuf:"bytes,3,opt,name=lastTransitionTime"`
	// Last time the condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime" protobuf:"bytes,4,opt,name=lastUpdateTime"`
	// The reason for the condition's last transition.
	Reason string `json:"reason" protobuf:"bytes,5,opt,name=reason"`
	// A human readable message indicating details about the transition.
	Message string `json:"message" protobuf:"bytes,6,opt,name=message"`
	// Well-defined error codes in case the condition reports a problem.
	// +optional
	Codes []ErrorCode `json:"codes,omitempty" protobuf:"bytes,7,rep,name=codes,casttype=ErrorCode"`
}


