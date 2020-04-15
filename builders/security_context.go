// Copyright 2020 VEXXHOST, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package builders

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

// SecurityContextBuilder defines the interface to build a securityContext
type SecurityContextBuilder struct {
	obj *corev1.SecurityContext
}

// SecurityContext returns a new SecurityContext builder
func SecurityContext() *SecurityContextBuilder {
	securityContext := &corev1.SecurityContext{}
	return &SecurityContextBuilder{
		obj: securityContext,
	}
}

// RunAsUser sets the RunAsUser inside this SecurityContext
func (sc *SecurityContextBuilder) RunAsUser(userID int64) *SecurityContextBuilder {
	sc.obj.RunAsUser = pointer.Int64Ptr(userID)
	return sc
}

// RunAsGroup sets the RunAsGroup inside this SecurityContext
func (sc *SecurityContextBuilder) RunAsGroup(groupID int64) *SecurityContextBuilder {
	sc.obj.RunAsGroup = pointer.Int64Ptr(groupID)
	return sc
}

// RunAsNonRoot sets the RunAsNonRoot inside this SecurityContext
func (sc *SecurityContextBuilder) RunAsNonRoot(flag bool) *SecurityContextBuilder {
	sc.obj.RunAsNonRoot = pointer.BoolPtr(flag)
	return sc
}

// Build returns a complete ConfigMap object
func (sc *SecurityContextBuilder) Build() (corev1.SecurityContext, error) {
	return *sc.obj, nil
}
