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
	infrastructurev1alpha1 "opendev.org/vexxhost/openstack-operator/api/v1alpha1"
)

// McrouterPoolSpecBuilder defines the interface to build a McrouterPoolSpec
type McrouterPoolSpecBuilder struct {
	obj *infrastructurev1alpha1.McrouterPoolSpec
}

// McrouterPoolSpec returns a new mcrouterPoolSpec builder
func McrouterPoolSpec() *McrouterPoolSpecBuilder {
	poolSpec := &infrastructurev1alpha1.McrouterPoolSpec{
		Servers: []string{},
	}
	return &McrouterPoolSpecBuilder{
		obj: poolSpec,
	}
}

// Servers specifies servers for the McrouterPoolSpec
func (ps *McrouterPoolSpecBuilder) Servers(servers []string) *McrouterPoolSpecBuilder {
	ps.obj.Servers = servers
	return ps
}

// Build returns a complete McrouterPoolSpec object
func (ps *McrouterPoolSpecBuilder) Build() (infrastructurev1alpha1.McrouterPoolSpec, error) {
	return *ps.obj, nil
}
