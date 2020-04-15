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
	"strconv"

	"k8s.io/apimachinery/pkg/util/intstr"
	monitoringv1 "opendev.org/vexxhost/openstack-operator/api/monitoring/v1"
)

// RuleBuilder provides an interface to build rule
type RuleBuilder struct {
	obj *monitoringv1.Rule
}

// Rule returns a new podmonitor builder
func Rule() *RuleBuilder {
	Rule := &monitoringv1.Rule{
		Annotations: map[string]string{},
	}
	return &RuleBuilder{
		obj: Rule,
	}
}

func (r *RuleBuilder) Alert(alert string) *RuleBuilder {
	r.obj.Alert = alert
	return r
}

func (r *RuleBuilder) Expr(expr string) *RuleBuilder {
	r.obj.Expr = intstr.FromString(expr)
	return r
}

func (r *RuleBuilder) For(duration string) *RuleBuilder {
	r.obj.For = duration
	return r
}

func (r *RuleBuilder) Priority(p int) *RuleBuilder {
	r.obj.Annotations["priority"] = "P" + strconv.Itoa(p)
	return r
}

func (r *RuleBuilder) Message(m string) *RuleBuilder {
	r.obj.Annotations["message"] = m
	return r
}

// Build returns the object after making certain assertions
func (r *RuleBuilder) Build() (monitoringv1.Rule, error) {
	return *r.obj, nil
}
